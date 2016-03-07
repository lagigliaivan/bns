package infrastructure

import (
	"github.com/pos/domain"
	"database/sql"
	_"github.com/go-sql-driver/mysql"
//	"fmt"
	"fmt"
)

type DB interface {
	GetItem(string) (domain.Item)
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

func (catDb CatalogDB) GetItem(id string) (domain.Item){
	item := domain.NewItem(id)
	return *item
}