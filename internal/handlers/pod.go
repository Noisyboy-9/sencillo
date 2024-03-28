package handlers

import "C"
import (
	"github.com/noisyboy-9/random-k8s-scheduler/internal/config"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/connector"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/log"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/model"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/scheduler"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/util"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

type PodEventHandler struct {
	State        *model.ClusterState
	PodScheduler scheduler.Scheduler
	NodeInformer cache.SharedIndexInformer
}

func (p PodEventHandler) OnAdd(obj interface{}, _ bool) {
	podKubernetesObj, ok := obj.(*v1.Pod)
	if !ok {
		log.App.Panic("unexpected pod event type")
	}

	if !p.isEventForThisScheduler(podKubernetesObj) {
		return
	}

	if !p.State.IsNodesSynced() {
		return
	}

	log.App.WithField("name", podKubernetesObj.GetName()).Infof("got pod add event")
	pod := model.NewPod(
		podKubernetesObj.GetUID(),
		podKubernetesObj.GetName(),
		podKubernetesObj.GetNamespace(),
		util.RequiredCpuSum(podKubernetesObj.Spec.Containers),
		util.RequiredMemorySum(podKubernetesObj.Spec.Containers),
	)
	p.State.AddPod(pod)

	selectedNode, err := p.PodScheduler.Run(pod, p.State.Nodes())
	if err != nil {
		log.App.WithError(err).WithField("pod", pod).Error("error in selecting node for pod")
		return
	}

	if err = connector.C.BindPodToNode(pod, selectedNode); err != nil {
		log.App.WithError(err).Error("error in binding pod to node")
	}

	p.State.SaveSelectedNodeForPod(selectedNode, pod)
	p.State.AllocateResources(selectedNode, pod)

	if err := connector.C.EmitScheduledEvent(pod, selectedNode); err != nil {
		log.App.WithError(err).WithFields(logrus.Fields{
			"selected_node": selectedNode,
			"pod":           pod,
		}).Error("error in emitting pod scheduled event")
	}

	log.App.WithFields(logrus.Fields{
		"pod": pod,
	}).Info("pod creation handled")
}

func (p PodEventHandler) isEventForThisScheduler(podKubernetesObj *v1.Pod) bool {
	return podKubernetesObj.GetNamespace() == config.Scheduler.Namespace &&
		podKubernetesObj.Spec.SchedulerName == config.Scheduler.Name
}

func (p PodEventHandler) OnUpdate(oldObj interface{}, newObj interface{}) {
	oldPodKubernetesObj, ok := oldObj.(*v1.Pod)
	if !ok {
		log.App.Panic("unexpected pod event type")
	}

	newPodKubernetesObj, ok := newObj.(*v1.Pod)
	if !ok {
		log.App.Panic("unexpected pod event type")
	}

	if !p.isEventForThisScheduler(oldPodKubernetesObj) || !p.isEventForThisScheduler(newPodKubernetesObj) {
		return
	}

	if newPodKubernetesObj.Status.Phase == v1.PodPending && oldPodKubernetesObj.Status.Phase == v1.PodPending {
		p.OnAdd(newPodKubernetesObj, false)
		return
	}

	_, exists := p.State.GetPodByUID(oldPodKubernetesObj.GetUID())
	if !exists {
		newPod := model.NewPod(
			newPodKubernetesObj.GetUID(),
			newPodKubernetesObj.GetName(),
			newPodKubernetesObj.GetNamespace(),
			util.RequiredCpuSum(newPodKubernetesObj.Spec.Containers),
			util.RequiredMemorySum(newPodKubernetesObj.Spec.Containers),
		)
		nodeName := newPodKubernetesObj.Spec.NodeName
		node, _ := p.State.FindNodeByName(nodeName)

		p.State.AddPod(newPod)
		p.State.SaveSelectedNodeForPod(node, newPod)
	}

	log.App.WithFields(logrus.Fields{
		"old_pod_name":   oldPodKubernetesObj.GetName(),
		"new_pod_name":   newPodKubernetesObj.GetName(),
		"old_pod_status": oldPodKubernetesObj.Status.Phase,
		"new_pod_status": newPodKubernetesObj.Status.Phase,
	}).Infof("got pod update event")
}

func (p PodEventHandler) OnDelete(obj interface{}) {
	deletedPodKubernetesObj, ok := obj.(*v1.Pod)
	if !ok {
		log.App.Panic("unexpected pod event type")
	}

	if !p.isEventForThisScheduler(deletedPodKubernetesObj) {
		return
	}

	if !p.State.IsNodesSynced() {
		return
	}

	deletedPod, exists := p.State.GetPodByUID(deletedPodKubernetesObj.GetUID())
	if !exists {
		log.App.WithField("pod", deletedPod).Error("trying to delete pod that doesn't exist")
	}

	node, exist := p.State.GetSelectedNodeForPod(deletedPod)
	if !exist {
		log.App.WithField("pod", deletedPod).Panic("trying to deleted pod which is not owned")
	}
	p.State.DeAllocateResources(node, deletedPod)
	p.State.RemovePod(deletedPod)

	log.App.WithField("node", deletedPod).Info("handled pod delete")
}
