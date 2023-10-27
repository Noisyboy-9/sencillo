package model

import (
	"github.com/noisyboy-9/random-k8s-scheduler/internal/enum"
)

type Pod struct {
	id        string
	status    enum.PodStatus
	node      *Node
	namespace string
}

func NewPod(id, namespace string) *Pod {
	return &Pod{
		id:        id,
		namespace: namespace,
		status:    enum.PodStatusPendding,
		node:      nil,
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
