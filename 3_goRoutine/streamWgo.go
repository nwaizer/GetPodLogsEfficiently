package main

import (
	"GetPodLogsEfficiently/client"
	"bufio"
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"time"
)

var namespace = "hw-event-proxy-operator-system"

func GetPodLogs(namespace string, podName string, container string) error {
	clientSet := client.New("")
	//pod, err := clientSet.CoreV1Interface.Pods(namespace).Get(context.Background(), podName, metav1.GetOptions{})
	//if err != nil {
	//	return err
	//}

	podLogOpts := v1.PodLogOptions{}
	podLogOpts.Follow = true
	podLogOpts.TailLines = &[]int64{int64(100)}[0]
	podLogOpts.Container = container
	podLogs, err := clientSet.CoreV1Interface.Pods(namespace).GetLogs(podName, &podLogOpts).Stream(context.Background())
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
				fmt.Println("worker"+"/"+podLogOpts.Container, line)
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
	GetPodLogs(namespace, "get-pod-logs-efficiently", "podslogs")
}
