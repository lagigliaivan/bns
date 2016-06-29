package main

import (
	"sync"
	"log"
	"time"
)

type items map[string] Item


type Mem_DB struct {
	lockI *sync.RWMutex
	lockP  *sync.RWMutex
	items items
	/*purchases map[time.Time] Purchase
	purchasesByMonth map[time.Month] []Purchase*/
	purchasesByUser  map[string] map[time.Month] []Purchase
}

func NewMemDb() (Mem_DB) {
	db := Mem_DB{}
	db.items = make(map[string]Item)
	/*db.purchases = make(map[time.Time]Purchase)
	db.purchasesByMonth = make(map[time.Month][]Purchase)*/
	db.purchasesByUser = make (map[string] map[time.Month] []Purchase)
	db.lockP = new(sync.RWMutex)
	db.lockI = new(sync.RWMutex)
	return  db
}

func (db Mem_DB) GetItem(id string) (Item)  {
	db.lockI.Lock()
	defer db.lockI.Unlock()
	log.Printf("GetItem id:%s db.size now: %d", id, len(db.items))
	return db.items[id]
}

func (db Mem_DB) GetItems() ([]Item){

	var items []Item = make([]Item, 0)

	for _, item := range db.items {
		items = append(items, item)
	}

	return items
}

func (db Mem_DB) SaveItem(item Item) int  {
	db.lockI.Lock()
	defer db.lockI.Unlock()
	log.Printf("SaveItem id:%s db.size before: %d", item.Id, len(db.items))
	db.items[item.Id] = item
	log.Printf("SaveItem id:%s db.size now: %d", item.Id, len(db.items))
	return 0;
}
/*
func (db Mem_DB) GetPurchase(time time.Time) Purchase  {
	db.lockP.Lock()
	defer db.lockP.Unlock()
	purchases := db.purchases[time]
	return purchases
}*/

func (db Mem_DB) SavePurchase( p Purchase, userId string) error {

	db.lockP.Lock()
	defer db.lockP.Unlock()

	userPurchasesByMonth := db.purchasesByUser[userId]

	if  userPurchasesByMonth == nil {
		userPurchasesByMonth = make(map[time.Month] []Purchase, 0)
	}

	userPurchasesByMonth[p.Time.Month()] = append(userPurchasesByMonth[p.Time.Month()],p)

	db.purchasesByUser[userId] = userPurchasesByMonth


	return nil
}

func (db Mem_DB) GetPurchases(userId string) []Purchase  {

	purchases := make([]Purchase, 0)

	db.lockP.Lock()
	defer db.lockP.Unlock()

	for _, ps := range db.purchasesByUser[userId] {
		for _, p := range ps {
			purchases = append(purchases, p)
		}
	}

	return purchases
}

func (db Mem_DB) GetPurchasesGroupedByMonth(userId string) map[time.Month][]Purchase  {

	return db.purchasesByUser[userId]
}

func (db Mem_DB) GetPurchasesByUser(user string) []Purchase  {
	return []Purchase{}
}