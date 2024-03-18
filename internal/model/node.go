package model

import (
	"slices"

	"github.com/noisyboy-9/random-k8s-scheduler/internal/app"
	"k8s.io/apimachinery/pkg/api/resource"
)

type Node struct {
	id       string
	name     string
	memory   *resource.Quantity
	cores    *resource.Quantity
	isOnEdge bool
}

func NewNode(id string, name string, memory *resource.Quantity, cpu *resource.Quantity) *Node {
	return &Node{
		id:       id,
		name:     name,
		memory:   memory,
		cores:    cpu,
		isOnEdge: checkIfOnEdge(name),
	}
}

func checkIfOnEdge(name string) bool {
	return slices.Contains(app.EdgeNodeList, name)
}

func (node *Node) Id() string {
	return node.id
}

func (node *Node) Memory() *resource.Quantity {
	return node.memory
}

func (node *Node) Cores() *resource.Quantity {
	return node.cores
}
func (node *Node) Name() string {
	return node.name
}
func (node *Node) ReduceAllocatableMemory(q *resource.Quantity) {
	node.memory.Sub(*q)
}

func (node *Node) ReduceAllocatableCpu(q *resource.Quantity) {
	node.cores.Sub(*q)
}

func (node *Node) HasEnoughResourcesForPod(pod *Pod) bool {
	return node.Cores().Cmp(*pod.Cpu()) == 1 &&
		node.Memory().Cmp(*pod.Memory()) == 1
}

func (node *Node) IsOnEdge() bool {
	return node.isOnEdge
}
