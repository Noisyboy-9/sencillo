package log

import (
	"os"

	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	HttpServerLogPrefix = "httpServer"
)

var (
	App *logrus.Logger
)

func Init() {
	App = logrus.New()

	App.SetFormatter(&runtime.Formatter{
		ChildFormatter: &logrus.JSONFormatter{},
		Line:           true,
		File:           true,
		Package:        false,
		BaseNameOnly:   false,
	})

	configPrefix := "logging.app"

	if viper.GetBool(configPrefix + "stdout") {
		App.SetOutput(os.Stdout)
	}

	level, err := logrus.ParseLevel(viper.GetString("logging.level"))
	if err != nil {
		panic(err)
	}
	App.SetLevel(level)
}
