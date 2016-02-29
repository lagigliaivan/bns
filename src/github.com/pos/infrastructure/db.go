package infrastructure

import (
	"github.com/pos/domain/item"
)
type Db interface {

	Save(item.Item)
	Get(string)
}

