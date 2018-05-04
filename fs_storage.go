package main

import (
	"io"
	"os"
	"path"
	"strings"
)

type FileSystemStorage struct {
	Config StorageConfig
}

func NewFileSystemStorage(config StorageConfig) *FileSystemStorage {
	return &FileSystemStorage{
		Config: config,
	}
}

func (fs *FileSystemStorage) ReadFile(path string) (io.ReadCloser, error) {
	fullPath := resolvePath(fs.Config.BaseDir, path)
	return os.Open(fullPath)
}

func (fs *FileSystemStorage) WriteFile(path string, file io.ReadCloser) error {
	fullPath := resolvePath(fs.Config.BaseDir, path)
	directoryPath, _ := parseFilepath(fullPath)
	os.MkdirAll(directoryPath, os.ModePerm)
	outFile, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer outFile.Close()
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
