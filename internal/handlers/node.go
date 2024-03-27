package handlers

import (
	"github.com/noisyboy-9/random-k8s-scheduler/internal/log"
	v1 "k8s.io/api/core/v1"
)

type NodeEventHandler struct {
}

func (n NodeEventHandler) OnAdd(obj interface{}, isInInitialList bool) {
	if isInInitialList {
		node := obj.(*v1.Node)
		log.App.Info(node.Name)
	}
}

func (n NodeEventHandler) OnUpdate(oldObj, newObj interface{}) {
	log.App.Info("node update")
}

func (n NodeEventHandler) OnDelete(obj interface{}) {
	log.App.Info("node delete")
}
