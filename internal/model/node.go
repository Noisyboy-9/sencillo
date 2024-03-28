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
	Name     string             `json:"GetName,omitempty"`
	Memory   *resource.Quantity `json:"GetMemory,omitempty"`
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

func NewNode(id types.UID, name string, memory *resource.Quantity, cpu *resource.Quantity, isOnEdge bool) Node {
	return Node{
		ID:       id,
		Name:     name,
		Memory:   memory,
		Cores:    cpu,
		IsOnEdge: isOnEdge,
	}
}

func (node *Node) GetID() types.UID {
	return node.ID
}

func (node *Node) GetMemory() *resource.Quantity {
	return node.Memory
}

func (node *Node) GetCores() *resource.Quantity {
	return node.Cores
}
func (node *Node) GetName() string {
	return node.Name
}
func (node *Node) HasEnoughResourcesForPod(pod *Pod) bool {
	hasCpu := node.GetCores().Cmp(*pod.GetCores()) == 1
	hasMemory := node.GetMemory().Cmp(*pod.GetMemory()) == 1
	if !hasCpu {
		log.App.WithFields(logrus.Fields{
			"node_name":  node.GetName(),
			"node_cores": node.GetCores(),
			"is_on_edge": node.GetIsOnEdge(),
			"pod_cpu":    pod.GetCores(),
		}).Info("is out of GetCores")
	}

	if !hasMemory {
		log.App.WithFields(logrus.Fields{
			"node_name":   node.GetName(),
			"node_memory": node.GetMemory(),
			"is_on_edge":  node.GetIsOnEdge(),
			"pod_memory":  pod.GetMemory(),
		}).Info("is out of GetMemory")
	}

	return hasCpu && hasMemory
}

func (node *Node) GetIsOnEdge() bool {
	return node.IsOnEdge
}

func (node *Node) SetMemory(memory *resource.Quantity) {
	node.Memory = memory
}

func (node *Node) SetCores(cores *resource.Quantity) {
	node.Cores = cores
}
