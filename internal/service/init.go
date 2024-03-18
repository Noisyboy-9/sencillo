package service

import (
	"github.com/noisyboy-9/random-k8s-scheduler/internal/scheduler"
)

func Init() {
	scheduler.NewScheduler()
}

func Terminate() {
}
