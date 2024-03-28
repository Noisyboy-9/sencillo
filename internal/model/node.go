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

func (node *Node) HasEnoughResourcesForPod(pod *Pod) bool {
	hasCpu := node.Cores.Cmp(*pod.GetCores()) == 1
	hasMemory := node.Memory.Cmp(*pod.GetMemory()) == 1
	if !hasCpu {
		log.App.WithFields(logrus.Fields{
			"node_name":  node.Name,
			"node_cores": node.Cores,
			"is_on_edge": node.IsOnEdge,
			"pod_cpu":    pod.GetCores(),
		}).Info("is out of GetCores")
	}

	if !hasMemory {
		log.App.WithFields(logrus.Fields{
			"node_name":   node.Name,
			"node_memory": node.Memory,
			"is_on_edge":  node.IsOnEdge,
			"pod_memory":  pod.Memory,
		}).Info("is out of GetMemory")
	}

	return hasCpu && hasMemory
}
func (node *Node) SetMemory(memory *resource.Quantity) {
	node.Memory = memory
}

func (node *Node) SetCores(cores *resource.Quantity) {
	node.Cores = cores
}
