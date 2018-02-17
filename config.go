package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/nexus-minimal/")
	viper.AddConfigPath("$HOME/.nexus-minimal")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Warn("error during reading config file")
	}
}
