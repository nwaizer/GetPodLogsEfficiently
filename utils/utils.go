package utils

import (
	"GetPodLogsEfficiently/client"
	"context"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
)

var Namespace = "default"

func GetPods() *v1.PodList {
	Pods, err := client.Client.Pods(Namespace).List(context.Background(), v12.ListOptions{
		LabelSelector: "app=logsdemo"})
	if err != nil {
		logrus.Errorf("Failed to get demo pod from cluster %v", client.Client.Config.Host)
	}
	return Pods
}

func SetDebuglogLevel() {
	lvl, ok := os.LookupEnv("LOG_LEVEL")
	// LOG_LEVEL not set, let's default to debug
	if !ok {
		lvl = "debug"
	}
	// parse string, this is built-in feature of logrus
	ll, err := logrus.ParseLevel(lvl)
	if err != nil {
		ll = logrus.DebugLevel
	}
	// set global log level
	logrus.SetLevel(ll)
}
