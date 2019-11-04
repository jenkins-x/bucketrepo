package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/helm/pkg/repo"
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
		return errors.Wrap(err, "failed to marshal helm index")
	}

	relativePath, err := filepath.Rel(ci.BaseCacheDir, out)
	if err != nil {
		return errors.Wrapf(err, "failed to calculate relative path for %s relative to %s", out, ci.BaseCacheDir)
	}

	logrus.Debugf("writing updated chart index at %s", relativePath)

	err = cache.WriteFile(relativePath, ioutil.NopCloser(bytes.NewReader(data)))
	if err != nil {
		return errors.Wrap(err, "failed to write helm index to cache")
	}

	if cloud != nil {
		err = cloud.WriteFile(relativePath, ioutil.NopCloser(bytes.NewReader(data)))
		if err != nil {
			return errors.Wrap(err, "failed to write helm index to cloud")
		}
	}
	return nil
}
