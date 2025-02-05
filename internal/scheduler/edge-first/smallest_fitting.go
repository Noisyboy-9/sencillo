package edge_first

import (
	"errors"
	"math/rand"

	"github.com/noisyboy-9/sencillo/internal/log"
	"github.com/noisyboy-9/sencillo/internal/util"

	"github.com/noisyboy-9/sencillo/internal/model"
)

type SmallestFittingEdgeNodeScheduler struct{}

func (s SmallestFittingEdgeNodeScheduler) Run(pod model.Pod, nodes []model.Node) (node model.Node, err error) {
	edgeNodes, cloudNodes := s.Filter(pod, nodes)
	if len(edgeNodes) == 0 && len(cloudNodes) == 0 {
		return model.Node{}, errors.New("no eligible nodes found")
	}
	return s.Schedule(edgeNodes, cloudNodes), nil
}

func (s SmallestFittingEdgeNodeScheduler) Filter(pod model.Pod, nodes []model.Node) (eligibleEdgeNodes []model.Node, eligibleCloudNodes []model.Node) {
	eligibleEdgeNodes = make([]model.Node, 0)
	eligibleCloudNodes = make([]model.Node, 0)
	for _, node := range nodes {
		if node.IsMaster {
			continue
		}

		if node.HasEnoughResourcesForPod(pod) && node.IsOnEdge {
			eligibleEdgeNodes = append(eligibleEdgeNodes, node)
		}

		if node.HasEnoughResourcesForPod(pod) && !node.IsOnEdge {
			eligibleCloudNodes = append(eligibleCloudNodes, node)
		}
	}

	return eligibleEdgeNodes, eligibleCloudNodes
}

func (s SmallestFittingEdgeNodeScheduler) Schedule(edgeNodes []model.Node, cloudNodes []model.Node) (node model.Node) {
	if len(edgeNodes) != 0 {
		return util.FindSmallestNode(edgeNodes)
	}

	log.App.Info("no edge nodes were eligible, scheduling on random cloud node ...")
	return cloudNodes[rand.Intn(len(cloudNodes))]
}
