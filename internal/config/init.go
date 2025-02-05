package config

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/noisyboy-9/sencillo/internal/log"
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
		return
	}

	Connector = new(connector)
	err = viper.UnmarshalKey("connector", Connector)
	if err != nil {
		log.App.WithError(err).Panic("can't read connector configs")
		return
	}

	Cluster = new(cluster)
	err = viper.UnmarshalKey("cluster", Cluster)
	if err != nil {
		log.App.WithError(err).Panic("can't read cluster configs")
		return
	}
}
