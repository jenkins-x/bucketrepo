package main

import (
	"io"
)

// Storage interfaces for artifacts storage
type Storage interface {
	// ReadFile reads a file from the storage
	ReadFile(path string) (io.ReadCloser, error)
	// WriteFile wrietes a file into the storage
	WriteFile(path string, file io.ReadCloser) error
}

// NewStorage creates a new storage
func NewStorage(config StorageConfig) Storage {
	if config.Enabled {
		return NewCloudStorage(config)
	}
	return nil
}
