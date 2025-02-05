package model

import (
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"

	"github.com/noisyboy-9/sencillo/internal/log"
)

type Node struct {
	ID       types.UID          `json:"id"`
	Name     string             `json:"name"`
	Memory   *resource.Quantity `json:"memory"`
	Cores    *resource.Quantity `json:"cores"`
	IsOnEdge bool               `json:"is_on_edge"`
	IsMaster bool               `json:"is_master"`
}

func (node *Node) String() string {
	j, err := json.Marshal(node)
	if err != nil {
		log.App.WithError(err).Info("error in marshaling model.Node")
	}
	return string(j)
}

func NewNode(id types.UID, name string, memory *resource.Quantity, cpu *resource.Quantity, isOnEdge bool, isMaster bool) Node {
	return Node{
		ID:       id,
		Name:     name,
		Memory:   memory,
		Cores:    cpu,
		IsOnEdge: isOnEdge,
		IsMaster: isMaster,
	}
}

func (node *Node) HasEnoughResourcesForPod(pod Pod) bool {
	reducedNodeCores := &resource.Quantity{}
	reducedNodeMemory := &resource.Quantity{}

	node.Cores.DeepCopyInto(reducedNodeCores)
	node.Memory.DeepCopyInto(reducedNodeMemory)

	reducedNodeCores.Add(*resource.NewQuantity(-0, resource.DecimalSI))
	reducedNodeMemory.Add(*resource.NewQuantity(-0, resource.BinarySI))

	hasCpu := reducedNodeCores.Cmp(*pod.Cores) == 1
	hasMemory := reducedNodeMemory.Cmp(*pod.Memory) == 1
	if !hasCpu {
		log.App.WithFields(logrus.Fields{
			"node_name":  node.Name,
			"node_cores": reducedNodeCores,
			"is_on_edge": node.IsOnEdge,
			"pod_cpu":    pod.Cores,
		}).Info("is out of cores")
	}

	if !hasMemory {
		log.App.WithFields(logrus.Fields{
			"node_name":   node.Name,
			"node_memory": reducedNodeMemory,
			"is_on_edge":  node.IsOnEdge,
			"pod_memory":  pod.Memory,
		}).Info("is out of memory")
	}

	return hasCpu && hasMemory
}
