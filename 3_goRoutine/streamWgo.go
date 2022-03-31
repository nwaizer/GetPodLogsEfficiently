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

func GetPodLogs(cancelCtx context.Context, PodName string) {
	PodLogsConnection := client.Client.Pods(utils.Namespace).GetLogs(PodName, &corev1.PodLogOptions{
		Follow:    true,
		TailLines: &[]int64{int64(10)}[0],
	})
	LogStream, _ := PodLogsConnection.Stream(context.Background())
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
				fmt.Printf("Pod: %v line: %v\n", PodName, line)
			}
		}
	}
}
func main() {
	ctx := context.Background()
	cancelCtx, endGofunc := context.WithCancel(ctx)
	for _, pod := range utils.GetPods().Items {
		fmt.Println(pod.Name)
		go GetPodLogs(cancelCtx, pod.Name)
	}
	time.Sleep(10 * time.Second)
	endGofunc()
}
