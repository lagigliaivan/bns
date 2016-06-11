package services

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
)

func NewRouter() *PreRouter{

	router := mux.NewRouter();
	router.Methods(http.MethodGet, http.MethodPut, http.MethodPost)
	subRouter := router.PathPrefix("/catalog/").Subrouter()

	return NewPreRouter(subRouter, validateUser)
}

func NewPreRouter(fwRouter *mux.Router, rule func(request *http.Request) bool) *PreRouter{

	preRouter := &PreRouter{router: fwRouter, rule: rule}
	return preRouter
}

type PreRouter struct {
	router *mux.Router
	rule func(request *http.Request) bool
}

func (preRouter *PreRouter) ServeHTTP(resp http.ResponseWriter, request *http.Request){

	if !preRouter.rule(request){
		http.NotFound(resp, request)
		return
	}

	preRouter.router.ServeHTTP(resp, request)

}

func (preRouter PreRouter) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route{
	return preRouter.router.HandleFunc(path, f)
}

type Router interface  {
	HandleFunc(string, func(http.ResponseWriter, *http.Request)) *mux.Route
}

func validateUser (request *http.Request) bool{
	securityHeader := request.Header.Get("Security")

	if len(securityHeader) == 0 {
		log.Printf("Security Header needs to be present")
		return false
	}
	
	log.Printf("Security Header: %s", securityHeader)
	return true
}