package model

import (
	"github.com/noisyboy-9/random-k8s-scheduler/internal/log"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
)

type Pod struct {
	id        types.UID `json:"id,omitempty"`
	name      string    `json:"name,omitempty"`
	namespace string    `json:"namespace,omitempty"`

	cpu    *resource.Quantity `json:"cpu,omitempty"`
	memory *resource.Quantity `json:"memory,omitempty"`
}

func (pod *Pod) String() string {
	j, err := json.Marshal(pod)
	if err != nil {
		log.App.WithError(err).Panic("can't marshal node")
	}
	return string(j)
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
