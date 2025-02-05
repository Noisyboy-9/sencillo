package model

import (
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
)

type Pod struct {
	ID        types.UID          `json:"id,omitempty"`
	Name      string             `json:"name,omitempty"`
	Namespace string             `json:"namespace,omitempty"`
	NodeName  string             `json:"nodeName,omitempty"`
	Cores     *resource.Quantity `json:"cores,omitempty"`
	Memory    *resource.Quantity `json:"memory,omitempty"`
}

func NewPod(ID types.UID, name string, namespace string, nodeName string, cores *resource.Quantity, memory *resource.Quantity) Pod {
	return Pod{ID: ID, Name: name, Namespace: namespace, NodeName: nodeName, Cores: cores, Memory: memory}
}
