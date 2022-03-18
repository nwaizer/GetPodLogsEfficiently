package main

import (
	"GetPodLogsEfficiently/client"
	"bufio"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"strings"
	"time"
)

var Namespace = "hw-event-proxy-operator-system"

func init() {
	lvl, ok := os.LookupEnv("LOG_LEVEL")
	// LOG_LEVEL not set, let's default to debug
	if !ok {
		lvl = "debug"
	}
	// parse string, this is built-in feature of logrus
	ll, err := log.ParseLevel(lvl)
	if err != nil {
		ll = log.DebugLevel
	}
	// set global log level
	log.SetLevel(ll)
}

func Contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func VerifyEvents(PodsList *corev1.PodList, consumerOutChannel chan string, timeoutDuration int) bool {
	results := make(map[string][]string)
	ctx := context.Background()
	timeout, cancelTimeout := context.WithTimeout(ctx, time.Duration(timeoutDuration)*time.Second)
	var testingOk string
	for i := 0; i < len(PodsList.Items); i++ {
		select {
		case <-timeout.Done():
			log.Errorln("Timeout expired")
			log.Warnf("Pods that did not log any string: %v\n", len(PodsList.Items)-len(results))
			podsNames := []string{}
			for _, result := range PodsList.Items {
				podsNames = append(podsNames, result.Name)
			}
			verifiedPodsList := []string{}
			for result := range results {
				verifiedPodsList = append(verifiedPodsList, result)
			}
			unverifiredPods := []string{}
			for _, consumer := range podsNames {
				if !Contains(verifiedPodsList, consumer) {
					unverifiredPods = append(unverifiredPods, consumer)
				}
			}
			log.Warnf("Unverified pods: %v\n", unverifiredPods)
			cancelTimeout()
			return false
		case testingOk = <-consumerOutChannel:
			if testingOk[len(testingOk)-2:] == "OK" {
				log.Debugln(testingOk)
				verified := strings.Split(testingOk, "/")
				results[verified[0]] = append(results[verified[0]], verified[1])
			} else {
				log.Errorln(testingOk)
				cancelTimeout()
				return false
			}
		}
	}
	cancelTimeout()
	return true
}

func Checker(cancelCtx context.Context, Pod corev1.Pod, outch chan string) {
	req := client.Client.Pods(Namespace).GetLogs(Pod.Name, &corev1.PodLogOptions{
		Follow: true,
	})
	LogStream, _ := req.Stream(context.Background())

	scanner := bufio.NewScanner(LogStream)
	var expectedString = "Demoi"
	var line string

	for {
		select {
		case <-cancelCtx.Done():
			break
		default:
			for scanner.Scan() {
				line = scanner.Text()
				log.Debugln(line)
				// check the log line
				if strings.Contains(line, expectedString) {
					log.Infof("Pod %v got: %v\t", Pod.Name, expectedString)
					outch <- fmt.Sprintf("%v/%v/OK", Pod.Name, expectedString)
				} else {
					log.Errorf("Pod %v did not get string %v \n", Pod.Name, expectedString)
					outch <- fmt.Sprintf("%v/%v/FAIL", Pod.Name, expectedString)
				}
			}
			if scanner.Err() != nil {
				log.Errorln(scanner.Err())
			}
		}
	}
}

func main() {
	PodsList := GetPod(Namespace)
	ctx := context.Background()
	cancelCtx, endCheckers := context.WithCancel(ctx)

	outChannel := make(chan string, len(PodsList.Items))

	for _, aPod := range PodsList.Items {
		go Checker(cancelCtx, aPod, outChannel)
	}

	if !VerifyEvents(PodsList, outChannel, 180) {
		endCheckers()
		return
	}
}

func GetPod(NameSpace string) *corev1.PodList {
	Pods, err := client.Client.Pods(NameSpace).List(context.Background(), v1.ListOptions{
		LabelSelector: "app=demo"})
	if err != nil {
		log.Errorf("Failed to get demo pod from cluster %v", client.Client.Config.Host)
	}
	return Pods
}
