package handlers

import (
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"

	"github.com/noisyboy-9/sencillo/internal/config"
	"github.com/noisyboy-9/sencillo/internal/connector"
	"github.com/noisyboy-9/sencillo/internal/log"
	"github.com/noisyboy-9/sencillo/internal/model"
	"github.com/noisyboy-9/sencillo/internal/scheduler"
	"github.com/noisyboy-9/sencillo/internal/util"
)

type PodEventHandler struct {
	State        *model.ClusterState
	PodScheduler scheduler.Scheduler
}

func (p PodEventHandler) OnAdd(obj interface{}, _ bool) {
	podKubernetesObj, ok := obj.(*v1.Pod)
	if !ok {
		log.App.Panic("unexpected pod event type")
	}

	if !p.isEventForThisScheduler(podKubernetesObj) {
		return
	}

	log.App.WithField("name", podKubernetesObj.GetName()).Infof("got pod add event")
	pod := model.NewPod(
		podKubernetesObj.GetUID(),
		podKubernetesObj.GetName(),
		podKubernetesObj.GetNamespace(),
		"",
		util.RequiredCpuSum(podKubernetesObj.Spec.Containers),
		util.RequiredMemorySum(podKubernetesObj.Spec.Containers),
	)

	p.State.Lock()
	defer p.State.Unlock()

	pods, err := connector.C.GetAllPods()
	if err != nil {
		log.App.WithError(err).Error("failed to get pods")
		return
	}

	syncedNodes, err := p.State.Sync(pods)
	if err != nil {
		log.App.WithError(err).Error("failed to sync pods")
		return
	}

	selectedNode, err := p.PodScheduler.Run(pod, syncedNodes)
	if err != nil {
		log.App.WithError(err).WithField("pod", pod).Error("error in selecting node for pod")
		return
	}

	if err = connector.C.BindPodToNode(pod, selectedNode); err != nil {
		log.App.WithError(err).Error("error in binding pod to node")
		return
	}

	if err := connector.C.EmitScheduledEvent(pod, selectedNode); err != nil {
		log.App.WithError(err).WithFields(logrus.Fields{
			"selected_node": selectedNode,
			"pod":           pod,
		}).Error("error in emitting pod scheduled event")
		return
	}

	log.App.WithFields(logrus.Fields{"pod": pod, "selected_node": selectedNode}).Info("pod creation handled")
}

func (p PodEventHandler) isEventForThisScheduler(podKubernetesObj *v1.Pod) bool {
	return podKubernetesObj.GetNamespace() == config.Scheduler.Namespace &&
		podKubernetesObj.Spec.SchedulerName == config.Scheduler.Name
}

func (p PodEventHandler) OnUpdate(oldObj interface{}, newObj interface{}) {
	return
}

func (p PodEventHandler) OnDelete(obj interface{}) {
	return
}
