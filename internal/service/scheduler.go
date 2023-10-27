package service

import (
	"github.com/noisyboy-9/random-k8s-scheduler/internal/config"
	"github.com/noisyboy-9/random-k8s-scheduler/internal/model"
)

type scheduler struct {
	Name string
}

var Scheduler *scheduler

func NewScheduler() {
	Scheduler = &scheduler{
		Name: config.Scheduler.Name,
	}
}

func (scheduler *scheduler) FindNodeForBinding(pod *model.Pod, nodes []*model.Node) (node *model.Node, err error) {
	return nil, nil
}
