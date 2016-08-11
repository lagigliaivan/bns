package main

import (
	"sync"
	"log"
	"time"
	"fmt"
	"strings"
)

type items map[string] Item


type Mem_DB struct {
	lockI *sync.RWMutex
	lockP  *sync.RWMutex
	items items
	purchasesByUser  map[string] map[time.Month] map[string]Purchase
}

func NewMemDb() (Mem_DB) {
	db := Mem_DB{}
	db.items = make(map[string]Item)
	db.purchasesByUser = make (map[string] map[time.Month] map[string]Purchase)
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

func (db Mem_DB) SavePurchase( p Purchase, userId string) error {

	db.lockP.Lock()
	defer db.lockP.Unlock()

	userPurchasesByMonth := db.purchasesByUser[userId]

	if  userPurchasesByMonth == nil {
		userPurchasesByMonth = make(map[time.Month] map[string]Purchase, 0)
	}

	if userPurchasesByMonth[p.Time.Month()] == nil {
		userPurchasesByMonth[p.Time.Month()] = make(map[string]Purchase, 0)
	}
	time := fmt.Sprintf("%d", p.Time.Unix())
	userPurchasesByMonth[p.Time.Month()][time] = p

	db.purchasesByUser[userId] = userPurchasesByMonth


	log.Printf("saving purchase %d", p.Time.Unix())

	return nil
}

func (db Mem_DB) GetPurchases(userId string) []Purchase  {

	purchases := make([]Purchase, 0)

	db.lockP.Lock()
	defer db.lockP.Unlock()

	for _, ps := range db.purchasesByUser[userId] {
		for k := range ps {
			purchases = append(purchases, ps[k])
		}
	}

	return purchases
}

func (db Mem_DB) GetPurchasesByMonth(userId string, year int) map[time.Month] []Purchase  {

	purchases := make(map[time.Month] []Purchase)
	for t, p_by_month := range db.purchasesByUser[userId] {

		for _, p := range p_by_month {
			purchases[t] = append(purchases[t], p)
		}
	}

	return purchases
}

func (db Mem_DB) GetPurchasesByUser(user string) []Purchase  {
	return []Purchase{}
}

func (db Mem_DB) DeletePurchase(userId string, id string) {


	db.lockP.Lock()
	defer db.lockP.Unlock()

	for time, ps_by_month := range db.purchasesByUser[userId] {

		for k, purchase := range ps_by_month {

			if strings.Compare(purchase.Id, id) == 0 {
				log.Printf("Deleting item: %s for user: %s ", id, userId)
				delete(db.purchasesByUser[userId][time], k)
				return
			}
		}
	}

}