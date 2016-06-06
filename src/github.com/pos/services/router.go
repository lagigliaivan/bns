package services

import (
"github.com/gorilla/mux"
"net/http"
)

func NewRouter() *mux.Router{

	router := mux.NewRouter();
	router.Methods(http.MethodGet, http.MethodPut, http.MethodPost)
	subRouter := router.PathPrefix("/catalog/").Headers("Authorization", "").Subrouter()

	return subRouter
}