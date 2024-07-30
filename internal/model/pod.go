package model

import (
	"github.com/noisyboy-9/sencillo/internal/log"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
)

type Pod struct {
	ID        types.UID `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Namespace string    `json:"namespace,omitempty"`

	Cores  *resource.Quantity `json:"cores,omitempty"`
	Memory *resource.Quantity `json:"memory,omitempty"`
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
		ID:        id,
		Name:      name,
		Namespace: namespace,
		Cores:     cpu,
		Memory:    memory,
	}
}
