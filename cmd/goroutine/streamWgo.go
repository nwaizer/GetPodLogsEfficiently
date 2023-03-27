package main

import (
	"GetPodLogsEfficiently/client"
	"bufio"
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	Namespace     = "default"
	LabelSelector = "app=logsdemo"
)

var Client *client.Set

func GetPodLogs(cancelCtx context.Context, podName string) {
	k8session := client.New("")
	PodLogsConnection := k8session.Pods(Namespace).GetLogs(podName, &corev1.PodLogOptions{
		Follow:    true,
		TailLines: &[]int64{int64(10)}[0],
	})
	LogStream, err := PodLogsConnection.Stream(context.Background())

	if err != nil {
		logrus.Error(err)

		return
	}
	defer LogStream.Close()

	reader := bufio.NewScanner(LogStream)

	var line string

	for {
		for reader.Scan() {
			select {
			case <-cancelCtx.Done():
				break
			default:
				line = reader.Text()
				fmt.Printf("Pod: %v line: %v\n", podName, line)
			}
		}

		if reader.Err() != nil {
			logrus.Printf("error in logs inpput for pod: %v due to: %v\n", podName, reader.Err())

			break
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	Client = client.New("")
	pods, err := Client.Pods(Namespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: LabelSelector})

	if err != nil {
		logrus.Errorf("Failed to get demo pod from cluster %v", client.Client.Config.Host)

		return
	}

	for _, pod := range pods.Items {
		logrus.Println(pod.Name)

		go GetPodLogs(ctx, pod.Name)
	}

	time.Sleep(10 * time.Second)
	cancel()
}
