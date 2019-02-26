package main

import (
	"flag"
	"fmt"
	"os"
)

type flags struct {
	LogLevel string
}

func main() {
	flags := flags{}
	flag.StringVar(&flags.LogLevel, "log-level", "warning", "Defines the logs level (debug, info, warning, error, fatal, panic)")
	flag.Parse()

	err := InitLogger(flags.LogLevel)
	if err != nil {
		fmt.Printf("Invalid log-level option: %s", err)
		os.Exit(2)
	}

	config := NewConfig()

	storage := NewStorage(config.Storage)
	cache := NewFileSystemStorage(config.Cache)
	repository := NewRepository(config.Repository)
	controller := NewFileController(cache, storage, repository)

	InitHttp(config.HTTP, controller)
}
