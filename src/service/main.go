/**
This Package starts up a server which has the following APIs:
GET /items
GET /purchases
**/
package main

import (
	"log"
	"net/http"
)

func main() {

	router := NewRouter()
	db := NewMemDb()

	itemsService := NewItemService(db)
	itemsService.ConfigureRouter(router)

	purchasesService := NewPurchaseService(db)
	purchasesService.ConfigureRouter(router)

	log.Fatal(http.ListenAndServe(":8080", router))
}