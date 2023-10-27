package model

import (
	"github.com/noisyboy-9/random-k8s-scheduler/internal/enum"
)

type Pod struct {
	id        string
	name      string
	status    enum.PodStatus
	node      *Node
	namespace string
}

func NewPod(id, name, namespace string) *Pod {
	return &Pod{
		id:        id,
		name:      name,
		status:    enum.PodStatusPendding,
		node:      nil,
		namespace: namespace,
	}
}

func (pod *Pod) SetNode(node *Node) {
	pod.node = node
}

func (pod *Pod) SetStatus(newStatus enum.PodStatus) {
	pod.status = newStatus
}

func (pod *Pod) Id() string {
	return pod.id
}

func (pod *Pod) Status() enum.PodStatus {
	return pod.status
}

func (pod *Pod) Node() *Node {
	return pod.node
}
func (pod *Pod) Namespace() string {
	return pod.namespace
}
func (pod *Pod) Name() string {
	return pod.name
}
