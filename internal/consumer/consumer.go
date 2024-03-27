package consumer

import (
	"context"
	"fmt"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/scheduler"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/util"

	"github.com/noisyboy-9/random-k8s-scheduler/internal/config"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/connector"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/log"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/model"
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
		newNode := model.NewNode(
			string(node.GetUID()),
			node.Name,
			node.Status.Allocatable.Memory(),
			node.Status.Allocatable.Cpu(),
		)
		consumer.nodes = append(consumer.nodes, newNode)

		log.App.WithFields(logrus.Fields{
			"node_name":          newNode.Name(),
			"is_on_edge":         newNode.IsOnEdge(),
			"cores":              newNode.Cores(),
			"memory":             newNode.Memory(),
			"current_node_count": len(consumer.nodes),
		}).Info("node has been added")
	}
	return consumer
}

func (c *Consumer) Consume() {
	go func() {

		nodeWatcher, err := connector.ClusterConnection.Client().CoreV1().Nodes().Watch(context.Background(), metaV1.ListOptions{})
		if err != nil {
			log.App.WithError(err).Panic("can't watch node event queue")
		}

		for event := range nodeWatcher.ResultChan() {
			if event.Type != watch.Modified {
				continue
			}

			modifiedNodeSpec, ok := event.Object.(*coreV1.Node)
			if !ok {
				log.App.Panic("unexpected node event object type")
			}

			i, err := util.FindNodeIndexByName(c.nodes, modifiedNodeSpec.Name)
			if err != nil {
				log.App.WithError(err).Error("error in getting node by name")
			}

			c.nodes[i] = model.NewNode(
				string(modifiedNodeSpec.GetUID()),
				modifiedNodeSpec.Name,
				modifiedNodeSpec.Status.Allocatable.Memory(),
				modifiedNodeSpec.Status.Allocatable.Cpu(),
			)
		}
	}()

	eventQueue, err := connector.ClusterConnection.Client().CoreV1().Pods(config.Scheduler.Namespace).Watch(
		context.Background(),
		metaV1.ListOptions{
			FieldSelector: fmt.Sprintf("spec.schedulerName=%s", config.Scheduler.Name),
		},
	)
	if err != nil {
		log.App.WithError(err).Panic("can't watch pod event queue")
	}

	for event := range eventQueue.ResultChan() {
		switch event.Type {
		case watch.Added:
			c.handlePodAdd(event)
		case watch.Deleted:
			c.handlePodDelete(event)
		case watch.Modified:
			c.handlePodModified(event)
		default:
			log.App.WithField("type", event.Type).Info("unrelated event to scheduling received")
			continue
		}
	}
}

func (c *Consumer) handlePodAdd(event watch.Event) {
	podSpec, ok := event.Object.(*coreV1.Pod)
	if !ok {
		log.App.Error("unexpected event pod event object type")
	}

	newPod := model.NewPod(
		string(podSpec.GetUID()),
		podSpec.Name,
		podSpec.Namespace,
		util.RequiredCpuSum(podSpec.Spec.Containers),
		util.RequiredMemorySum(podSpec.Spec.Containers),
	)

	c.allocatePodToNode(newPod)
}

func (c *Consumer) allocatePodToNode(newPod *model.Pod) (*model.Node, bool) {
	c.pods = append(c.pods, newPod)

	log.App.WithFields(logrus.Fields{
		"id":     newPod.Id(),
		"name":   newPod.Name(),
		"cpu":    newPod.Cores(),
		"memory": newPod.Memory(),
	}).Info("New pod created")

	selectedNode, err := scheduler.S.Run(newPod, c.nodes)
	if err != nil {
		log.App.WithError(err).WithFields(logrus.Fields{"pod": newPod.Id()}).Error("error in finding node for pod")
		return nil, true
	}
	log.App.WithFields(logrus.Fields{
		"pod":           newPod.Name(),
		"selected_node": selectedNode.Name(),
	}).Info("node selected for pod")

	if err := connector.ClusterConnection.BindPodToNode(newPod, selectedNode); err != nil {
		log.App.WithError(err).Error("pod bind error")
		return nil, true
	}
	if err := connector.ClusterConnection.EmitScheduledEvent(newPod, selectedNode); err != nil {
		log.App.WithError(err).Error("scheduled event emit error")
		return nil, true
	}
	return selectedNode, false
}

func (c *Consumer) handlePodDelete(event watch.Event) {
	podEvent, ok := event.Object.(*coreV1.Pod)
	if !ok {
		log.App.Error("unexpected event object type")
	}

	pod, err := util.FindPodByName(c.pods, podEvent.Name)
	if err != nil {
		log.App.WithError(err).Error("error in finding pod by name")
		return
	}

	node, err := util.FindNodeByName(c.nodes, podEvent.Spec.NodeName)
	if err != nil {
		log.App.WithError(err).Error("error in finding node by name")
		return
	}

	c.deAllocatePodFromNode(pod, node)
}

func (c *Consumer) deAllocatePodFromNode(pod *model.Pod, node *model.Node) {
	c.removePodFromListOfPods(pod)
}

func (c *Consumer) removePodFromListOfPods(pod *model.Pod) {
	for i, p := range c.pods {
		if p.Name() == pod.Name() {
			c.pods = append(c.pods[:i], c.pods[i+1:]...)
			return
		}
	}
}

func (c *Consumer) handlePodModified(event watch.Event) {
	podEvent, ok := event.Object.(*coreV1.Pod)
	if !ok {
		log.App.Error("unexpected event object type")
	}

	pod := model.NewPod(
		string(podEvent.GetUID()),
		podEvent.Name,
		podEvent.Namespace,
		util.RequiredCpuSum(podEvent.Spec.Containers),
		util.RequiredMemorySum(podEvent.Spec.Containers),
	)

	node, err := util.FindNodeByName(c.nodes, podEvent.Spec.NodeName)
	if err != nil {
		log.App.WithError(err).Error("error in finding node by name")
		return
	}

	c.deAllocatePodFromNode(pod, node)
}
