package config

import (
	"fmt"

	"github.com/noisyboy-9/sencillo/internal/log"
	"github.com/spf13/viper"
)

func LoadViper() {
	configPath := "configs/general.yaml"
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("config file in path: %s can't be found", configPath))
	}
}

func Init() {
	var err error

	Scheduler = new(scheduler)
	err = viper.UnmarshalKey("scheduler", Scheduler)
	if err != nil {
		log.App.WithError(err).Panic("can't read scheduler configs")
	}

	Connector = new(connector)
	err = viper.UnmarshalKey("connector", Connector)
	if err != nil {
		log.App.WithError(err).Panic("can't read connector configs")
	}
}
