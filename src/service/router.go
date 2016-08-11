package main

import (
	"github.com/gorilla/mux"
	"net/http"
)


func NewPreRouter(preRoutingRule func(request *http.Request) bool) *PreRouter{

	router := mux.NewRouter();
	subRouter := router.PathPrefix("/catalog/").Subrouter().StrictSlash(true)

	preRouter := &PreRouter{router: subRouter, evaluationRule: preRoutingRule}
	return preRouter
}

type PreRouter struct {
	router         *mux.Router
	evaluationRule func(request *http.Request) bool
}

//This method is called by http.ListenAndServe when a request comes in
func (preRouter *PreRouter) ServeHTTP(resp http.ResponseWriter, request *http.Request){

	//The purpose of PreRouter is to be able to make any processing before requests are evaluated
	if !preRouter.evaluationRule(request){
		ForbiddenHandler(resp, request)
		return
	}

	//Here begins the processing of the different requests
	preRouter.router.ServeHTTP(resp, request)
}

func (preRouter *PreRouter) AddService(service Service) *PreRouter{

	service.ConfigureRouter(preRouter.GetRouter())

	return preRouter
}

func (preRouter PreRouter) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route{
	return preRouter.router.HandleFunc(path, f)
}

func (preRouter PreRouter) GetRouter() *mux.Router{
	return preRouter.router
}

func ForbiddenHandler(w http.ResponseWriter, r *http.Request) { http.Error(w, "503 Forbidden", http.StatusForbidden) }