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
	nodeKubernetesObject, ok := obj.(*v1.Node)
	if !ok {
		log.App.Panic("unexpected event object type")
		return
	}

	node := model.NewNode(
		nodeKubernetesObject.GetUID(),
		nodeKubernetesObject.GetName(),
		nodeKubernetesObject.Status.Allocatable.Memory(),
		nodeKubernetesObject.Status.Allocatable.Cpu(),
		util.IsNodeOnEdge(nodeKubernetesObject),
		util.IsMasterNode(nodeKubernetesObject),
	)

	n.State.AddNode(node)

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
