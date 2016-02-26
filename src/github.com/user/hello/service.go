package main

import (
	"fmt"
	"html"
	"net/http"
)

type Service struct {
	name string
}


func (service Service) HandleProducts(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "ProductId: %q %q %+d", service.name, html.EscapeString(r.URL.Path), )
}
func (service Service) HandleRoot(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Hello, %q %q", service.name, html.EscapeString(r.URL.Path))
}