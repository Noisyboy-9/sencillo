package connector

import (
	"context"
	"fmt"
	"k8s.io/client-go/dynamic"
	"time"

	"github.com/noisyboy-9/random-k8s-scheduler/internal/config"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/log"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/model"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type connector struct {
	client        *kubernetes.Clientset
	dynamicConfig *dynamic.DynamicClient
}

var C *connector

func Connect() {
	log.App.Info("connecting to Kubernetes cluster ...")
	C = new(connector)

	var err error
	C.client, err = createInClusterClient()
	if err != nil {
		log.App.WithError(err).Panic("can't create cluster client")
	}

	C.dynamicConfig, err = createInClusterDynamicClient()
	if err != nil {
		log.App.WithError(err).Panic("can't create cluster dynamic client")
	}

	log.App.Info("successfully connected to Kubernetes cluster")
}

func createInClusterDynamicClient() (*dynamic.DynamicClient, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	return dynamic.NewForConfig(config)
}

func (connector *connector) Client() *kubernetes.Clientset {
	return connector.client
}

func (connector *connector) BindPodToNode(pod *model.Pod, selectedNode *model.Node) error {
	return C.Client().CoreV1().Pods(pod.Namespace()).Bind(
		context.Background(),
		&v1.Binding{
			ObjectMeta: metaV1.ObjectMeta{
				Name:      pod.Name(),
				Namespace: pod.Namespace(),
			},
			Target: v1.ObjectReference{
				APIVersion: "v1",
				Kind:       "Node",
				Name:       selectedNode.Name(),
			},
		},
		metaV1.CreateOptions{},
	)

}

func (connector *connector) EmitScheduledEvent(pod *model.Pod, node *model.Node) error {
	timestamp := time.Now().UTC()
	_, err := C.Client().CoreV1().Events(pod.Namespace()).Create(
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
				UID:       pod.ID(),
			},
			ObjectMeta: metaV1.ObjectMeta{
				GenerateName: pod.Name() + "-",
			},
		},
		metaV1.CreateOptions{},
	)
	return err
}
func (connector *connector) DynamicConfig() *dynamic.DynamicClient {
	return connector.dynamicConfig
}
func createInClusterClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}
