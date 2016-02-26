package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter();
	root := Service{"service - > root"}
	foo := Service{"service -> foo"}

	router.HandleFunc("/", root.HandleRoot)
	router.HandleFunc("/catalog/products/{id}", foo.HandleProducts)

	log.Fatal(http.ListenAndServe(":8080", router))
}

