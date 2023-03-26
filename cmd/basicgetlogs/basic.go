package main

import (
	"GetPodLogsEfficiently/client"
	"bytes"
	"context"
	"fmt"

	"io"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

const (
	Namespace     = "default"
	LabelSelector = "app=logsdemo"
)

var Client *client.Set

func getPodLogs(pod corev1.Pod) error {
	req := Client.Pods(pod.Namespace).GetLogs(pod.Name, &corev1.PodLogOptions{})
	podLogs, err := req.Stream(context.Background())

	if err != nil {
		return fmt.Errorf("failed to open pod logs due to: %w", err)
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)

	if err != nil {
		return fmt.Errorf("failed to decode logs binary input due to: %w", err)
	}

	logrus.Print(buf.String())

	return err
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
		err := getPodLogs(pod)
		// We expect no error to occur.
		if err != nil {
			logrus.Println(err)

			break
		}
	}

}
