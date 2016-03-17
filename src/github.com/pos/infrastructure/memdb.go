package infrastructure

import (
	"github.com/pos/dto"
	"sync"
	"log"
)

type Mem_DB struct {
	lock *sync.RWMutex
	m map[string]dto.Item
}

func NewMemDb() (Mem_DB) {
	db := Mem_DB{}
	db.m = make(map[string]dto.Item)
	db.lock = new(sync.RWMutex)
	return  db
}

func (db Mem_DB) GetItem(id string) (dto.Item)  {
	db.lock.Lock()
	defer db.lock.Unlock()
	log.Printf("GetItem id:%s db.size now: %d", id, len(db.m))
	return db.m[id]
}

func (db Mem_DB) SaveItem(item dto.Item) int  {
	db.lock.Lock()
	defer db.lock.Unlock()
	log.Printf("SaveItem id:%s db.size before: %d", item.Id, len(db.m))
	db.m[item.Id] = item
	log.Printf("SaveItem id:%s db.size now: %d", item.Id, len(db.m))
	return 0;
}
