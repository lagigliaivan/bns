package main

import (
	"testing"
	"github.com/pos/domain"
)

type DB_MOCK struct {

}

func (db DB_MOCK) GetItem(string) (domain.Item)  {
	item := domain.NewItem("1")
   	return *item
}

func Test_Return_An_Existing_ItemId(t *testing.T) {

	service := Service{"get_item_service_test", DB_MOCK{}};
	item := service.GetItem("1")

	if item.GetId() != "1" {
		t.Fail()
	}
}

func Test_Return_An_Error_When_ItemId_Does_NOT_Exist (t *testing.T) {

	service := Service{"get_item_service_test", DB_MOCK{}};
	item := service.GetItem("2")

	if item.GetId() == "2" {
		t.Fail()
	}
}

func Test_Return_An_ItemId_Just_Saved (t *testing.T) {


	id := "1"
	price := float32(10)
	descr := "milk 100 cm3"

	service := Service{"get_item_service_test", DB_MOCK{}};

	service.PutItem(id, price, descr)


	item := service.GetItem("2")

	if item.GetId() == "2" {
		t.Fail()
	}
}