package random

import (
	"errors"
	"math/rand"

	"github.com/noisyboy-9/sencillo/internal/model"
)

type RandScheduler struct{}

func (r RandScheduler) Run(pod model.Pod, nodes []model.Node) (node model.Node, err error) {
	eligibleNodes := r.Filter(pod, nodes)
	if len(eligibleNodes) == 0 {
		return model.Node{}, errors.New("no eligible nodes found")
	}
	return r.Schedule(eligibleNodes)
}

func (r RandScheduler) Filter(pod model.Pod, nodes []model.Node) []model.Node {
	eligibleNodes := make([]model.Node, 0)
	for _, node := range nodes {
		if node.HasEnoughResourcesForPod(pod) {
			eligibleNodes = append(eligibleNodes, node)
		}
	}

	return eligibleNodes
}

func (r RandScheduler) Schedule(eligibleNodes []model.Node) (node model.Node, err error) {
	return eligibleNodes[rand.Intn(len(eligibleNodes))], nil
}
