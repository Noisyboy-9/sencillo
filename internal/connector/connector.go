package connector

import (
	"github.com/noisyboy-9/random-k8s-scheduler/internal/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type connector struct {
	client *kubernetes.Clientset
}

var ClusterConnection *connector

func createInClusterClientset() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

func Connect() {
	log.App.Info("connecting to Kubernetes cluster ...")
	Connector := new(connector)

	var err error
	Connector.client, err = createInClusterClientset()
	if err != nil {
		log.App.WithError(err).Panic("can't connect to k8s cluster")
	}

	log.App.Info("successfully connected to Kubernetes cluster")
}

func (connector *connector) Client() *kubernetes.Clientset {
	return connector.client
}
