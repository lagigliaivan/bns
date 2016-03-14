package main

import (
	"testing"
	"github.com/pos/infrastructure"
	"net/http"
	"net/http/httptest"
	"log"
	//"io/ioutil"
)

//DB MOCK
var db infrastructure.Mem_DB

func init() {
	db = infrastructure.NewMemDb()
}

//TESTS

func Test_returns_404_when_item_does_not_exist(t *testing.T){

	service := NewService(db)
	ts := httptest.NewServer(http.HandlerFunc(service.HandleGetItem))
	defer ts.Close()

	url := ts.URL + "/catalog/products/2";

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != http.StatusNotFound {
		t.Fail()
	}
}

func Test_returns_200_when_item_exists(t *testing.T){

	item_id := "123133131"
	item_desc := "item for testing purposes"
	item_price :=  float32(2.0)

	service := NewService(db)
	service.PutItem(item_id, item_desc , item_price)

	ts := httptest.NewServer(http.HandlerFunc(service.HandleGetItem))
	defer ts.Close()

	url := ts.URL + "/catalog/products/" + item_id;

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("GET URL: %s StatusCode different from what was expected: %d <> %d", url, res.StatusCode, http.StatusOK)
		t.Fail()
	}
}

func Test_return_an_error_when_itemid_does_NOT_exist (t *testing.T) {

	service := NewService(db);
	item := service.GetItem("1021")

	if item.Id == "1021" {
		t.Fail()
	}
}

func Test_return_an_itemid_just_saved (t *testing.T) {

	id := "2"
	price := float32(10)
	descr := "milk 100 cm3"

	service := NewService(db);

	service.PutItem(id, descr, price)

	item := service.GetItem("2")

	if item.Id != id || item.Desc != descr || item.Price != price {
		t.Fail()
	}
}

func Test_return_an_empty_item_if_it_does_not_exist(t *testing.T){
	service := NewService(db);
	item := service.GetItem("non_existing_item")

	if item.Id != "" {
		t.Fail()
	}
}