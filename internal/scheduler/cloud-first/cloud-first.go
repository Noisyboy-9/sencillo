package cloud_first

import (
	"errors"
	"math/rand"

	"github.com/noisyboy-9/random-k8s-scheduler/internal/model"
)

type CloudFirstScheduler struct{}

func (c CloudFirstScheduler) Run(pod *model.Pod, nodes []*model.Node) (node *model.Node, err error) {
	eligibleNodes := c.Filter(pod, nodes)
	if len(eligibleNodes) == 0 {
		return nil, errors.New("no eligible nodes found")
	}
	return c.Schedule(eligibleNodes)
}
func (c CloudFirstScheduler) Filter(pod *model.Pod, nodes []*model.Node) []*model.Node {
	eligibleNodes := make([]*model.Node, 0)
	for _, node := range nodes {
		if node.IsOnEdge() && node.HasEnoughResourcesForPod(pod) {
			eligibleNodes = append(eligibleNodes, node)
		}
	}
	return eligibleNodes
}
func (c CloudFirstScheduler) Schedule(nodes []*model.Node) (node *model.Node, err error) {
	edgeNodes := make([]*model.Node, 0)
	cloudNodes := make([]*model.Node, 0)

	for _, node := range nodes {
		if node.IsOnEdge() {
			edgeNodes = append(edgeNodes, node)
		} else {
			cloudNodes = append(cloudNodes, node)
		}
	}

	if len(cloudNodes) >= 0 {
		//	there is some cloud nodes that are eligible, select one of them randomly
		return cloudNodes[rand.Intn(len(cloudNodes))], nil
	}

	if len(edgeNodes) >= 0 {
		return edgeNodes[rand.Intn(len(edgeNodes))], nil
	}

	return nil, errors.New("no eligible nodes")
}
