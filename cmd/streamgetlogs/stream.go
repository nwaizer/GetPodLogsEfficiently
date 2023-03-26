package main

import (
	"GetPodLogsEfficiently/client"
	"bufio"
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	Namespace     = "default"
	LabelSelector = "app=logsdemo"
)

var Client *client.Set

func GetPodLogs(podName string) error {
	k8session := client.New("")
	podLogOpts := corev1.PodLogOptions{}
	podLogOpts.Follow = true
	tailLines := int64(100)
	podLogOpts.TailLines = &tailLines // get 100 last lines
	podLogs, err := k8session.CoreV1Interface.Pods(
		Namespace).GetLogs(podName, &podLogOpts).Stream(context.Background())

	if err != nil {
		return fmt.Errorf("failed to get logs for pod: %v due to: %w", podName, err)
	}
	defer podLogs.Close()

	for i := 0; i < 10; i++ {
		reader := bufio.NewScanner(podLogs)
		for reader.Scan() {
			line := reader.Text()
			logrus.Printf("Pod: %v line: %v\n", podName, line)
		}
		time.Sleep(1 * time.Millisecond)

		if reader.Err() != nil {
			return fmt.Errorf("caught an error while processing logs due to: %w", reader.Err())
		}
	}

	return nil
}

func main() {
	Client = client.New("")
	pods, err := Client.Pods(Namespace).List(context.Background(), v12.ListOptions{
		LabelSelector: LabelSelector})

	if err != nil {
		logrus.Errorf("Failed to get demo pod from cluster %v", client.Client.Config.Host)

		return
	}

	for _, pod := range pods.Items {
		logrus.Println(pod.Name)
		err := GetPodLogs(pod.Name)

		if err != nil {
			logrus.Println(err)
		}
	}
}
