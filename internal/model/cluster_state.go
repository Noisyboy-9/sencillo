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

func (state *ClusterState) Sync(podList []Pod) ([]Node, error) {
	for _, pod := range podList {
		if pod.NodeName == "" {
			// The pod is currently unscheduled, it will be scheduled later.
			continue
		}

		node, err := state.findNodeByName(pod.NodeName)
		if err != nil {
			return nil, err
		}

		state.AllocationMap[node] = append(state.AllocationMap[node], pod)
	}

	syncedNodes := make([]Node, 0, len(state.Nodes))
	for _, node := range state.Nodes {
		pods, ok := state.AllocationMap[node]
		if !ok {
			//no pod is scheduled on the node
			syncedNodes = append(syncedNodes, node)
			continue
		}

		coresSum := resource.NewQuantity(0, resource.DecimalSI)
		memorySum := resource.NewQuantity(0, resource.BinarySI)

		for _, pod := range pods {
			coresSum.Add(pod.Cores.DeepCopy())
			memorySum.Add(pod.Memory.DeepCopy())
		}

		node.RemainingCores.Sub(coresSum.DeepCopy())
		node.RemainingMemory.Sub(memorySum.DeepCopy())

		syncedNodes = append(syncedNodes, node)
	}

	return syncedNodes, nil
}
