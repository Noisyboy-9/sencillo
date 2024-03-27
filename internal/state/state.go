package state

import (
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

func (s *ClusterState) AddPod(pod model.Pod) {
	s.pods[pod.ID()] = pod
}

func (s *ClusterState) RemovePod(pod model.Pod) {
	delete(s.pods, pod.ID())
}

func (s *ClusterState) EditPodWithUID(id types.UID, editedPod model.Pod) {
	s.pods[id] = editedPod
}

func (s *ClusterState) GetPodByUID(id types.UID) (pod model.Pod, exists bool) {
	pod, ok := s.pods[id]
	return pod, ok
}
