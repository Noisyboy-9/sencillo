package random

import (
	"math/rand"

	"github.com/noisyboy-9/random-k8s-scheduler/internal/config"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/model"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/scheduler"
)

type randomScheduler struct {
	Name      string
	Namespace string
}

func newRandomScheduler() scheduler.Scheduler {
	rs := &randomScheduler{
		Name:      config.Scheduler.Name,
		Namespace: config.Scheduler.Namespace,
	}

	return rs
}

func (r randomScheduler) Run(pod *model.Pod, nodes []*model.Node) (node *model.Node, err error) {
	//filtering step
	eligibleNodes := make([]*model.Node, 0)
	for _, node := range nodes {
		if node.HasEnoughResourcesForPod(pod) {
			eligibleNodes = append(eligibleNodes, node)
		}
	}

	return r.Schedule(pod, eligibleNodes)
}

func (r randomScheduler) Schedule(pod *model.Pod, eligibleNodes []*model.Node) (node *model.Node, err error) {
	return eligibleNodes[rand.Intn(len(eligibleNodes))], nil
}
