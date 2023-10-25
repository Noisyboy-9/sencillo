package config

type httpServer struct {
	Host string `default:"localhost"`
	Port int    `default:"8080"`
}

var HttpServer *httpServer
