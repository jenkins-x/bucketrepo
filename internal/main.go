package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type flags struct {
	logLevel   string
	configPath string
}

func main() {
	flags := flags{}
	flag.StringVar(&flags.logLevel, "log-level", "warning", "Defines the logs level (debug, info, warning, error, fatal, panic)")
	flag.StringVar(&flags.configPath, "config-path", ".", "Define the absolute path to the config.yaml file")
	flag.Parse()

	err := InitLogger(flags.logLevel)
	if err != nil {
		fmt.Printf("Invalid log-level option: %s", err)
		os.Exit(2)
	}

	config := NewConfig(flags.configPath)

	storage := NewStorage(config.Storage)
	cache := NewFileSystemStorage(config.Cache)
	repositories := make([]Repository, len(config.Repositories))
	for i, r := range config.Repositories {
		repositories[i] = NewRepository(r)
	}

	controller, err := NewFileController(cache, storage, repositories, config)
	if err != nil {
		logrus.Fatalf("failed to initialise controller: %s", err.Error())
	}

	if config.Cache.CacheTime != 0 {
		ticker := time.NewTicker(config.Cache.CleanInterval)
		go func(ticker *time.Ticker) {
			cache.RemoveUnusedArtifacts(controller)
			for range ticker.C {
				cache.RemoveUnusedArtifacts(controller)
			}
		}(ticker)
	}

	logrus.Infof("serving http")
	InitHTTP(config.HTTP, controller)
}
