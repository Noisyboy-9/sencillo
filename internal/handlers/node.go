package handlers

import (
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"

	"github.com/noisyboy-9/sencillo/internal/log"
	"github.com/noisyboy-9/sencillo/internal/model"
	"github.com/noisyboy-9/sencillo/internal/util"
)

type NodeEventHandler struct {
	State *model.ClusterState
}

func (n NodeEventHandler) OnAdd(obj interface{}, isInInitialList bool) {
	kubeNode, ok := obj.(*v1.Node)
	if !ok {
		log.App.Panic("unexpected event object type")
		return
	}

	node := model.NewNode(
		kubeNode.GetUID(),
		kubeNode.GetName(),
		kubeNode.Status.Allocatable.Cpu(),
		kubeNode.Status.Allocatable.Memory(),
		util.IsNodeOnEdge(kubeNode),
		util.IsMasterNode(kubeNode),
	)

	n.State.Lock()
	n.State.AddNode(node)
	n.State.Unlock()

	log.App.WithFields(logrus.Fields{
		"node":         node,
		"is_init_list": isInInitialList,
	}).Info("new node has been added")
}

func (n NodeEventHandler) OnUpdate(oldObj interface{}, newObj interface{}) {
	return
}

func (n NodeEventHandler) OnDelete(obj interface{}) {
	return
}
