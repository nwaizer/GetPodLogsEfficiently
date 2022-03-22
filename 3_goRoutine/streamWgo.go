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
	podLogOpts.TailLines = &[]int64{int64(100)}[0]
	podLogs, err := client.Client.CoreV1Interface.Pods(utils.Namespace).GetLogs(podName, &podLogOpts).Stream(context.Background())
	if err != nil {
		return err
	}
	defer podLogs.Close()
	ctx := context.Background()
	cancelCtx, endGofunc := context.WithCancel(ctx)
	go func(cancelCtx context.Context) {
		reader := bufio.NewScanner(podLogs)
		for reader.Scan() {
			select {
			case <-cancelCtx.Done():
				return
			default:
				line := reader.Text()
				fmt.Printf("Pod: %v line: %v\n", podName, line)
			}
		}
	}(cancelCtx)
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Millisecond)
	}
	endGofunc()
	return nil
}
func main() {
	for _, pod := range utils.GetPods().Items {
		fmt.Println(pod.Name)
		GetPodLogs(pod.Name)
	}
}
