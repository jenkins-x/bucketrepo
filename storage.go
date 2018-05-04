package main

import (
	"io"
)

type Storage interface {
	ReadFile(path string) (io.ReadCloser, error)
	WriteFile(path string, file io.ReadCloser) error
}

func NewStorage(config StorageConfig) Storage {
	switch config.Type {
	case "s3":
		return NewS3Storage(config)
	case "fs":
		return NewFileSystemStorage(config)
	default:
		panic("Unknown storage type")
	}
}
