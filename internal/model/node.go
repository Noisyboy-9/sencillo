package model

import (
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"

	"github.com/noisyboy-9/sencillo/internal/log"
)

type Node struct {
	ID                types.UID         `json:"id"`
	Name              string            `json:"name"`
	AllocatableCores  resource.Quantity `json:"allocatable_cpu"`
	AllocatableMemory resource.Quantity `json:"allocatable_memory"`
	RemainingCores    resource.Quantity `json:"cores"`
	RemainingMemory   resource.Quantity `json:"memory"`
	IsOnEdge          bool              `json:"is_on_edge"`
	IsMaster          bool              `json:"is_master"`
}

func NewNode(id types.UID, name string, allocatableCores resource.Quantity, allocatableMemory resource.Quantity, isOnEdge bool, isMaster bool) Node {
	return Node{
		ID:   id,
		Name: name,

		AllocatableCores:  allocatableCores,
		AllocatableMemory: allocatableMemory,

		RemainingCores:  allocatableCores,
		RemainingMemory: allocatableMemory,

		IsOnEdge: isOnEdge,
		IsMaster: isMaster,
	}
}

func (node Node) HasEnoughResourcesForPod(pod Pod) bool {
	node.RemainingCores.Sub(resource.NewQuantity(1, resource.DecimalSI).DeepCopy())
	node.RemainingMemory.Sub(resource.NewQuantity(-1024, resource.BinarySI).DeepCopy())

	hasCpu := node.RemainingCores.Cmp(pod.Cores.DeepCopy()) == 1
	hasMemory := node.RemainingMemory.Cmp(pod.Memory.DeepCopy()) == 1
	if !hasCpu {
		log.App.WithFields(logrus.Fields{
			"node_name":  node.Name,
			"node_cores": node.RemainingCores.String(),
			"is_on_edge": node.IsOnEdge,
			"pod_cpu":    pod.Cores,
		}).Info("is out of cores")
	}

	if !hasMemory {
		log.App.WithFields(logrus.Fields{
			"node_name":   node.Name,
			"node_memory": node.RemainingMemory.String(),
			"is_on_edge":  node.IsOnEdge,
			"pod_memory":  pod.Memory,
		}).Info("is out of memory")
	}

	return hasCpu && hasMemory
}
