package main

import (
	"testing"
	"github.com/pos/infrastructure"
	"net/http"
	"net/http/httptest"
	"log"
	"strings"
	"bytes"
	"fmt"
)

//DB MOCK
var db infrastructure.Mem_DB

func init() {
	db = infrastructure.NewMemDb()
}

//TESTs to check GET /catalog/products/{id} results
func Test_returns_404_when_item_does_not_exist(t *testing.T){

	service := NewService(db)
	ts := httptest.NewServer(http.HandlerFunc(service.HandleGetItem))
	defer ts.Close()

	item_id := "2"
	url := ts.URL + "/catalog/products/" + item_id;

	res, err := http.Get(url)
	if (err != nil) || (res.StatusCode != http.StatusNotFound) {
		Log("GET", url, res.StatusCode, http.StatusOK)
		t.Fail()
	}

}

func Test_returns_200_when_item_exists(t *testing.T){

	item_id := "123133131"
	item_desc := "item for testing purposes"
	item_price :=  float32(2.0)

	service := NewService(db)
	service.PutItem(item_id, item_desc , item_price)

	f := func (r *http.Request) map[string]string {
		req_vars := make(map[string]string)
		req_vars["id"]= item_id
		log.Printf("Returning itemid:%s from a testing callback.", req_vars["id"])
		return req_vars
	}
	service.GetRequestParameters = f

	ts := httptest.NewServer(http.HandlerFunc(service.HandleGetItem))
	defer ts.Close()

	url := ts.URL + "/catalog/products/" + item_id;
	res, err := http.Get(url)
	if (err != nil) || (res.StatusCode != http.StatusOK)  {
		Log("GET", url, res.StatusCode, http.StatusOK)
		t.Fail()
	}
}

//TESTs to check PUT /catalog/products/{id} results
func Test_returns_201_when_item_is_created (t *testing.T) {

	item_id := "123133131"
	item_desc := "item for testing purposes"
	item_price :=  "2.0"

	service := NewService(db)

	ts := httptest.NewServer(http.HandlerFunc(service.HandlePutItem))
	defer ts.Close()

	url := ts.URL + "/catalog/products/" + item_id;
	body_as_string := "{\"description\":\"" + item_desc + "\", \"price\":" + item_price + "}"
	body := strings.NewReader(body_as_string)
	req,err := http.NewRequest("PUT", url, body)
	if err != nil {
		log.Printf("Error when creating PUT request %d.", err)
		t.Fail()
	}

	res, err := http.DefaultClient.Do(req)


	if (err != nil) || (res.StatusCode != http.StatusCreated)  {
		Log("PUT", url, res.StatusCode, http.StatusCreated)
		t.Fail()
	}
}

func Test_returns_the_item_after_it_is_created(t *testing.T){
	item_id := "112"
	item_desc := "item for testing purposes"
	item_price :=  "2.0"

	service := NewService(db)

	f := func (r *http.Request) map[string]string {
		req_vars := make(map[string]string)
		req_vars["id"]= item_id
		log.Printf("Returning itemid:%s from a testing callback.", req_vars["id"])
		return req_vars
	}
	service.GetRequestParameters = f

	tsPUT := httptest.NewServer(http.HandlerFunc(service.HandlePutItem))
	defer tsPUT.Close()

	url := tsPUT.URL + "/catalog/products/" + item_id;
	body_as_string := "{\"description\":\"" + item_desc + "\", \"price\":" + item_price + "}"
	body := strings.NewReader(body_as_string)
	req,err := http.NewRequest("PUT", url, body)
	if err != nil {
		log.Printf("Error when creating PUT request %d.", err)
		t.Fail()
	}
	res, err := http.DefaultClient.Do(req)

	if (err != nil) || (res.StatusCode != http.StatusCreated)  {
		Log("PUT", url, res.StatusCode, http.StatusCreated)
		t.Fail()
	}

	tsGET := httptest.NewServer(http.HandlerFunc(service.HandleGetItem))
	defer tsGET.Close()
	url = tsGET.URL + "/catalog/products/" + item_id;
	res, err = http.Get(url)
	if (err != nil) || (res.StatusCode != http.StatusOK)  {
		Log("GET", url, res.StatusCode, http.StatusOK)
		t.Fail()
	}
}

func Test_returns_an_error_when_item_does_NOT_exist (t *testing.T) {

	service := NewService(db);
	item := service.GetItem("1021")

	if item.Id == "1021" {
		t.Fail()
	}
}

func Test_return_an_item_just_saved (t *testing.T) {

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

func Log(method string, url string, expectedStatusCode int, receivedStatusCode int){

	var buf bytes.Buffer
	logger := log.New(&buf, "logger: ", log.Lshortfile)
	logger.Printf("%s URL: %s StatusCode %d different from what was expected %d", method, url, expectedStatusCode, receivedStatusCode)
	fmt.Print(&buf)
}