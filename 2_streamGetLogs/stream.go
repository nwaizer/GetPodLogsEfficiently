package main

import (
	"GetPodLogsEfficiently/client"
	"bufio"
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"time"
)

var namespace = "default"

func GetPodLogs(namespace string, podName string, container string) error {
	clientSet := client.New("")

	podLogOpts := v1.PodLogOptions{}
	podLogOpts.Follow = true
	podLogOpts.TailLines = &[]int64{int64(100)}[0]
	podLogOpts.Container = container
	podLogs, err := clientSet.CoreV1Interface.Pods(namespace).GetLogs(podName, &podLogOpts).Stream(context.Background())
	if err != nil {
		return err
	}
	defer podLogs.Close()
	for i := 0; i < 10; i++ {
		reader := bufio.NewScanner(podLogs)
		for reader.Scan() {
			line := reader.Text()
			fmt.Println("worker"+"/"+podLogOpts.Container, line)
		}
		time.Sleep(1 * time.Millisecond)
	}
	return nil
}

func main() {
	GetPodLogs(namespace, "get-pod-logs-efficiently", "podslogs")
}
