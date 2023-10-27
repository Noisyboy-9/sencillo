package connector

import (
	"context"
	"fmt"
	"time"

	"github.com/noisyboy-9/random-k8s-scheduler/internal/config"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/log"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/model"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
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

func (connector *connector) BindPodToNode(pod *model.Pod, selectedNode *model.Node) error {
	return ClusterConnection.Client().CoreV1().Pods(pod.Namespace()).Bind(
		context.Background(),
		&v1.Binding{
			ObjectMeta: metaV1.ObjectMeta{
				Name:      pod.Name(),
				Namespace: pod.Namespace(),
			},
			Target: v1.ObjectReference{
				APIVersion: "v1",
				Kind:       "node",
				Name:       selectedNode.Name(),
			},
		},
		metaV1.CreateOptions{},
	)

}

func (connector *connector) EmitScheduledEvent(pod *model.Pod, node *model.Node) error {
	timestamp := time.Now().UTC()
	_, err := ClusterConnection.Client().CoreV1().Events(pod.Namespace()).Create(
		context.Background(),
		&v1.Event{
			Count:          1,
			Message:        fmt.Sprintf("pod: %s has been bound to node: %s", pod.Name(), node.Name()),
			Reason:         "Scheduled",
			LastTimestamp:  metaV1.NewTime(timestamp),
			FirstTimestamp: metaV1.NewTime(timestamp),
			Type:           "Normal",
			Source: v1.EventSource{
				Component: config.Scheduler.Name,
			},
			InvolvedObject: v1.ObjectReference{
				Kind:      "Pod",
				Name:      pod.Name(),
				Namespace: pod.Namespace(),
				UID:       types.UID(pod.Id()),
			},
			ObjectMeta: metaV1.ObjectMeta{
				GenerateName: pod.Name() + "-",
			},
		},
		metaV1.CreateOptions{},
	)
	return err
}
