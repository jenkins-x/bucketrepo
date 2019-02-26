package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

// InitLogger initializes the logger
func InitLogger(level string) error {
	logLevel, err := log.ParseLevel(level)
	if err != nil {
		return fmt.Errorf("invalid log level: %s", level)
	}
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(logLevel)
	return nil
}
