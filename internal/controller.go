package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
)

const (
	// ChartPackageFileExtension is the file extension used for chart packages
	ChartPackageFileExtension = "tgz"

	// ChartFolder the folder relative to the `index.yaml` file where we store charts for a repository
	ChartFolder = "files"

	// DefaultWritePermissions default permissions for new folders
	DefaultWritePermissions = 0760
)

// FileController controller which handles the artifacts files serving and updating
type FileController struct {
	config           Config
	cache            Storage
	cloudStorage     Storage
	repositories     []Repository
	chartsPath       string
	chartsDir        string
	chartIndexer     *ChartIndexer
	operationChannel chan string
}

var (
	// defaultChartIndex an empty chart index
	defaultChartIndex = `apiVersion: v1
generated: "2019-11-01T17:04:16Z"
entries:`
)

// NewFileController creates a new file controller
func NewFileController(cache Storage, storage Storage, repositories []Repository, config Config) (*FileController, error) {
	chartsPath := config.HTTP.ChartPath
	ctrl := &FileController{
		config:       config,
		cache:        cache,
		cloudStorage: storage,
		repositories: repositories,
		chartsPath:   chartsPath,
	}
	if chartsPath != "" {
		ctrl.chartsDir = filepath.Join(config.Cache.BaseDir, chartsPath)
		ctrl.operationChannel = make(chan string)
		err := os.MkdirAll(ctrl.chartsDir, DefaultWritePermissions)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create charts dir %s", ctrl.chartsDir)
		}

		ctrl.chartIndexer = &ChartIndexer{BaseCacheDir: config.Cache.BaseDir}

		go ctrl.backgroundOperations()

		log.Debugf("now triggering a background reindex")
		ctrl.operationChannel <- "reindex"
	}
	return ctrl, nil
}

// GetFile handlers which returns an artifacts file either from the local file cache, cloud storage or
// central repository
func (ctrl *FileController) GetFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	filename := ps.ByName("filepath")
	log.Debugf("GetFile, filename: %s", filename)
	err := ctrl.readFileFromCache(w, filename, r.Method != "HEAD")
	if err == nil {
		return
	}

	err = ctrl.updateCache(filename)
	if err != nil {
		w.WriteHeader(404)
		msg := fmt.Sprintf("Error when downloading the file: %s", err)
		log.Error(msg)
		fmt.Fprint(w, msg)
		return
	}

	err = ctrl.readFileFromCache(w, filename, r.Method != "HEAD")
	if err != nil {
		w.WriteHeader(404)
		msg := fmt.Sprintf("Error when serving the file from cache: %s", err)
		log.Error(msg)
		fmt.Fprint(w, msg)
		return
	}
}

// PutFile handler which stores an artifact file either into a local file cache or cloud storage
func (ctrl *FileController) PutFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	filename := ps.ByName("filepath")
	log.Debugf("PutFile, filename: %s\n", filename)
	err := ctrl.writeFileToCache(filename, r.Body)
	if err != nil {
		msg := fmt.Sprintf("Error when saving the file into cache: %s", err)
		w.WriteHeader(500)
		log.Error(msg)
		fmt.Fprint(w, msg)
		return
	}

	err = ctrl.updateCloudStorage(filename)
	if err != nil {
		msg := fmt.Sprintf("Error when saving the file into cloud storage: %s", err)
		w.WriteHeader(500)
		log.Error(msg)
		fmt.Fprint(w, msg)
		return
	}

	w.WriteHeader(200)
}

// PostChart handler which stores an artifact file either into a local file cache or cloud storage
func (ctrl *FileController) PostChart(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	repo := ps.ByName("repo")
	log.Debugf("PostChart, repo: %s\n", repo)

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := fmt.Sprintf("Failed to load payload: %s", err)
		w.WriteHeader(500)
		log.Error(msg)
		fmt.Fprint(w, msg)
		return
	}

	filename, err := ctrl.chartPackageFilenameFromContent(content)
	if err != nil {
		msg := fmt.Sprintf("Failed to extract chart from payload: %s", err)
		w.WriteHeader(500)
		log.Error(msg)
		fmt.Fprint(w, msg)
		return
	}

	if path.Base(filename) != filename {
		// Name wants to break out of current directory
		msg := fmt.Sprintf("%s is improperly formattedd ", filename)
		w.WriteHeader(400)
		log.Error(msg)
		fmt.Fprint(w, msg)
		return
	}

	folder := ctrl.chartsPath
	if repo != "" {
		folder = filepath.Join(folder, repo)
	}
	filename = filepath.Join(folder, ChartFolder, filename)
	log.Debugf("PostChart, filename: %s\n", filename)

	err = ctrl.writeFileToCache(filename, ioutil.NopCloser(bytes.NewReader(content)))
	if err != nil {
		msg := fmt.Sprintf("Error when saving the file into cache: %s", err)
		w.WriteHeader(500)
		log.Error(msg)
		fmt.Fprint(w, msg)
		return
	}

	err = ctrl.updateCloudStorage(filename)
	if err != nil {
		msg := fmt.Sprintf("Error when saving the file into cloud storage: %s", err)
		w.WriteHeader(500)
		log.Error(msg)
		fmt.Fprint(w, msg)
	}

	// trigger a reindex
	ctrl.operationChannel <- "reindex"

	w.WriteHeader(200)
}

