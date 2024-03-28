package model

import (
	"errors"
	"k8s.io/apimachinery/pkg/types"
	"sync"
)

type ClusterState struct {
	mux           *sync.RWMutex
	nodes         map[types.UID]Node
	isNodesSynced bool

	pods         map[types.UID]Pod
	isPodsSynced bool
}

func NewClusterState() *ClusterState {
	return &ClusterState{
		mux: new(sync.RWMutex),
	}
}

func (s *ClusterState) IsPodsSynced() bool {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return s.isPodsSynced
}

func (s *ClusterState) IsNodesSynced() bool {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return s.isNodesSynced
}

func (s *ClusterState) SetIsPodsSynced(isPodsSynced bool) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.isPodsSynced = isPodsSynced
}

func (s *ClusterState) SetIsNodesSynced(isNodesSynced bool) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.isNodesSynced = isNodesSynced
}

func (s *ClusterState) AddPod(p Pod) {
	s.pods[p.ID] = p
}

func (s *ClusterState) RemovePod(p Pod) {
	delete(s.pods, p.ID)
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
	s.nodes[n.ID] = n
}

func (s *ClusterState) RemoveNode(n Node) {
	delete(s.nodes, n.ID)
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
