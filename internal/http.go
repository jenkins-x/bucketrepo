package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

func BasicAuth(h httprouter.Handle, config HttpConfig) httprouter.Handle {
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

func InitHttp(config HttpConfig, controller *FileController) {
	router := httprouter.New()
	router.GET("/*filepath", BasicAuth(controller.GetFile, config))
	router.PUT("/*filepath", BasicAuth(controller.PutFile, config))

	log.Infof("Start http server on %s", config.Address)
	if config.HTTPS {
		log.Fatal(http.ListenAndServeTLS(config.Address, config.Certificate, config.Key, router))
		return
	}
	log.Fatal(http.ListenAndServe(config.Address, router))
}
