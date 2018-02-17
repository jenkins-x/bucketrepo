package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func BasicAuth(h httprouter.Handle, requiredUser, requiredPassword string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user, password, hasAuth := r.BasicAuth()

		if hasAuth && user == requiredUser && password == requiredPassword {
			h(w, r, ps)
		} else {
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}

func InitHttp(controller *FileController) {
	username := viper.GetString("http.username")
	password := viper.GetString("http.password")

	router := httprouter.New()
	router.GET("/*filepath", BasicAuth(controller.GetFile, username, password))
	router.PUT("/*filepath", BasicAuth(controller.PutFile, username, password))

	httpAddr := viper.GetString("http.addr")
	log.Infof("Start http server on %s", httpAddr)
	log.Fatal(http.ListenAndServe(httpAddr, router))
}
