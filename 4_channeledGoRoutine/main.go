package main

import (
	"GetPodLogsEfficiently/client"
	"GetPodLogsEfficiently/utils"
	"bufio"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"runtime/pprof"
	"strings"
	"time"
)

func init() {
	utils.SetDebuglogLevel()
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
	timeout, cancelTimeout := context.WithTimeout(ctx, time.Duration(timeoutDuration)*time.Millisecond)
	var testingOk string
	for {
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
	req := client.Client.Pods(utils.Namespace).GetLogs(Pod.Name, &corev1.PodLogOptions{
		Follow:    true,
		TailLines: &[]int64{int64(10)}[0],
	})
	LogStream, _ := req.Stream(context.Background())

	scanner := bufio.NewScanner(LogStream)
	var expectedString = "Good"
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
					break
				}
			}
			if scanner.Err() != nil {
				log.Errorln(scanner.Err())
			}
		}
	}
}

func main() {
	PodsList := utils.GetPods()
	ctx := context.Background()
	cancelCtx, endCheckers := context.WithCancel(ctx)

	outChannel := make(chan string, len(PodsList.Items))

	for _, aPod := range PodsList.Items {
		labels := pprof.Labels("request", "automated", "pod", aPod.Name)
		pprof.Do(context.Background(), labels, func(_ context.Context) {
			go Checker(cancelCtx, aPod, outChannel)
		})
	}

	if !VerifyEvents(PodsList, outChannel, 100) {
		endCheckers()
		log.Errorf("Test failed due to error printed above")
	}
}
