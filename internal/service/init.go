package service

import "context"

func Init() {
	NewScheduler()
}

func Terminate(cancelCtx context.Context) {
}
