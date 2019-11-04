package main

import (
	"fmt"
	"os"

	stackdriver "github.com/TV4/logrus-stackdriver-formatter"
	log "github.com/sirupsen/logrus"
)

// InitLogger initializes the logger
func InitLogger(level string) error {
	logLevel, err := log.ParseLevel(level)
	if err != nil {
		return fmt.Errorf("invalid log level: %s", level)
	}
	log.SetFormatter(stackdriver.NewFormatter(
		stackdriver.WithService("bucketrepo"),
		stackdriver.WithVersion("1.0"),
	))
	log.SetOutput(os.Stdout)
	log.SetLevel(logLevel)
	return nil
}
