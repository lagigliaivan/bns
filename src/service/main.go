/**
This Package starts up a server which has the following APIs:
GET /items
GET /purchases
**/
package main

import (
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	BNS_DB = "BNS_DB"
	LOCALDB = "LOCALDB"
	MEMDB = "MEMDB"
)

func main() {

	dbType := os.Getenv("BNS_DB")


	var db DB

	if strings.Compare(dbType,LOCALDB) == 0 {

		log.Print("Using LOCALDB")
		db, _ = NewDynamoDB("http://localhost:8000", "us-west-2")

	} else if strings.Compare(dbType, MEMDB) == 0 {

		db = NewMemDb()
		log.Print("Using MEMDB")
	} else {

		db, _ = NewDynamoDB("", "us-west-2")
		log.Print("Using DYNAMODB")
	}

	router := NewRouter()

	/*
	itemsService := NewItemService(db)
	itemsService.ConfigureRouter(router)
	*/

	purchasesService := NewPurchaseService(db)
	purchasesService.ConfigureRouter(router.GetRouter())

	log.Fatal(http.ListenAndServe(":8080", router))
}