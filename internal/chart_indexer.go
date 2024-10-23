package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"

	"github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/repo"
)

// ChartIndexer indexes the charts
type ChartIndexer struct {
	BaseCacheDir string
	BaseURL      string
}

// Reindex reindexes the chart repository
func (ci *ChartIndexer) Reindex(dir string, out string, cache Storage, cloud Storage) error {
	url := ci.BaseURL

	i, err := repo.IndexDirectory(dir, url)
	if err != nil {
		return err
	}
	mergeTo := out
	// if index.yaml is missing then create an empty one to merge into
	var i2 *repo.IndexFile
	if _, err := os.Stat(mergeTo); os.IsNotExist(err) {
		i2 = repo.NewIndexFile()
	} else {
		i2, err = repo.LoadIndexFile(out)
		if err != nil {
			return fmt.Errorf("Merge failed: %s", err)
		}
	}
	i.Merge(i2)
	i.SortEntries()

	data, err := yaml.Marshal(i)
	if err != nil {
		return fmt.Errorf("failed to marshal helm index: %w", err)
	}

	relativePath, err := filepath.Rel(ci.BaseCacheDir, out)
	if err != nil {
		return fmt.Errorf("failed to calculate relative path for %s relative to %s: %w", out, ci.BaseCacheDir, err)
	}

	logrus.Debugf("writing updated chart index at %s", relativePath)

	err = cache.WriteFile(relativePath, io.NopCloser(bytes.NewReader(data)))
	if err != nil {
		return fmt.Errorf("failed to write helm index to cache: %w", err)
	}

	if cloud != nil {
		err = cloud.WriteFile(relativePath, io.NopCloser(bytes.NewReader(data)))
		if err != nil {
			return fmt.Errorf("failed to write helm index to cloud: %w", err)
		}
	}
	return nil
}
