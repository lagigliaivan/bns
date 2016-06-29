package main

import (
	"fmt"
	"time"
)

type DB interface {
	GetItem(string) Item
	SaveItem(Item) int
	GetItems() []Item

	//GetPurchase(time.Time) Purchase
	SavePurchase(Purchase, string) error
	GetPurchases(string) []Purchase
	GetPurchasesGroupedByMonth(string) map[time.Month][]Purchase
	//GetPurchasesByUser(user string)
}

type CatalogDB struct {
	db *sql.DB
}

func (catDb CatalogDB) init() {
	db, err := sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		panic(fmt.Sprintf("Error when opening database connection: %v", err))
	}
	catDb.db = db
}

func (catDb CatalogDB) GetItem(id string) (Item){
	item := Item{}
	item.Id = id
	return item
}

func (catDb CatalogDB) GetItems() []Item{
	return nil
}

func (catDb CatalogDB) SaveItem(Item) int {

	return 0
}

func (catDb CatalogDB) GetPurchase(time time.Time) Purchase  {

	return Purchase{}
}

func (catDb CatalogDB) SavePurchase( p Purchase, userId string) error {


	return nil
}

func (catDb CatalogDB) GetPurchases() []Purchase  {


	return []Purchase{}
}

func (catDb CatalogDB) GetPurchasesGroupedByMonth() map[time.Month][]Purchase  {

	return make(map[time.Month][]Purchase)
}

func (catDb CatalogDB) GetPurchasesByUser(user string) []Purchase  {
	return []Purchase{}
}