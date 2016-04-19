package infrastructure

import (
	"database/sql"
	_"github.com/go-sql-driver/mysql"
	"fmt"
	"github.com/pos/dto/item"
	"time"
	"github.com/pos/dto/purchase"

)

type DB interface {
	GetItem(string) item.Item
	SaveItem(item.Item) int
	GetItems() []item.Item

	GetPurchase(time.Time) purchase.Purchase
	SavePurchase(purchase.Purchase) error
	GetPurchases() []purchase.Purchase
	GetPurchasesGroupedByMonth() map[time.Month][]purchase.Purchase
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

func (catDb CatalogDB) GetItem(id string) (item.Item){
	item := item.Item{}
	item.Id = id
	return item
}

func (catDb CatalogDB) GetItems() []item.Item{
	return nil
}

func (catDb CatalogDB) SaveItem(item.Item) int {

	return 0
}

func (catDb CatalogDB) GetPurchase(time time.Time) purchase.Purchase  {

	return purchase.Purchase{}
}

func (catDb CatalogDB) SavePurchase( p purchase.Purchase) error {


	return nil
}

func (catDb CatalogDB) GetPurchases() []purchase.Purchase  {


	return []purchase.Purchase{}
}

func (catDb CatalogDB) GetPurchasesGroupedByMonth() map[time.Month][]purchase.Purchase  {

	return make(map[time.Month][]purchase.Purchase)

}