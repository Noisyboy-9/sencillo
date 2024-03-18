package scheduler

import "github.com/noisyboy-9/random-k8s-scheduler/internal/model"

type Scheduler interface {
	Run(pod *model.Pod, nodes []*model.Node) (node *model.Node, err error)
}

var S Scheduler

//func NewScheduler() {
//
//}
