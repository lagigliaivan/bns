package dto

import (
	"github.com/pos/domain"
)

type Item struct {
	Id    string
	Desc  string
	Price float32
}

func (itemDto Item) GetDto(item domain.Item) Item {

	itemDto.Id = item.GetId()
	itemDto.Desc = item.GetDescription()
	itemDto.Price = item.GetPrice()

	return itemDto
}