package edge_first

import (
	"errors"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/log"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/model"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/util"
	"math/rand"
)

type BiggestFittingEdgeNodeScheduler struct{}

func (b BiggestFittingEdgeNodeScheduler) Run(pod *model.Pod, nodes []*model.Node) (node *model.Node, err error) {
	edgeNodes, cloudNodes := b.Filter(pod, nodes)
	if len(edgeNodes) == 0 && len(cloudNodes) == 0 {
		return nil, errors.New("no eligible nodes found")
	}
	return b.Schedule(edgeNodes, cloudNodes), nil
}

func (b BiggestFittingEdgeNodeScheduler) Filter(pod *model.Pod, nodes []*model.Node) (eligibleEdgeNodes []*model.Node, eligibleCloudNodes []*model.Node) {
	eligibleEdgeNodes = make([]*model.Node, 0)
	eligibleCloudNodes = make([]*model.Node, 0)
	for _, node := range nodes {
		if node.HasEnoughResourcesForPod(pod) && node.GetIsOnEdge() {
			eligibleEdgeNodes = append(eligibleEdgeNodes, node)
		}

		if node.HasEnoughResourcesForPod(pod) && !node.GetIsOnEdge() {
			eligibleCloudNodes = append(eligibleCloudNodes, node)
		}
	}

	return eligibleEdgeNodes, eligibleCloudNodes
}
func (b BiggestFittingEdgeNodeScheduler) Schedule(edgeNodes []*model.Node, cloudNodes []*model.Node) (node *model.Node) {
	if len(edgeNodes) != 0 {
		return util.FindLargestNode(edgeNodes)
	}

	log.App.Info("no edge nodes were eligible, scheduling on random cloud node ...")
	return cloudNodes[rand.Intn(len(cloudNodes))]
}
