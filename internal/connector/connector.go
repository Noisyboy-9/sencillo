package connector

import (
	"context"
	"fmt"
	"time"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"

	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/noisyboy-9/sencillo/internal/config"
	"github.com/noisyboy-9/sencillo/internal/log"
	"github.com/noisyboy-9/sencillo/internal/model"
	"github.com/noisyboy-9/sencillo/internal/util"
)

type connector struct {
	client        *kubernetes.Clientset
	dynamicConfig *dynamic.DynamicClient
}

var C *connector

func Connect() {
	log.App.Info("connecting to Kubernetes cluster ...")
	C = new(connector)

	log.App.Info(config.Connector.Mode)

	var err error
	if config.Connector.Mode == "outside" {
		C.client, err = createOutsideClusterClient()
		if err != nil {
			log.App.WithError(err).Panic("can't create cluster client")
		}
	} else {
		C.client, err = createInClusterClient()
		if err != nil {
			log.App.WithError(err).Panic("can't create cluster client")
		}

		C.dynamicConfig, err = createInClusterDynamicClient()
		if err != nil {
			log.App.WithError(err).Panic("can't create cluster dynamic client")
		}
	}

	log.App.Info("successfully connected to Kubernetes cluster")
}

func createInClusterDynamicClient() (*dynamic.DynamicClient, error) {
	c, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	return dynamic.NewForConfig(c)
}

func (connector *connector) Client() *kubernetes.Clientset {
	return connector.client
}

func (connector *connector) BindPodToNode(pod model.Pod, selectedNode model.Node) error {
	return C.Client().CoreV1().Pods(pod.Namespace).Bind(
		context.Background(),
		&v1.Binding{
			ObjectMeta: metaV1.ObjectMeta{
				Name:      pod.Name,
				Namespace: pod.Namespace,
			},
			Target: v1.ObjectReference{
				APIVersion: "v1",
				Kind:       "Node",
				Name:       selectedNode.Name,
			},
		},
		metaV1.CreateOptions{},
	)

}

func (connector *connector) EmitScheduledEvent(pod model.Pod, node model.Node) error {
	timestamp := time.Now().UTC()
	_, err := C.Client().CoreV1().Events(pod.Namespace).Create(
		context.Background(),
		&v1.Event{
			Count:          1,
			Message:        fmt.Sprintf("pod: %s has been bound to node: %s", pod.Name, node.Name),
			Reason:         "Scheduled",
			LastTimestamp:  metaV1.NewTime(timestamp),
			FirstTimestamp: metaV1.NewTime(timestamp),
			Type:           "Normal",
			Source: v1.EventSource{
				Component: config.Scheduler.Name,
			},
			InvolvedObject: v1.ObjectReference{
				Kind:      "Pod",
				Name:      pod.Name,
				Namespace: pod.Namespace,
				UID:       pod.ID,
			},
			ObjectMeta: metaV1.ObjectMeta{
				GenerateName: pod.Name + "-",
			},
		},
		metaV1.CreateOptions{},
	)
	return err
}
func (connector *connector) DynamicConfig() *dynamic.DynamicClient {
	return connector.dynamicConfig
}

func (connector *connector) GetAllPods() ([]model.Pod, error) {
	list, err := connector.client.CoreV1().Pods(config.Scheduler.Namespace).List(context.Background(), metaV1.ListOptions{})
	if err != nil {
		return nil, err
	}

	pods := make([]model.Pod, len(list.Items))
	for _, item := range list.Items {
		pod := model.NewPod(
			item.UID,
			item.Name,
			item.Namespace,
			item.Spec.NodeName,
			util.RequiredCpuSum(item.Spec.Containers),
			util.RequiredMemorySum(item.Spec.Containers),
		)

		pods = append(pods, pod)
	}

	return pods, nil
}

func createInClusterClient() (*kubernetes.Clientset, error) {
	c, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(c)
}

func createOutsideClusterClient() (*kubernetes.Clientset, error) {
	c, err := clientcmd.BuildConfigFromFlags("", config.Connector.KubeConfigPath)
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(c)
}
