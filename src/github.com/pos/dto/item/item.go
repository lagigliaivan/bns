package item

import (
	"encoding/json"
)

type Item struct {
	Id          string  `json:"id"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
}


func (item Item) IsEmpty() bool{
	return item.Id == ""
}

func (item Item) IsNOTEmpty() bool{
	return !(item.Id == "")
}

func (item Item) ToJsonString() string{

	itemAsJson, _ := json.Marshal(item)
	return string(itemAsJson)
}

type Container struct {

	Items []Item  `json:"items"`

}

func NewContainer() Container {

	return Container{Items: make([]Item, 0)}

}

func (items *Container) Add(item Item) {

	items.Items = append(items.Items, item)
}

func (items Container) GetItems() []Item{

	return items.Items

}

func (items Container) ToJsonString() string{

	itemsAsJson, _ := json.Marshal(items)
	return string(itemsAsJson)

}

func (items Container) IsEmpty() bool{

	if len(items.GetItems()) < 1 {

		return true
	}

	return false

}