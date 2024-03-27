package model

import (
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
)

type Pod struct {
	id        types.UID
	name      string
	namespace string

	cpu    *resource.Quantity
	memory *resource.Quantity
}

func NewPod(id types.UID, name string, namespace string, cpu *resource.Quantity, memory *resource.Quantity) *Pod {
	return &Pod{
		id:        id,
		name:      name,
		namespace: namespace,
		cpu:       cpu,
		memory:    memory,
	}
}

func (pod *Pod) ID() types.UID {
	return pod.id
}

func (pod *Pod) Namespace() string {
	return pod.namespace
}

func (pod *Pod) Name() string {
	return pod.name
}

func (pod *Pod) Memory() *resource.Quantity {
	return pod.memory
}

func (pod *Pod) Cores() *resource.Quantity {
	return pod.cpu
}
