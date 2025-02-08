package model

import (
	"fmt"
	"sync"

	"k8s.io/apimachinery/pkg/api/resource"
)

type ClusterState struct {
	sync.Mutex

	nodes         []Node
	allocationMap map[Node][]Pod
}

func NewClusterState() *ClusterState {
	return &ClusterState{
		nodes:         make([]Node, 0),
		allocationMap: make(map[Node][]Pod),
	}
}

func (state *ClusterState) AddNode(node Node) {
	state.nodes = append(state.nodes, node)
}

func (state *ClusterState) findNodeByName(nodeName string) (Node, error) {
	for _, node := range state.nodes {
		if node.Name == nodeName {
			return node, nil
		}
	}

	return Node{}, fmt.Errorf("node: %s not found", nodeName)
}

func (state *ClusterState) Sync(podList []Pod) ([]Node, error) {
	state.allocationMap = make(map[Node][]Pod)

	for _, pod := range podList {
		if pod.NodeName == "" {
			// The pod is currently unscheduled, it will be scheduled later.
			continue
		}

		node, err := state.findNodeByName(pod.NodeName)
		if err != nil {
			return nil, err
		}

		state.allocationMap[node] = append(state.allocationMap[node], pod)
	}

	syncedNodes := make([]Node, 0, len(state.nodes))
	for _, node := range state.nodes {
		pods, ok := state.allocationMap[node]
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
