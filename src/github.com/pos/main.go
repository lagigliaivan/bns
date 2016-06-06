/**
This Package starts up a server which has the following APIs:
GET /items
GET /purchases
**/
package main

import (
	"log"
	"net/http"
	"github.com/pos/infrastructure"
	"github.com/pos/services"
)

func main() {

	router := services.NewRouter()
	db := infrastructure.NewMemDb()

	itemsService := services.NewItemService(db)
	itemsService.ConfigureRouter(router)

	purchasesService := services.NewPurchaseService(db)
	purchasesService.ConfigureRouter(router)

	log.Fatal(http.ListenAndServe(":8080", router))
}