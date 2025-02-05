package model

import (
	"fmt"
	"sync"

	"k8s.io/apimachinery/pkg/api/resource"
)

type ClusterState struct {
	sync.Mutex

	Nodes         []Node
	AllocationMap map[Node][]Pod
}

func NewClusterState() *ClusterState {
	return &ClusterState{
		Nodes:         make([]Node, 0),
		AllocationMap: make(map[Node][]Pod),
	}
}

func (state *ClusterState) AddNode(node Node) {
	state.Lock()
	defer state.Unlock()

	state.Nodes = append(state.Nodes, node)
}

func (state *ClusterState) findNodeByName(nodeName string) (Node, error) {
	for _, node := range state.Nodes {
		if node.Name == nodeName {
			return node, nil
		}
	}

	return Node{}, fmt.Errorf("node: %s not found", nodeName)
}

func (state *ClusterState) Sync(podList []Pod) error {
	for _, pod := range podList {
		if pod.NodeName == "" {
			// The pod is currently unscheduled, it will be scheduled later.
			continue
		}

		node, err := state.findNodeByName(pod.NodeName)
		if err != nil {
			return err
		}

		state.AllocationMap[node] = append(state.AllocationMap[node], pod)
	}

	for node, pods := range state.AllocationMap {
		node.Cores = resource.NewQuantity(0.0, resource.DecimalSI)
		node.Memory = resource.NewQuantity(0.0, resource.BinarySI)

		for _, pod := range pods {
			node.Memory.Add(pod.Memory.DeepCopy())
			node.Cores.Add(pod.Cores.DeepCopy())
		}
	}

	return nil
}
