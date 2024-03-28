package model

import (
	"github.com/noisyboy-9/random-k8s-scheduler/internal/log"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
)

type Node struct {
	ID       types.UID          `json:"id,omitempty"`
	Name     string             `json:"name,omitempty"`
	Memory   *resource.Quantity `json:"memory,omitempty"`
	Cores    *resource.Quantity `json:"cores,omitempty"`
	IsOnEdge bool               `json:"is_on_edge,omitempty"`
}

func (node *Node) String() string {
	j, err := json.Marshal(node)
	if err != nil {
		log.App.WithError(err).Info("error in marshaling model.Node")
	}
	return string(j)
}

func NewNode(id types.UID, name string, memory *resource.Quantity, cpu *resource.Quantity, isOnEdge bool) *Node {
	return &Node{
		ID:       id,
		Name:     name,
		Memory:   memory,
		Cores:    cpu,
		IsOnEdge: isOnEdge,
	}
}

func (node *Node) HasEnoughResourcesForPod(pod *Pod) bool {
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
func (node *Node) SetMemory(memory *resource.Quantity) {
	node.Memory = memory
}

func (node *Node) SetCores(cores *resource.Quantity) {
	node.Cores = cores
}

func (node *Node) allocateMemory(memory resource.Quantity) {
	node.Memory.Sub(memory)
}

func (node *Node) allocateCpu(cpu resource.Quantity) {
	node.Cores.Sub(cpu)
}

func (node *Node) deallocateMemory(memory resource.Quantity) {
	node.Memory.Add(memory)
}

func (node *Node) deAllocateCpu(cpu resource.Quantity) {
	node.Cores.Add(cpu)
}
