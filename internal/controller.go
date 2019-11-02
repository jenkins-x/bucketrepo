package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

// FileController controller which handles the artifacts files serving and updating
type FileController struct {
	cache        Storage
	cloudStorage Storage
	repositories []Repository
}

var (
	// defaultChartIndex an empty chart index
	defaultChartIndex = `apiVersion: v1
generated: "2019-11-01T17:04:16Z"
entries:`
)

// NewFileController creates a new file controller
func NewFileController(cache Storage, storage Storage, repositories []Repository, chartsPath string) *FileController {
	ctrl := &FileController{
		cache:        cache,
		cloudStorage: storage,
		repositories: repositories,
	}
	if chartsPath != "" {
		ctrl.ensureChartIndex(filepath.Join(chartsPath, "index.yaml"))
	}
	return ctrl
}

// GetFile handlers which returns an artifacts file either from the local file cache, cloud storage or
// central repository
func (ctrl *FileController) GetFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	filepath := ps.ByName("filepath")
	log.Debugf("GetFile, filepath: %s", filepath)
	err := ctrl.readFileFromCache(w, filepath, r.Method != "HEAD")
	if err == nil {
		return
	}

	err = ctrl.updateCache(filepath)
	if err != nil {
		w.WriteHeader(404)
		msg := fmt.Sprintf("Error when downlaoding the file: %s", err)
		log.Error(msg)
		fmt.Fprint(w, msg)
		return
	}

	err = ctrl.readFileFromCache(w, filepath, r.Method != "HEAD")
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
	filepath := ps.ByName("filepath")
	log.Debugf("PutFile, filepath: %s\n", filepath)
	err := ctrl.writeFileToCache(filepath, r.Body)
	if err != nil {
		msg := fmt.Sprintf("Error when saving the file into cache: %s", err)
		w.WriteHeader(500)
		log.Error(msg)
		fmt.Fprint(w, msg)
		return
	}

	err = ctrl.updateCloudStorage(filepath)
	if err != nil {
		msg := fmt.Sprintf("Error when saving the file into cloud storage: %s", err)
		w.WriteHeader(500)
		log.Error(msg)
		fmt.Fprint(w, msg)
		return
	}

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

func (ctrl *FileController) readFileFromCache(w http.ResponseWriter, filepath string, writeBody bool) error {
	log.Debugf("Read file form cache: %s\n", filepath)
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
	return nil, errors.New("no cloud storage availalbe")
}

func (ctrl *FileController) writeFileToCloudStorage(filepath string, file io.ReadCloser) error {
	if ctrl.cloudStorage != nil {
		log.Debugf("Write file to cloud storage: %s", filepath)
		return ctrl.cloudStorage.WriteFile(filepath, file)
	}
	return nil
}

func (ctrl *FileController) downloadFile(filepath string) (io.ReadCloser, error) {
	log.Debugf("Read file form repository: %s", filepath)
	for _, r := range ctrl.repositories {
		log.Debugf("Trying to download from repository: %s", r.BaseURL())
		b, err := r.DownloadFile(filepath)
		if err == nil {
			return b, nil
		}
	}

	return nil, fmt.Errorf("unable to download %s from any configured repository", filepath)
}

func (ctrl *FileController) ensureChartIndex(filepath string) error {
	// ignore errors if no cloud storage
	ctrl.updateCache(filepath)

	file, err := ctrl.cache.ReadFile(filepath)
	if err != nil {
		log.Infof("no charts file %s so generating it", filepath)
		return ctrl.cache.WriteFile(filepath, ioutil.NopCloser(strings.NewReader(defaultChartIndex)))
	}
	defer file.Close()
	return nil
}
