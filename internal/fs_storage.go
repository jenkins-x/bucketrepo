package main

import (
	"io"
	"os"
	"path"
	"strings"
)

// FileSystemStorage file storage backend
type FileSystemStorage struct {
	config CacheConfig
}

// NewFileSystemStorage creates a new file system file storage
func NewFileSystemStorage(config CacheConfig) *FileSystemStorage {
	return &FileSystemStorage{
		config: config,
	}
}

// ReadFile reads a file from the local file system
func (fs *FileSystemStorage) ReadFile(path string) (io.ReadCloser, error) {
	fullPath := resolvePath(fs.config.BaseDir, path)
	return os.Open(fullPath) // #nosec
}

// WriteFile writes a file into the local file system
func (fs *FileSystemStorage) WriteFile(path string, file io.ReadCloser) error {
	fullPath := resolvePath(fs.config.BaseDir, path)
	directoryPath, _ := parseFilepath(fullPath)
	err := os.MkdirAll(directoryPath, 0750)
	if err != nil {
		return err
	}
	outFile, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer outFile.Close() // #nosec
	_, err = io.Copy(outFile, file)
	return err
}

func resolvePath(basedir string, filepath string) string {
	return path.Join(basedir, filepath)
}

func parseFilepath(filepath string) (directoryPath, filename string) {
	segments := strings.Split(filepath, "/")
	filename = segments[len(segments)-1]
	directoryPath = strings.Join(segments[:len(segments)-1], "/")
	return
}
