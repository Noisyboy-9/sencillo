package consumer

import (
	"context"
	"fmt"
	"sync"

	"github.com/noisyboy-9/random-k8s-scheduler/internal/config"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/connector"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/log"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/model"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/service"
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
		if node.Status.Phase != coreV1.NodeRunning {
			// don't add an unready node to list of nodes
			continue
		}

		consumer.nodes = append(consumer.nodes, model.NewNode(
			string(node.GetUID()),
			node.Status.Allocatable.Cpu().AsApproximateFloat64(),
			node.Status.Allocatable.Memory().AsApproximateFloat64(),
		))

		log.App.WithFields(logrus.Fields{"nodes_count": len(consumer.nodes)}).Info("List of nodes has been added")
	}

	return consumer
}

func (consumer *Consumer) Consume() {
	eventQueue, err := connector.ClusterConnection.Client().CoreV1().Pods(config.Scheduler.Namespace).Watch(context.Background(), metaV1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.schedulerName=%s,space.nodeName=", config.Scheduler.Name),
	})

	if err != nil {
		log.App.WithError(err).Panic("can't watch pod event queue")
	}

	for event := range eventQueue.ResultChan() {
		log.App.Info("go new pod event!")

		if event.Type != watch.Added {
			// some other event not related to pod creation
			log.App.WithFields(logrus.Fields{"event_type": event.Type}).Info("event wasn't related to pod creation, ignoring event . . .")
			continue
		}

		podSpec, ok := event.Object.(*coreV1.Pod)
		if !ok {
			log.App.Error("unexpected event object type")
		}
		newPod := model.NewPod(string(podSpec.GetUID()))
		newPod := model.NewPod(
			string(podSpec.GetUID()),
			podSpec.Namespace,
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

		wg := new(sync.WaitGroup)
		wg.Add(2)
		go connector.BindPodToNode(wg, newPod, selectedNode)
		go connector.EmitScheduledEvent(wg, newPod, selectedNode)
		wg.Wait()
	}
}
