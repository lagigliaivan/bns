package services

import (
"github.com/gorilla/mux"
"net/http"
)

type Service interface {
	ConfigureRouter(router Router)
}

func getPathParams(r *http.Request) map[string]string {
	return mux.Vars(r)
}