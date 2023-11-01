package consumer

import (
	"context"
	"fmt"

	"github.com/noisyboy-9/random-k8s-scheduler/internal/config"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/connector"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/enum"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/log"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/model"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/service"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/util"
	"github.com/sirupsen/logrus"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type Consumer struct {
	nodes []*model.Node
	pods  []*model.Pod
}

func Init() *Consumer {
	consumer := &Consumer{
		nodes: make([]*model.Node, 0),
		pods:  make([]*model.Pod, 0),
	}

	nodeList, err := connector.ClusterConnection.Client().CoreV1().Nodes().List(context.Background(), metaV1.ListOptions{})
	if err != nil {
		log.App.WithError(err).Panic("can't get node list")
	}

	for _, node := range nodeList.Items {
		consumer.nodes = append(consumer.nodes, model.NewNode(
			string(node.GetUID()),
			node.Name,
			node.Status.Allocatable.Cpu(),
			node.Status.Allocatable.Memory(),
		))

		log.App.WithFields(logrus.Fields{
			"node_name":          node.Name,
			"current_node_count": len(consumer.nodes),
			"resources": map[string]string{
				"cpu":    node.Status.Allocatable.Cpu().String(),
				"memory": node.Status.Allocatable.Memory().String(),
			},
		}).Info("node has been added")
	}

	log.App.WithFields(logrus.Fields{"nodes_count": len(consumer.nodes)}).Info("List of nodes has been added")
	return consumer
}

func (consumer *Consumer) Consume() {
	eventQueue, err := connector.ClusterConnection.Client().CoreV1().Pods(config.Scheduler.Namespace).Watch(
		context.Background(),
		metaV1.ListOptions{
			FieldSelector: fmt.Sprintf("spec.schedulerName=%s", config.Scheduler.Name),
		},
	)

	if err != nil {
		log.App.WithError(err).Panic("can't watch pod event queue")
	}

	log.App.Info("watching for pod events . . .")
	for event := range eventQueue.ResultChan() {
		log.App.Info("got pod event!")
		if event.Type != watch.Added {
			log.App.WithFields(logrus.Fields{"event_type": event.Type}).Info("event wasn't related to pod creation, ignoring event . . .")
			continue
		}

		podSpec, ok := event.Object.(*coreV1.Pod)
		if !ok {
			log.App.Error("unexpected event object type")
		}
		newPod := model.NewPod(
			string(podSpec.GetUID()),
			podSpec.Name,
			podSpec.Namespace,
			util.RequiredCpuSum(podSpec.Spec.Containers),
			util.RequiredMemorySum(podSpec.Spec.Containers),
		)
		consumer.pods = append(consumer.pods, newPod)

		selectedNode, err := service.Scheduler.FindNodeForBinding(newPod, consumer.nodes)
		if err != nil {
			log.App.WithError(err).WithFields(logrus.Fields{"pod": newPod.Id()}).Error("error in finding node for pod")
		}

		log.App.WithFields(logrus.Fields{
			"pod":           newPod.Id(),
			"selected_node": selectedNode.Id(),
		}).Info("node selected for pod")

		if err := connector.ClusterConnection.BindPodToNode(newPod, selectedNode); err != nil {
			log.App.WithError(err).Error("pod bind error")
			continue
		}
		if err := connector.ClusterConnection.EmitScheduledEvent(newPod, selectedNode); err != nil {
			log.App.WithError(err).Error("scheduled event emit error")
		}

		consumer.UpdateClusterStateAfterBinding(newPod, selectedNode)
	}
}

func (consumer *Consumer) UpdateClusterStateAfterBinding(pod *model.Pod, node *model.Node) {
	pod.SetStatus(enum.PodStatusRunning)
	pod.SetNode(node)
	node.ReduceAllocateableCpu(pod.Cpu())
	node.ReduceAllocateableMemory(pod.Memory())
}
