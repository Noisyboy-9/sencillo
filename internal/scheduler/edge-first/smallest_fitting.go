package edge_first

import (
	"errors"

	"github.com/noisyboy-9/random-k8s-scheduler/internal/model"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/util"
)

type SmallestFittingEdgeNodeScheduler struct{}

func (s SmallestFittingEdgeNodeScheduler) Run(pod *model.Pod, nodes []*model.Node) (node *model.Node, err error) {
	eligibleNodes := s.Filter(pod, nodes)
	if len(eligibleNodes) == 0 {
		return nil, errors.New("no eligible nodes found")
	}
	return s.Schedule(eligibleNodes)
}

func (s SmallestFittingEdgeNodeScheduler) Filter(pod *model.Pod, nodes []*model.Node) []*model.Node {
	eligibleNodes := make([]*model.Node, 0)
	for _, node := range nodes {
		if node.IsOnEdge() && node.HasEnoughResourcesForPod(pod) {
			eligibleNodes = append(eligibleNodes, node)
		}
	}

	return eligibleNodes
}
func (s SmallestFittingEdgeNodeScheduler) Schedule(nodes []*model.Node) (node *model.Node, err error) {
	smallestNode := nodes[0]
	smallestNodeResources := util.GetNodeResourceSum(nodes[0])

	for _, node := range nodes {
		nodeResources := util.GetNodeResourceSum(node)

		if smallestNodeResources.Cmp(*nodeResources) == 1 {
			smallestNode = node
			smallestNodeResources = nodeResources
		}
	}

	return smallestNode, nil
}
