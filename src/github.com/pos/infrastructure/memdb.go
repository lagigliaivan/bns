package infrastructure

import (
	"github.com/pos/dto"
	"sync"
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
	return db.m[id]
}

func (db Mem_DB) SaveItem(item dto.Item) int  {
	db.lock.Lock()
	defer db.lock.Unlock()
	db.m[item.Id] = item
	return 0;
}
