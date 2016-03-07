package main

import (
	"testing"
	"github.com/pos/infrastructure"
)
var db infrastructure.Mem_DB

func init() {
	db = infrastructure.NewMemDb()
}

func Test_Return_An_Error_When_ItemId_Does_NOT_Exist (t *testing.T) {

	service := Service{"get_item_service_test", db};
	item := service.GetItem("2")

	if item.Id == "2" {
		t.Fail()
	}
}

func Test_Return_An_ItemId_Just_Saved (t *testing.T) {

	id := "2"
	price := float32(10)
	descr := "milk 100 cm3"

	service := Service{"get_item_service_test", db};

	service.PutItem(id, descr, price)

	item := service.GetItem("2")

	if item.Id != id || item.Desc != descr || item.Price != price {
		t.Fail()
	}
}