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
	"github.com/pos/dto"
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

	item_to_be_added := createItem()

	//GETting URL
	url := getServiceURLToBeTested(ts.URL, item_to_be_added.Id);
	res, err := http.Get(url)
	if (err != nil) || (res.StatusCode != http.StatusNotFound) {
		Log("GET", url, res.StatusCode, http.StatusOK)
		t.Fail()
	}
}

func Test_returns_200_when_item_exists(t *testing.T){

	item_to_be_added := createItem()

	service := NewService(db)
	service.PutItem(item_to_be_added)
	service.GetRequestParameters = getItemIdFromURL(item_to_be_added.Id)

	ts := httptest.NewServer(http.HandlerFunc(service.HandleGetItem))
	defer ts.Close()

	//GETting URL
	url := getServiceURLToBeTested(ts.URL, item_to_be_added.Id);
	res, err := http.Get(url)
	if (err != nil) || (res.StatusCode != http.StatusOK)  {
		Log("GET", url, res.StatusCode, http.StatusOK)
		t.Fail()
	}
}

//TESTs to check PUT /catalog/products/{id} results
func Test_returns_201_when_item_is_created (t *testing.T) {

	item_to_be_added := createItem()
	service := NewService(db)

	ts := httptest.NewServer(http.HandlerFunc(service.HandlePutItem))
	defer ts.Close()

	url := getServiceURLToBeTested(ts.URL, item_to_be_added.Id);
	body_as_string := getRequestBody(item_to_be_added)
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

	item_to_be_added := createItem()

	service := NewService(db)

	service.GetRequestParameters = getItemIdFromURL(item_to_be_added.Id)

	tsPUT := httptest.NewServer(http.HandlerFunc(service.HandlePutItem))
	defer tsPUT.Close()

	url := getServiceURLToBeTested(tsPUT.URL, item_to_be_added.Id);

	body_as_string := getRequestBody(item_to_be_added)

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
	url = getServiceURLToBeTested(tsGET.URL, item_to_be_added.Id);
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

	item_to_be_added := createItem()

	service := NewService(db);
	service.PutItem(item_to_be_added)

	item := service.GetItem(item_to_be_added.Id)

	if item.Id != item_to_be_added.Id || item.Desc != item_to_be_added.Desc || item.Price != item_to_be_added.Price {
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

func createItem() dto.Item {

	id := "2"
	price := float32(10)
	descr := "milk 100 cm3"

	return dto.Item{id, descr, price}
}

func getRequestBody(item dto.Item) string {
	body := "{\"description\":\"" + item.Desc + "\", \"price\":" + fmt.Sprintf("%.2f", item.Price) + "}"
	return body
}

func getItemIdFromURL(item_id string) func (r *http.Request) map[string]string {

	f := func (r *http.Request) map[string]string {
		req_vars := make(map[string]string)
		req_vars["id"]= item_id
		log.Printf("Returning itemid:%s from a testing callback.", req_vars["id"])
		return req_vars
	}

	return f
}

//API to be tested
func getServiceURLToBeTested(base_url, item_id string) string {

	catalog_api := "/catalog/products/"
	return base_url + catalog_api + item_id;
}