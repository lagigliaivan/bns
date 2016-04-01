package dto

import (
	"encoding/json"
)

type Item struct {
	Id          string  `json:"id"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
}


func (dto Item) IsEmpty() bool{
	return dto.Id == ""
}

func (dto Item) IsNOTEmpty() bool{
	return !(dto.Id == "")
}

func (dto Item) ToJsonString() string{

	itemAsJson, _ := json.Marshal(dto)
	return string(itemAsJson)
}

type ItemsContainer struct {

	Items []Item  `json:"items"`

}

func NewContainer() ItemsContainer {

	return ItemsContainer{Items: make([]Item, 0)}

}

func (items *ItemsContainer) Add(item Item) {

	items.Items = append(items.Items, item)
}

func (items ItemsContainer) GetItems() []Item{

	return items.Items

}

func (items ItemsContainer) ToJsonString() string{

	itemsAsJson, _ := json.Marshal(items)
	return string(itemsAsJson)

}

func (items ItemsContainer) IsEmpty() bool{

	if len(items.GetItems()) < 1 {

		return true
	}

	return false

}