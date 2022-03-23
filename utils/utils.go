package utils

import (
	"GetPodLogsEfficiently/client"
	"context"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Namespace = "default"

func GetPods() *v1.PodList {
	Pods, err := client.Client.Pods(Namespace).List(context.Background(), v12.ListOptions{
		LabelSelector: "app: logsdemo"})
	if err != nil {
		logrus.Errorf("Failed to get demo pod from cluster %v", client.Client.Config.Host)
	}
	return Pods
}
