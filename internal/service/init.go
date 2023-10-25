package service

import "context"

func Init() {
	InitHttpServer()
}

func Terminate(cancelCtx context.Context) {
	TerminateHttpServer(cancelCtx)
}
