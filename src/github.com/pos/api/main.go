package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/pos/infrastructure"
)

func main() {

	router := mux.NewRouter();
	db := infrastructure.CatalogDB{}
	root := Service{"api - > root", db}
	foo := Service{"api -> foo", db}

	router.HandleFunc("/", root.HandleRoot)
	router.HandleFunc("/catalog/products/{id}", foo.HandleProducts)

	log.Fatal(http.ListenAndServe(":8080", router))
}

