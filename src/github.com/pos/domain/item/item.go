package item

type Item struct {
	id string
	desc string
	price float32
	quantity int32
}

func New(id string) *Item {
	p := new(Item)
	p.id = id
	return p
}