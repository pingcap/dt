package controller

import (
	"github.com/gorilla/mux"
	"net/http"

	"testingframe/pkg/util"
)

func runHttpServer(addr string, ctrl *Controller) {
	m := mux.NewRouter()

	m.HandleFunc(util.UrlRegisterAgent, ctrl.apiRegisterAgent).Methods("POST", "PUT")

	http.Handle("/", m)
	http.ListenAndServe(addr, nil)
}

// TODO: implement
func (ctrl *Controller) apiRegisterAgent(w http.ResponseWriter, r *http.Request) {
	panic("apiRegisterAgent hasn't implemented")
}
