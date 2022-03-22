package main

import (
	"GetPodLogsEfficiently/client"
	"GetPodLogsEfficiently/utils"
	"bufio"
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"time"
)

func GetPodLogs(podName string) error {
	podLogOpts := corev1.PodLogOptions{}
	podLogOpts.Follow = true
	podLogOpts.TailLines = &[]int64{int64(100)}[0] //get 100 last lines
	podLogs, err := client.Client.CoreV1Interface.Pods(utils.Namespace).GetLogs(podName, &podLogOpts).Stream(context.Background())
	if err != nil {
		return err
	}
	defer podLogs.Close()
	for i := 0; i < 10; i++ {
		reader := bufio.NewScanner(podLogs)
		for reader.Scan() {
			line := reader.Text()
			fmt.Printf("Pod: %v line: %v\n", podName, line)
		}
		time.Sleep(1 * time.Millisecond)
	}
	return nil
}

func main() {
	for _, pod := range utils.GetPods().Items {
		fmt.Println(pod.Name)
		GetPodLogs(pod.Name)
	}
}
