package edge_first

import (
	"errors"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/util"
	"math/rand"

	"github.com/noisyboy-9/random-k8s-scheduler/internal/model"
)

type SmallestFittingEdgeNodeScheduler struct{}

func (s SmallestFittingEdgeNodeScheduler) Run(pod *model.Pod, nodes []*model.Node) (node *model.Node, err error) {
	edgeNodes, cloudNodes := s.Filter(pod, nodes)
	if len(edgeNodes) == 0 && len(cloudNodes) == 0 {
		return nil, errors.New("no eligible nodes found")
	}
	return s.Schedule(edgeNodes, cloudNodes)
}

func (s SmallestFittingEdgeNodeScheduler) Filter(pod *model.Pod, nodes []*model.Node) (eligibleEdgeNodes []*model.Node, eligibleCloudNodes []*model.Node) {
	eligibleEdgeNodes = make([]*model.Node, 0)
	eligibleCloudNodes = make([]*model.Node, 0)
	for _, node := range nodes {
		if node.HasEnoughResourcesForPod(pod) && node.IsOnEdge() {
			eligibleEdgeNodes = append(eligibleEdgeNodes, node)
		}

		if node.HasEnoughResourcesForPod(pod) && !node.IsOnEdge() {
			eligibleEdgeNodes = append(eligibleCloudNodes, node)
		}
	}

	return eligibleEdgeNodes, eligibleCloudNodes
}
func (s SmallestFittingEdgeNodeScheduler) Schedule(edgeNodes []*model.Node, cloudNodes []*model.Node) (node *model.Node, err error) {
	if len(edgeNodes) != 0 {
		return s.FindSmallestEdgeNode(edgeNodes), nil
	}

	return cloudNodes[rand.Intn(len(cloudNodes))], nil
}

func (s SmallestFittingEdgeNodeScheduler) FindSmallestEdgeNode(nodes []*model.Node) *model.Node {
	smallestNode := nodes[0]
	smallestNodeResources := util.GetNodeResourceSum(smallestNode)

	for _, node := range nodes {
		resourceSum := util.GetNodeResourceSum(node)

		if smallestNodeResources.Cmp(*resourceSum) == 1 {
			smallestNode = node
			smallestNodeResources = resourceSum
		}
	}

	return smallestNode
}
