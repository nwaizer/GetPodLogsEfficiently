package main

import (
	"GetPodLogsEfficiently/client"
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/sirupsen/logrus"
)

const (
	Namespace     = "default"
	LabelSelector = "app=logsdemo"
)

var Client *client.Set

func Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}

	return false
}

func VerifyEvents(podsList *corev1.PodList, consumerOutChannel chan string, timeoutDuration int) bool {
	results := make(map[string][]string)
	ctx := context.Background()

	timeout, cancelTimeout := context.WithTimeout(ctx, time.Duration(timeoutDuration)*time.Second)

	// End successfully after these lines are verified:
	expectedEvents := len(podsList.Items) * timeoutDuration / 2

	var testingOk string

	for {
		select {
		case <-timeout.Done():
			logrus.Errorf("timeout expired")
			logrus.Warnf("Pods that did not logrus any string: %v\n", len(podsList.Items)-len(results))

			unVerifiedPods := getUnVerifiedPods(podsList, results)

			cancelTimeout()

			if len(unVerifiedPods) > 0 {
				return false
			}

			return true

		case testingOk = <-consumerOutChannel:
			if testingOk[len(testingOk)-2:] != "OK" {
				logrus.Errorf("verifier got bad msg: %v\n", testingOk)
				cancelTimeout()

				return false
			}

			logrus.Println(testingOk)
			verified := strings.Split(testingOk, "/")
			results[verified[0]] = append(results[verified[0]], verified[1])

			if len(results) == expectedEvents {
				cancelTimeout()

				return true
			}
		}
	}
}

func getUnVerifiedPods(podsList *corev1.PodList, results map[string][]string) []string {
	var podsNames []string

	for _, result := range podsList.Items {
		podsNames = append(podsNames, result.Name)
	}

	var verifiedPodsList []string

	for result := range results {
		verifiedPodsList = append(verifiedPodsList, result)
	}

	var unverifiredPods []string

	for _, consumer := range podsNames {
		if !Contains(verifiedPodsList, consumer) {
			unverifiredPods = append(unverifiredPods, consumer)
		}
	}

	logrus.Printf("Unverified pods: %v\n", unverifiredPods)

	return unverifiredPods
}

func Checker(cancelCtx context.Context, podName string, outChannel chan string) {
	PodLogsConnection := Client.Pods(Namespace).GetLogs(podName, &corev1.PodLogOptions{
		Follow:    true,
		TailLines: &[]int64{int64(10)}[0],
	})
	LogStream, _ := PodLogsConnection.Stream(context.Background())

	defer LogStream.Close()

	reader := bufio.NewScanner(LogStream)

	var expectedString = "Good"

	var line string

	for {
		for reader.Scan() {
			select {
			case <-cancelCtx.Done():
				break
			default:
				line = reader.Text()
				logrus.Debugln(line)
				// check the logrus line
				if strings.Contains(line, expectedString) {
					logrus.Infof("Pod %v got: %v\t", podName, expectedString)
					outChannel <- fmt.Sprintf("%v/%v/OK", podName, expectedString)
				} else {
					logrus.Errorf("Pod %v did not get string %v \n", podName, expectedString)
					outChannel <- fmt.Sprintf("%v/%v/FAIL", podName, expectedString)

					return
				}
			}

			if reader.Err() != nil {
				logrus.Errorln(reader.Err())
			}
		}
	}
}

func SetDebuglogLevel() {
	lvl, ok := os.LookupEnv("LOG_LEVEL")
	// LOG_LEVEL not set, let's default to debug
	if !ok {
		lvl = "debug"
	}
	// parse string, this is built-in feature of logrus
	logLevel, err := logrus.ParseLevel(lvl)
	if err != nil {
		logLevel = logrus.DebugLevel
	}
	// set global logrus level
	logrus.SetLevel(logLevel)
}

func main() {
	SetDebuglogLevel() // set logrus level to debug

	Client = client.New("")
	ctx := context.Background()
	cancelCtx, endCheckers := context.WithCancel(ctx)

	pods, err := client.Client.Pods(Namespace).List(context.Background(), v12.ListOptions{
		LabelSelector: LabelSelector})

	if err != nil {
		logrus.Errorf("Failed to get demo pod from cluster %v", client.Client.Config.Host)

		return
	}

	outChannel := make(chan string, len(pods.Items))

	for _, Pod := range pods.Items {
		go Checker(cancelCtx, Pod.Name, outChannel)
	}

	if !VerifyEvents(pods, outChannel, 10) {
		endCheckers()
		logrus.Errorf("Test failed due to error printed above")
	}
}
