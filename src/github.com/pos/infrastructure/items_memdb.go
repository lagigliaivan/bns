package infrastructure

import (
	"github.com/pos/dto"
	"sync"
	"log"
	"time"
)

type items map[string] dto.Item
type purchases map[time.Time] []dto.Purchase

type Mem_DB struct {
	lockI *sync.RWMutex
	lockP  *sync.RWMutex
	items items
	purchases purchases
}

func NewMemDb() (Mem_DB) {
	db := Mem_DB{}
	db.items = make(map[string]dto.Item)
	db.lockP = new(sync.RWMutex)
	db.lockI = new(sync.RWMutex)
	return  db
}

func (db Mem_DB) GetItem(id string) (dto.Item)  {
	db.lockI.Lock()
	defer db.lockI.Unlock()
	log.Printf("GetItem id:%s db.size now: %d", id, len(db.items))
	return db.items[id]
}

func (db Mem_DB) GetItems() ([]dto.Item){

	var items []dto.Item = make([]dto.Item, 0)

	for _, item := range db.items {
		items = append(items, item)
	}

	return items
}

func (db Mem_DB) SaveItem(item dto.Item) int  {
	db.lockI.Lock()
	defer db.lockI.Unlock()
	log.Printf("SaveItem id:%s db.size before: %d", item.Id, len(db.items))
	db.items[item.Id] = item
	log.Printf("SaveItem id:%s db.size now: %d", item.Id, len(db.items))
	return 0;
}

func (db Mem_DB) GetPurchases(time time.Time) []dto.Purchase  {
	db.lockP.Lock()
	defer db.lockP.Unlock()
	purchases := db.purchases[time]
	return purchases
}