func (ctrl *FileController) updateCache(filepath string) error {
	file, err := ctrl.downloadFile(filepath)
	updateCloud := false
	if err != nil {
		file, err = ctrl.readFileFromCloudStorage(filepath)
		if err != nil {
			return fmt.Errorf("reading file from cloud storage: %v", err)
		}
		updateCloud = true
	}
	defer file.Close()
	err = ctrl.writeFileToCache(filepath, file)
	if err != nil {
		return fmt.Errorf("writing file to cache: %v", err)
	}

	if updateCloud && ctrl.cloudStorage != nil {
		file, err := ctrl.cache.ReadFile(filepath)
		if err != nil {
			return fmt.Errorf("reading file from cache: %v", err)
		}
		defer file.Close()
		err = ctrl.writeFileToCloudStorage(filepath, file)
		if err != nil {
			log.Warnf("Error when storing the file into cloud storage: %s", err)
			return nil
		}
	}
	return nil
}

func (ctrl *FileController) updateCloudStorage(filepath string) error {
	file, err := ctrl.cache.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("reading file from cache: %v", err)
	}
	defer file.Close()

	return ctrl.writeFileToCloudStorage(filepath, file)
}

func (ctrl *FileController) readFileFromCache(w io.Writer, filepath string, writeBody bool) error {
	log.Debugf("Read file from cache: %s\n", filepath)
	file, err := ctrl.cache.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("reading file from cache: %v", err)
	}
	defer file.Close()
	if !writeBody {
		return nil
	}
	_, err = io.Copy(w, file)
	if err != nil {
		return fmt.Errorf("copying file from cache: %v", err)
	}
	return nil
}

func (ctrl *FileController) writeFileToCache(filepath string, file io.ReadCloser) error {
	log.Debugf("Write file to cache: %s\n", filepath)
	return ctrl.cache.WriteFile(filepath, file)
}

func (ctrl *FileController) readFileFromCloudStorage(filepath string) (io.ReadCloser, error) {
	if ctrl.cloudStorage != nil {
		log.Debugf("Read file from cloud storage: %s", filepath)
		return ctrl.cloudStorage.ReadFile(filepath)
	}
	return nil, errors.New("no cloud storage available")
}

func (ctrl *FileController) writeFileToCloudStorage(filepath string, file io.ReadCloser) error {
	if ctrl.cloudStorage != nil {
		log.Debugf("Write file to cloud storage: %s", filepath)
		return ctrl.cloudStorage.WriteFile(filepath, file)
	}
	return nil
}

func (ctrl *FileController) downloadFile(filepath string) (io.ReadCloser, error) {
	log.Debugf("Read file from repository: %s", filepath)
	for _, r := range ctrl.repositories {
		log.Debugf("Trying to download from repository: %s", r.BaseURL())
		b, err := r.DownloadFile(filepath)
		if err == nil {
			return b, nil
		}
	}

	return nil, fmt.Errorf("unable to download %s from any configured repository", filepath)
}

func (ctrl *FileController) chartReindex() error {
	dir := ctrl.chartsDir
	filename := filepath.Join(dir, "index.yaml")

	err := ctrl.updateCache(filename)
	if err != nil {
		if ctrl.cloudStorage == nil {
			log.Debugf("no cloud storage so failed to update cache for %s: %s", filename, err.Error())
		} else {
			log.Errorf("failed to update cache for %s: %s", filename, err.Error())
		}
	}

	err = ctrl.chartIndexer.Reindex(dir, filename, ctrl.cache, ctrl.cloudStorage)
	if err != nil {
		log.Errorf("failed to reindex the charts: %s", err.Error())
		return err
	}
	return nil
}

// ChartPackageFilenameFromContent returns a chart filename from binary content
func (ctrl *FileController) chartPackageFilenameFromContent(content []byte) (string, error) {
	chart, err := chartFromContent(content)
	if err != nil {
		return "", err
	}
	meta := chart.Metadata
	filename := fmt.Sprintf("%s-%s.%s", meta.Name, meta.Version, ChartPackageFileExtension)
	return filename, nil
}

func (ctrl *FileController) backgroundOperations() {
	for {
		<-ctrl.operationChannel

		log.Debugf("reindexing charts")

		err := ctrl.chartReindex()
		if err != nil {
			log.Errorf("failed to reindex charts: %s", err.Error())
		}
	}
}

func chartFromContent(content []byte) (*chart.Chart, error) {
	chart, err := loader.LoadArchive(bytes.NewBuffer(content))
	return chart, err
}
