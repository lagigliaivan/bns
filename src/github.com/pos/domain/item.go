package domain

type Item struct {
	id string
	desc string
	price float32
	quantity int32
}

func NewItem(id string) *Item {
	p := new(Item)
	p.id = id
	return p
}

func (item *Item) GetId() string {
	return item.id;
}