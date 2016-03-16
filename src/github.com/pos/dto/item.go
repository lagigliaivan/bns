package dto

import (
	"github.com/pos/domain"
)

type Item struct {
	Id    string
	Desc  string
	Price float32
}

func (dto Item) GetDto(item domain.Item) Item {

	dto.Id = item.GetId()
	dto.Desc = item.GetDescription()
	dto.Price = item.GetPrice()

	return dto
}

func (dto Item) IsEmpty() bool{
	return dto.Id == ""
}

func (dto Item) IsNOTEmpty() bool{
	return !(dto.Id == "")
}