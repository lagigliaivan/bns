package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"strings"
)

const(
	HEADER = "Authorization"
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
		ForbiddenHandler(resp, request)
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

var users [2]string =  [...]string{"d563af2d08b4f672a11b3ed9065b7890a6412cab", "107cbb20a1d1e156beac1a9a7a331b36321300d4"}

func validateUser (request *http.Request) bool{

	securityHeader := request.Header.Get(HEADER)

	if len(securityHeader) == 0 {
		log.Printf("Security Header needs to be present")
		return false
	}else {

		for _, allowedUser := range users {
			if strings.Compare(securityHeader, allowedUser) == 0 {
				log.Printf("user:%s\n", allowedUser)
				return true
			}
		}

		log.Printf("User %s is not allowed", securityHeader)
	}

	return false
}

func ForbiddenHandler(w http.ResponseWriter, r *http.Request) { http.Error(w, "503 Forbidden", http.StatusForbidden) }