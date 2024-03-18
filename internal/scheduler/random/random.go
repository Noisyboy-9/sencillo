package random

import (
	"math/rand"

	"github.com/noisyboy-9/random-k8s-scheduler/internal/model"
)

type RandScheduler struct {
	Name      string
	Namespace string
}

func (r RandScheduler) Run(pod *model.Pod, nodes []*model.Node) (node *model.Node, err error) {
	eligibleNodes := r.Filter(pod, nodes)
	return r.Schedule(eligibleNodes)
}

func (r RandScheduler) Filter(pod *model.Pod, nodes []*model.Node) []*model.Node {
	eligibleNodes := make([]*model.Node, 0)
	for _, node := range nodes {
		if node.HasEnoughResourcesForPod(pod) {
			eligibleNodes = append(eligibleNodes, node)
		}
	}
	return eligibleNodes
}

func (r RandScheduler) Schedule(eligibleNodes []*model.Node) (node *model.Node, err error) {
	return eligibleNodes[rand.Intn(len(eligibleNodes))], nil
}
