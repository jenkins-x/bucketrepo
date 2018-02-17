package main

import (
	"io"

	"github.com/spf13/viper"
)

type Storage interface {
	ReadFile(path string) (io.ReadCloser, error)
	WriteFile(path string, file io.ReadCloser) error
}

func NewStorage(storageType string) Storage {
	switch storageType {
	case "s3":
		bucket := viper.GetString("storage.bucket")
		accessKey := viper.GetString("storage.access_key")
		secretKey := viper.GetString("storage.secret_key")
		return NewS3Storage(bucket, accessKey, secretKey)
	case "fs":
		baseDir := viper.GetString("storage.base_dir")
		return NewFileSystemStorage(baseDir)
	default:
		panic("Unknown storage type")
	}
}
