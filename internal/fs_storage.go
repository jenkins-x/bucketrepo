package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/djherbis/atime"
	"github.com/sirupsen/logrus"
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
// TODO: If writing fails due to out of disk files with oldest access timestamp should be removed
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

// RemoveUnusedArtifacts cleans away artifacts that have not been used for configurable amount of time
func (fs *FileSystemStorage) RemoveUnusedArtifacts(ctrl *FileController) {
	logrus.Info("starting removal of unused artifacts")
	cacheTime := fs.config.CacheTime
	maxAccessTime := time.Now().Add(-cacheTime)
	var err error
	// Don't remove charts directory or index
	var keepRegExp *regexp.Regexp
	if ctrl.chartsDir != "" {
		keepString := fmt.Sprintf("%s(?:%cindex.yaml)?", regexp.QuoteMeta(ctrl.chartsDir), os.PathSeparator)
		keepRegExp, err = regexp.Compile(keepString)
		if err != nil {
			logrus.WithError(err).Errorf("can't compile string %s to regexp", keepString)
		}
	}
	err = filepath.Walk(fs.config.BaseDir, func(path string, info os.FileInfo, err error) error {
		if keepRegExp != nil && keepRegExp.MatchString(path) {
			return nil
		}
		aTime := atime.Get(info)
		if aTime.Before(maxAccessTime) {
			logrus.Debugf("removing %s that has not been accessed since %s", path, aTime)
			err := os.RemoveAll(path)
			if err != nil {
				return err
			}
			if info.IsDir() {
				return filepath.SkipDir
			}
		}
		return nil
	})
	if err != nil {
		logrus.WithError(err).Errorf("failed to remove unused artifacts")
	}
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
