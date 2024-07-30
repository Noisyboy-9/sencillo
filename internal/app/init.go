package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/noisyboy-9/sencillo/internal/config"
	"github.com/noisyboy-9/sencillo/internal/connector"
	"github.com/noisyboy-9/sencillo/internal/consumer"
	"github.com/noisyboy-9/sencillo/internal/log"
	"github.com/noisyboy-9/sencillo/internal/service"
)

var Done = make(chan os.Signal, 1)

func InitApp() {
	config.LoadViper()
	log.Init()
	config.Init()
	service.Init()
	connector.Connect()
	consumer.Start()
}

func SetupGracefulShutdown() {
	signal.Notify(Done, syscall.SIGINT, syscall.SIGTERM)
	<-Done
	log.App.Info("Exit signal has been received. Shutting down . . .")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	terminateApp(ctx)
}

func terminateApp(_ context.Context) {
	service.Terminate()
}
