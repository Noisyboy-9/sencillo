package model

import (
	"errors"
	"k8s.io/apimachinery/pkg/types"
)

type ClusterState struct {
	nodes map[types.UID]Node
	pods  map[types.UID]Pod
}

func NewClusterState() *ClusterState {
	return &ClusterState{}
}

func (s *ClusterState) AddPod(p Pod) {
	s.pods[p.ID()] = p
}

func (s *ClusterState) RemovePod(p Pod) {
	delete(s.pods, p.ID())
}
func (s *ClusterState) RemovePodByID(id types.UID) {
	delete(s.pods, id)
}

func (s *ClusterState) EditPodWithUID(id types.UID, edited Pod) error {
	if _, ok := s.pods[id]; !ok {
		return errors.New("pod with given uid doesn't exist")
	}
	s.pods[id] = edited
	return nil
}

func (s *ClusterState) GetPodByUID(id types.UID) (pod Pod, exists bool) {
	pod, exists = s.pods[id]
	return
}

func (s *ClusterState) AddNode(n Node) {
	s.nodes[n.ID()] = n
}

func (s *ClusterState) RemoveNode(n Node) {
	delete(s.nodes, n.ID())
}

func (s *ClusterState) RemoveNodeByID(id types.UID) {
	delete(s.nodes, id)
}

func (s *ClusterState) EditNodeWithUID(id types.UID, edited Node) error {
	if _, ok := s.nodes[id]; !ok {
		return errors.New("node with given id doesn't exist")
	}
	s.nodes[id] = edited
	return nil
}
func (s *ClusterState) GetNodeWithID(id types.UID) (node Node, exists bool) {
	node, exists = s.nodes[id]
	return
}

func (s *ClusterState) Nodes() []Node {
	nodes := make([]Node, 0, len(s.nodes))
	i := 0
	for _, n := range s.nodes {
		nodes[i] = n
		i++
	}

	return nodes
}

func (s *ClusterState) Pods() []Pod {
	pods := make([]Pod, 0, len(s.pods))
	i := 0
	for _, p := range s.pods {
		pods[i] = p
		i++
	}

	return pods
}
