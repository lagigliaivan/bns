package infrastructure

import (
	"database/sql"
	_"github.com/go-sql-driver/mysql"
	"fmt"
	"github.com/pos/dto"
	"time"

)

type DB interface {
	GetItem(string) dto.Item
	SaveItem(dto.Item) int
	GetItems() []dto.Item

	GetPurchase(time.Time) dto.Purchase
	SavePurchase(dto.Purchase) error
	GetPurchases() []dto.Purchase
	GetPurchasesGroupedByMonth() map[time.Month][]dto.Purchase
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

func (catDb CatalogDB) GetItem(id string) (dto.Item){
	item := dto.Item{}
	item.Id = id
	return item
}

func (catDb CatalogDB) GetItems() []dto.Item{
	return nil
}

func (catDb CatalogDB) SaveItem(dto.Item) int {

	return 0
}

func (catDb CatalogDB) GetPurchase(time time.Time) dto.Purchase  {

	return dto.Purchase{}
}

func (catDb CatalogDB) SavePurchase( p dto.Purchase) error {


	return nil
}

func (catDb CatalogDB) GetPurchases() []dto.Purchase  {


	return []dto.Purchase{}
}

func (catDb CatalogDB) GetPurchasesGroupedByMonth() map[time.Month][]dto.Purchase  {

	return make(map[time.Month][]dto.Purchase)

}