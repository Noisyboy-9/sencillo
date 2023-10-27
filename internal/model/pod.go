package model

import (
	"github.com/noisyboy-9/random-k8s-scheduler/internal/enum"
)

type Pod struct {
	id         int
	status     enum.PodStatus
	node       *Node
	deployment *Deployment
}

func NewPod(id int, deployment *Deployment) *Pod {
	return &Pod{
		id:         id,
		status:     enum.PodStatusPendding,
		node:       nil,
		deployment: deployment,
	}
}

func (pod *Pod) SetNode(node *Node) {
	pod.node = node
}

func (pod *Pod) SetStatus(newStatus enum.PodStatus) {
	pod.status = newStatus
}

func (pod *Pod) Id() int {
	return pod.id
}

func (pod *Pod) Status() enum.PodStatus {
	return pod.status
}

func (pod *Pod) Node() *Node {
	return pod.node
}

func (pod *Pod) Deployment() *Deployment {
	return pod.deployment
}
