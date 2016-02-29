package main

import (
	"fmt"
	"html"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/pos/domain/item"
)

type Service struct {
	name string
}


func (service Service) HandleProducts(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	prodId := vars["id"]

	item :=  item.New(prodId)



	fmt.Fprintf(w, "ProductId: %q %q %s", service.name, html.EscapeString(r.URL.Path), prodId)
}
func (service Service) HandleRoot(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Hello, %q %q", service.name, html.EscapeString(r.URL.Path))
}