package scheduler

import (
	"github.com/noisyboy-9/random-k8s-scheduler/internal/config"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/enum"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/log"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/model"
)

type Scheduler interface {
	Run(pod *model.Pod, nodes []*model.Node) (node *model.Node, err error)
}

var S Scheduler

func NewScheduler() {
	schedulerType, err := enum.ParseAlgorithmName(config.Scheduler.Algorithm)
	if err != nil {
		log.App.WithError(err).Panic("Error in creating scheduler")
	}

	switch schedulerType {
	case enum.RandomScheduler:
		S = newRandomScheduler()
		log.App.Info("random scheduler created successfully")
	default:
		log.App.Panic("not known scheduler type")
	}
}
