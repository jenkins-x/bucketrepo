package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

// basicAuth creates a wrapper for basic authentication
func basicAuth(h httprouter.Handle, config HTTPConfig) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user, password, hasAuth := r.BasicAuth()

		if hasAuth && user == config.Username && password == config.Password {
			h(w, r, ps)
		} else {
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}

// noAuth creates a wrapper without any authentication
func noAuth(h httprouter.Handle, config HTTPConfig) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		h(w, r, ps)
	}
}

// health handles the health check requests
func health(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "OK")
}

// InitHTTP initializes the HTTP server routes
func InitHTTP(config HTTPConfig, controller *FileController) {
	router := httprouter.New()
	// Enable basic auth only for https
	auth := noAuth
	if config.HTTPS {
		auth = basicAuth
	}
	router.GET("/healthz", health)
	router.GET("/bucketrepo/*filepath", auth(controller.GetFile, config))
	router.PUT("/bucketrepo/*filepath", auth(controller.PutFile, config))

	log.Infof("Starting http server on %q", config.Address)
	if config.HTTPS {
		log.Fatal(http.ListenAndServeTLS(config.Address, config.Certificate, config.Key, router))
		return
	}
	log.Fatal(http.ListenAndServe(config.Address, router))
}
