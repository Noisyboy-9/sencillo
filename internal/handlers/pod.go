package handlers

import (
	"github.com/noisyboy-9/random-k8s-scheduler/internal/log"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/model"
)

type PodEventHandler struct {
	State *model.ClusterState
}

func (p PodEventHandler) OnAdd(obj interface{}, isInInitialList bool) {
	log.App.Info("pod add")
}

func (p PodEventHandler) OnUpdate(oldObj, newObj interface{}) {
	log.App.Info("pod update")
}

func (p PodEventHandler) OnDelete(obj interface{}) {
	log.App.Info("pod delete")
}
