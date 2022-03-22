package main

import (
	"GetPodLogsEfficiently/client"
	"GetPodLogsEfficiently/utils"
	"bytes"
	"context"
	"fmt"
	"io"
	corev1 "k8s.io/api/core/v1"
)

func getPodLogs(pod corev1.Pod) string {
	req := client.Client.Pods(pod.Namespace).GetLogs(pod.Name, &corev1.PodLogOptions{})
	podLogs, err := req.Stream(context.Background())
	if err != nil {
		return "error in opening pod logs"
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "error in copy information from podLogs to buf"
	}
	str := buf.String()
	return str
}

func main() {
	for _, pod := range utils.GetPods().Items {
		fmt.Println(pod.Name)
		fmt.Printf(getPodLogs(pod))
	}
}
