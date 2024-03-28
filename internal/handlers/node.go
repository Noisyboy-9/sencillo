package handlers

import (
	"github.com/noisyboy-9/random-k8s-scheduler/internal/consumer"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/log"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/model"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/util"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

type NodeEventHandler struct {
}

func (n NodeEventHandler) OnAdd(obj interface{}, isInInitialList bool) {
	if isInInitialList {
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
		)

		consumer.C.State.AddNode(node)

		log.App.WithFields(logrus.Fields{
			"node":         node,
			"is_init_list": isInInitialList,
		}).Info("new node has been added")
	}

	if consumer.C.State.IsNodesSynced() && consumer.C.State.IsPodsSynced() {
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
		)

		consumer.C.State.AddNode(node)

		log.App.WithField("node", node).Info("node")
		consumer.C.State.AddNode(node)

		log.App.WithFields(logrus.Fields{
			"node":         node,
			"is_init_list": isInInitialList,
		}).Info("new node has been added")
	}
}

func (n NodeEventHandler) OnUpdate(oldObj interface{}, newObj interface{}) {
	oldNodeKubernetesObj, ok := oldObj.(*v1.Node)
	if !ok {
		log.App.Panic("unexpected event object type")
		return
	}
	newNodeKubernetesObj, ok := newObj.(*v1.Node)
	if !ok {
		log.App.Panic("unexpected event object type")
		return
	}

	if consumer.C.State.IsNodesSynced() && consumer.C.State.IsPodsSynced() {
		oldNode := model.NewNode(
			oldNodeKubernetesObj.GetUID(),
			oldNodeKubernetesObj.GetName(),
			oldNodeKubernetesObj.Status.Allocatable.Memory(),
			oldNodeKubernetesObj.Status.Allocatable.Cpu(),
			util.IsNodeOnEdge(oldNodeKubernetesObj),
		)

		newNode := model.NewNode(
			newNodeKubernetesObj.GetUID(),
			newNodeKubernetesObj.GetName(),
			newNodeKubernetesObj.Status.Allocatable.Memory(),
			newNodeKubernetesObj.Status.Allocatable.Cpu(),
			util.IsNodeOnEdge(newNodeKubernetesObj),
		)

		err := consumer.C.State.EditNodeWithUID(oldNode.ID(), newNode)
		if err != nil {
			log.App.WithError(err).Error("error in updating with UID")
		}

		log.App.WithFields(logrus.Fields{
			"old_node": oldNode,
			"new_node": newNode,
		}).Info("updated node status")
	}
}

func (n NodeEventHandler) OnDelete(obj interface{}) {
	deletedNodeKubernetesObject, ok := obj.(*v1.Node)
	if !ok {
		log.App.Panic("unexpected event object type")
		return
	}

	if consumer.C.State.IsNodesSynced() && consumer.C.State.IsPodsSynced() {
		node := model.NewNode(
			deletedNodeKubernetesObject.GetUID(),
			deletedNodeKubernetesObject.GetName(),
			deletedNodeKubernetesObject.Status.Allocatable.Memory(),
			deletedNodeKubernetesObject.Status.Allocatable.Cpu(),
			util.IsNodeOnEdge(deletedNodeKubernetesObject),
		)

		consumer.C.State.RemoveNode(node)
		log.App.WithField("node", node).Info("deleted node")
	}
}
