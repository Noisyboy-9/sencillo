package service

import (
	"github.com/noisyboy-9/random-k8s-scheduler/internal/config"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type scheduler struct {
	Name   string
	Client *kubernetes.Clientset
}

var Scheduler *scheduler

func createInClusterClientset() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

func NewScheduler() {
	log.App.Info("connecting to Kubernetes cluster ...")
	clientset, err := createInClusterClientset()
	if err != nil {
		log.App.WithError(err).Panic("can't connect to k8s cluster")
	}

	Scheduler = &scheduler{
		Name:   config.Scheduler.Name,
		Client: clientset,
	}
	log.App.Info("successfully connected to Kubernetes cluster")
}
