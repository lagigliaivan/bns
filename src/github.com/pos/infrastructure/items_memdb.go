package infrastructure

import (
	"github.com/pos/dto/item"
	"github.com/pos/dto/purchase"
	"sync"
	"log"
	"time"
)

type items map[string] item.Item
type purchases map[time.Time] purchase.Purchase

type Mem_DB struct {
	lockI *sync.RWMutex
	lockP  *sync.RWMutex
	items items
	purchases purchases
}

func NewMemDb() (Mem_DB) {
	db := Mem_DB{}
	db.items = make(map[string]item.Item)
	db.purchases = make(map[time.Time]purchase.Purchase)
	db.lockP = new(sync.RWMutex)
	db.lockI = new(sync.RWMutex)
	return  db
}

func (db Mem_DB) GetItem(id string) (item.Item)  {
	db.lockI.Lock()
	defer db.lockI.Unlock()
	log.Printf("GetItem id:%s db.size now: %d", id, len(db.items))
	return db.items[id]
}

func (db Mem_DB) GetItems() ([]item.Item){

	var items []item.Item = make([]item.Item, 0)

	for _, item := range db.items {
		items = append(items, item)
	}

	return items
}

func (db Mem_DB) SaveItem(item item.Item) int  {
	db.lockI.Lock()
	defer db.lockI.Unlock()
	log.Printf("SaveItem id:%s db.size before: %d", item.Id, len(db.items))
	db.items[item.Id] = item
	log.Printf("SaveItem id:%s db.size now: %d", item.Id, len(db.items))
	return 0;
}

func (db Mem_DB) GetPurchase(time time.Time) purchase.Purchase  {
	db.lockP.Lock()
	defer db.lockP.Unlock()
	purchases := db.purchases[time]
	return purchases
}

func (db Mem_DB) SavePurchase( purchase purchase.Purchase) error {

	db.lockP.Lock()
	defer db.lockP.Unlock()
	db.purchases[purchase.Time] = purchase
	log.Printf("SavePurchase time: %s", purchase.Time)
	return nil
}

func (db Mem_DB) GetPurchases() []purchase.Purchase  {

	purchases := make([]purchase.Purchase, 0)

	db.lockP.Lock()
	defer db.lockP.Unlock()

	for _, p := range db.purchases {
		purchases = append(purchases, p)
	}

	return purchases
}