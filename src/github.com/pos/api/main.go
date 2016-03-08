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
	service := NewService(db)

	router.HandleFunc("/catalog/products/{id}", service.HandleRequest)

	router.Methods("GET", "PUT")

	log.Fatal(http.ListenAndServe(":8080", router))
}

