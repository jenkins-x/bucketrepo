package main

import (
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/julienschmidt/httprouter"
)

type FileController struct {
	Storage    Storage
	Repository Repository
}

func NewFileController(storage Storage, repository Repository) *FileController {
	return &FileController{
		Storage:    storage,
		Repository: repository,
	}
}

func (ctrl *FileController) GetFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	filepath := ps.ByName("filepath")
	log.Debugf("GetFile, filepath: %s", filepath)
	file, err := ctrl.Storage.ReadFile(filepath)
	if err != nil {
		w.WriteHeader(404)
		msg := fmt.Sprintf("Error when reading the file: %s", err)
		log.Error(msg)
		fmt.Fprint(w, msg)
		return
	}
	defer file.Close()

	_, err = io.Copy(w, file)
	if err != nil {
		w.WriteHeader(404)
		msg := fmt.Sprintf("Error when serving the file: %s", err)
		log.Error(msg)
		fmt.Fprintf(w, msg)
		return
	}
}

func (ctrl *FileController) PutFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	filepath := ps.ByName("filepath")
	log.Debugf("PutFile, filepath: %s\n", filepath)
	err := ctrl.Storage.WriteFile(filepath, r.Body)
	if err != nil {
		msg := fmt.Sprintf("Error when saving the file: %s", err)
		w.WriteHeader(500)
		log.Error(msg)
		fmt.Fprint(w, msg)
		return
	}
	w.WriteHeader(200)
}

func (ctrl *FileController) DownloadFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	filepath := ps.ByName("filepath")
	log.Debugf("DownloadFile, filepath: %s\n", filepath)
	err := ctrl.getLocalFile(w, r, ps)
	if err == nil {
		return
	}

	// The file is not available locally. It should be downloaded from the remote repository.
	file, err := ctrl.Repository.DownloadFile(filepath)
	if err != nil {
		w.WriteHeader(404)
		msg := fmt.Sprintf("Error when downloading the file: %s", err)
		log.Error(msg)
		fmt.Fprint(w, msg)
		return
	}
	defer file.Close()
	err = ctrl.Storage.WriteFile(filepath, file)
	if err != nil {
		w.WriteHeader(404)
		msg := fmt.Sprintf("Error when saving the file: %s", err)
		log.Error(msg)
		fmt.Fprint(w, msg)
		return
	}

	err = ctrl.getLocalFile(w, r, ps)
	if err != nil {
		w.WriteHeader(404)
		msg := fmt.Sprintf("Error when reading the local file: %s", err)
		log.Error(msg)
		fmt.Fprint(w, msg)
	}
}

func (ctrl *FileController) getLocalFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	filepath := ps.ByName("filepath")
	log.Debugf("GetLocalFile, filepath: %s\n", filepath)
	file, err := ctrl.Storage.ReadFile(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(w, file)
	if err != nil {
		w.WriteHeader(404)
		msg := fmt.Sprintf("Error when serving the file: %s", err)
		log.Error(msg)
		fmt.Fprintf(w, msg)
		return nil
	}
	return nil
}
