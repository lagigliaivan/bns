package infrastructure

import (
	"database/sql"
	_"github.com/go-sql-driver/mysql"
	"fmt"
	"github.com/pos/dto"
)

type DB interface {
	GetItem(string) (dto.Item)
	SaveItem(dto.Item) int
}

type CatalogDB struct {
	db *sql.DB
}

func (catDb *CatalogDB) init() {
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

func (catDb CatalogDB) SaveItem(dto.Item) int {

	return 0
}
