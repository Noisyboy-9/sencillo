package cloud_first

import (
	"errors"
	"math/rand"

	"github.com/noisyboy-9/sencillo/internal/log"
	"github.com/noisyboy-9/sencillo/internal/model"
)

type CloudFirstScheduler struct{}

func (c CloudFirstScheduler) Run(pod *model.Pod, nodes []*model.Node) (node *model.Node, err error) {
	edgeNodes, cloudNodes := c.Filter(pod, nodes)
	if len(edgeNodes) == 0 && len(cloudNodes) == 0 {
		return nil, errors.New("no eligible nodes found")

	}
	return c.Schedule(edgeNodes, cloudNodes), nil
}
func (c CloudFirstScheduler) Filter(pod *model.Pod, nodes []*model.Node) (eligibleEdgeNodes []*model.Node, eligibleCloudNodes []*model.Node) {
	eligibleEdgeNodes = make([]*model.Node, 0)
	eligibleCloudNodes = make([]*model.Node, 0)
	for _, node := range nodes {
		if node.HasEnoughResourcesForPod(pod) && node.IsOnEdge {
			eligibleEdgeNodes = append(eligibleEdgeNodes, node)
		}

		if node.HasEnoughResourcesForPod(pod) && !node.IsOnEdge {
			eligibleCloudNodes = append(eligibleCloudNodes, node)
		}
	}

	return eligibleEdgeNodes, eligibleCloudNodes
}

func (c CloudFirstScheduler) Schedule(edgeNodes []*model.Node, cloudNodes []*model.Node) (node *model.Node) {
	if len(cloudNodes) != 0 {
		return cloudNodes[rand.Intn(len(cloudNodes))]
	}

	log.App.Info("no cloud nodes were eligible, scheduling on random edge node ...")
	return edgeNodes[rand.Intn(len(edgeNodes))]
}
