package main

import (
	"GetPodLogsEfficiently/client"
	"bytes"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	corev1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var NameSpace = "default"

func getPodLogs() string {
	Pods, err := client.Client.Pods(NameSpace).List(context.Background(), metaV1.ListOptions{
		LabelSelector: "app=demo"})
	if err != nil {
		log.Errorf("Failed to get demo pod from cluster %v\n", client.Client.Config.Host)
	}
	pod := Pods.Items[0]
	podLogOpts := corev1.PodLogOptions{}

	req := client.Client.Pods(pod.Namespace).GetLogs(pod.Name, &podLogOpts)
	podLogs, err := req.Stream(context.Background())
	if err != nil {
		return "error in opening stream"
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
	fmt.Printf(getPodLogs())
}
