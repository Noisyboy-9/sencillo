package service

import (
	"github.com/noisyboy-9/sencillo/internal/scheduler"
)

func Init() {
	scheduler.NewScheduler()
}

func Terminate() {
}
