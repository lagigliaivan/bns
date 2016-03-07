package domain

type Item struct {
	id string
	desc string
	price float32
}

func NewItem(id string) *Item {
	item := new(Item)
	item.id = id
	return item
}

func (item Item) GetId() string {
	return item.id;
}

func (item Item) SetDescription(desc string) {
	item.desc = desc
}

func (item Item) GetDescription() string {
	return item.desc
}

func (item Item) SetPrice(price float32) {
	item.price = price
}

func (item Item) GetPrice() float32 {
	return item.price
}