package edge_first

import (
	"errors"

	"github.com/noisyboy-9/random-k8s-scheduler/internal/model"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/util"
)

type BiggestFittingEdgeNodeScheduler struct{}

func (b BiggestFittingEdgeNodeScheduler) Run(pod *model.Pod, nodes []*model.Node) (node *model.Node, err error) {
	eligibleNodes := b.Filter(pod, nodes)
	if len(eligibleNodes) == 0 {
		return nil, errors.New("no eligible nodes found")
	}
	return b.Schedule(eligibleNodes)
}

func (b BiggestFittingEdgeNodeScheduler) Filter(pod *model.Pod, nodes []*model.Node) []*model.Node {
	eligibleNodes := make([]*model.Node, 0)
	for _, node := range nodes {
		if node.IsOnEdge() && node.HasEnoughResourcesForPod(pod) {
			eligibleNodes = append(eligibleNodes, node)
		}
	}
	return eligibleNodes
}

func (b BiggestFittingEdgeNodeScheduler) Schedule(nodes []*model.Node) (node *model.Node, err error) {
	biggestNode := nodes[0]
	biggestNodeResources, err := util.GetNodeResourceSum(nodes[0])
	if err != nil {
		return nil, err
	}

	for _, node := range nodes {
		nodeResources, err := util.GetNodeResourceSum(node)
		if err != nil {
			return nil, err
		}

		if biggestNodeResources.Cmp(*nodeResources) == -1 {
			biggestNode = node
			biggestNodeResources = nodeResources
		}
	}

	return biggestNode, nil
}
