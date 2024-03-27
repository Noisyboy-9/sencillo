package state

import (
	"errors"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/model"
	"k8s.io/apimachinery/pkg/types"
)

type ClusterState struct {
	nodes map[types.UID]model.Node
	pods  map[types.UID]model.Pod
}

func NewClusterState() *ClusterState {
	return &ClusterState{}
}

func (s *ClusterState) AddPod(p model.Pod) {
	s.pods[p.ID()] = p
}

func (s *ClusterState) RemovePod(p model.Pod) {
	delete(s.pods, p.ID())
}
func (s *ClusterState) RemovePodByID(id types.UID) {
	delete(s.pods, id)
}

func (s *ClusterState) EditPodWithUID(id types.UID, edited model.Pod) error {
	if _, ok := s.pods[id]; !ok {
		return errors.New("pod with given uid doesn't exist")
	}
	s.pods[id] = edited
	return nil
}

func (s *ClusterState) GetPodByUID(id types.UID) (pod model.Pod, exists bool) {
	pod, exists = s.pods[id]
	return
}

func (s *ClusterState) AddNode(n model.Node) {
	s.nodes[n.ID()] = n
}

func (s *ClusterState) RemoveNode(n model.Node) {
	delete(s.nodes, n.ID())
}

func (s *ClusterState) RemoveNodeByID(id types.UID) {
	delete(s.nodes, id)
}

func (s *ClusterState) EditNodeWithUID(id types.UID, edited model.Node) error {
	if _, ok := s.nodes[id]; !ok {
		return errors.New("node with given id doesn't exist")
	}
	s.nodes[id] = edited
	return nil
}
func (s *ClusterState) GetNodeWithID(id types.UID) (node model.Node, exists bool) {
	node, exists = s.nodes[id]
	return
}
