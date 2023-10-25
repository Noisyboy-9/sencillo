package service

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/noisyboy-9/golang_api_template/internal/config"
	"github.com/noisyboy-9/golang_api_template/internal/handler"
	"github.com/noisyboy-9/golang_api_template/internal/log"
)

type httpServer struct {
	e echo.Echo
}

var HttpServer *httpServer

func InitHttpServer() {
	HttpServer = new(httpServer)
	HttpServer.e = *echo.New()
	HttpServer.e.HideBanner = true

	HttpServer.setupMiddlewares()
	HttpServer.registerRoutes()

	go func() {
		serverUrl := fmt.Sprintf("%s:%d", config.HttpServer.Host, config.HttpServer.Port)
		if err := HttpServer.e.Start(serverUrl); err != nil {
			log.App.WithField("err", err.Error()).Fatalf("can't start web server")
		}
	}()
}

func (server *httpServer) registerRoutes() {
	HttpServer.e.GET("/", handler.SayHello)
}

func (server *httpServer) setupMiddlewares() {
	HttpServer.e.Use(middleware.Logger())
}

func TerminateHttpServer(ctx context.Context) {
	HttpServer.e.Shutdown(ctx)
}
