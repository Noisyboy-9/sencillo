package scheduler

import (
	"math/rand"

	"github.com/noisyboy-9/random-k8s-scheduler/internal/config"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/model"
)

type randomScheduler struct {
	Name      string
	Namespace string
}

func newRandomScheduler() Scheduler {
	rs := &randomScheduler{
		Name:      config.Scheduler.Name,
		Namespace: config.Scheduler.Namespace,
	}

	return rs
}

func (r randomScheduler) Run(pod *model.Pod, nodes []*model.Node) (node *model.Node, err error) {
	eligibleNodes := r.Filter(pod, nodes)
	return r.Schedule(eligibleNodes)
}

func (r randomScheduler) Filter(pod *model.Pod, nodes []*model.Node) []*model.Node {
	eligibleNodes := make([]*model.Node, 0)
	for _, node := range nodes {
		if node.HasEnoughResourcesForPod(pod) {
			eligibleNodes = append(eligibleNodes, node)
		}
	}
	return eligibleNodes
}

func (r randomScheduler) Schedule(eligibleNodes []*model.Node) (node *model.Node, err error) {
	return eligibleNodes[rand.Intn(len(eligibleNodes))], nil
}
