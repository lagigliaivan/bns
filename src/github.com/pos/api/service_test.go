package main

import (
	"log"
	"testing"
	"github.com/pos/infrastructure"
	"net/http"
	"net/http/httptest"
	"strings"
	"bytes"
	"fmt"
	"github.com/pos/dto"
	"io/ioutil"
	"encoding/json"
//	"reflect"
	"io"
)

//DB MOCK
var db infrastructure.Mem_DB


func init() {
	db = infrastructure.NewMemDb()
}

//Testing service to check GET /catalog/productestingServer/{id}
func Test_returns_404_when_item_does_not_exist(t *testing.T){

	service := NewService(db)
	testingServer := httptest.NewServer(http.HandlerFunc(service.HandleGetItem))
	defer testingServer.Close()

	itemToBeAdded := createItemDto()

	//GETting URL
	url := getURLToBeTested(testingServer.URL, itemToBeAdded.Id);

	res, err := http.Get(url)
	if !isHTTPStatus(http.StatusNotFound, res, err){
		debug("GET", url, res.StatusCode, http.StatusOK)
		t.Fail()
	}
}

func Test_returns_200_when_item_exists(t *testing.T){

	itemToBeAdded := createItemDto()

	service := NewService(db)
	service.PutItem(itemToBeAdded)
	service.GetRequestParameters = returnItemIdFromURL(itemToBeAdded.Id)

	testingServer := httptest.NewServer(http.HandlerFunc(service.HandleGetItem))
	defer testingServer.Close()

	//GETting URL
	url := getURLToBeTested(testingServer.URL, itemToBeAdded.Id);

	res, err := http.Get(url)
	if !isHTTPStatus(http.StatusOK, res, err){
		debug("GET", url, res.StatusCode, http.StatusOK)
		t.Fail()
	}
}

//TEStestingServer to check PUT /catalog/productestingServer/{id} resultestingServer
func Test_returns_201_when_item_is_created (t *testing.T) {

	itemToBeAdded := createItemDto()
	service := NewService(db)

	testingServer := httptest.NewServer(http.HandlerFunc(service.HandlePutItem))
	defer testingServer.Close()
	url := getURLToBeTested(testingServer.URL, itemToBeAdded.Id);

	res, err := httpPut(url, itemToBeAdded)

	if !isHTTPStatus(http.StatusCreated, res, err){
		debug("PUT", url, res.StatusCode, http.StatusCreated)
		t.Fail()
	}
}

func Test_returns_200_when_item_is_updated (t *testing.T) {

	itemToBeAdded := createItemDto()
	service := NewService(db)

	testingServer := httptest.NewServer(http.HandlerFunc(service.HandlePutItem))
	defer testingServer.Close()
	url := getURLToBeTested(testingServer.URL, itemToBeAdded.Id);

	res, err := httpPut(url, itemToBeAdded)

	if !isHTTPStatus(http.StatusOK, res, err){
		debug("PUT", url, res.StatusCode, http.StatusOK)
		t.Fail()
	}
	//GET Item
	testingServerGET := httptest.NewServer(http.HandlerFunc(service.HandleGetItem))
	defer testingServerGET.Close()

	res, err = http.Get(url)

	if !isHTTPStatus(http.StatusOK, res, err){
		debug("GET", url, res.StatusCode, http.StatusOK)
		t.Fail()
	}

	if !areItemsEquals(itemToBeAdded, createItemFromJson(res.Body)){
		log.Printf("Error when GETting item to contrast it with the saved one")
		t.Fail()
	}


}

func Test_returns_the_same_item_after_it_is_created(t *testing.T){

	itemToBeAdded := createItemDto()

	service := NewService(db)
	service.GetRequestParameters = returnItemIdFromURL(itemToBeAdded.Id)

	testingServerPUT := httptest.NewServer(http.HandlerFunc(service.HandlePutItem))
	defer testingServerPUT.Close()

	//PUT Item
	url := getURLToBeTested(testingServerPUT.URL, itemToBeAdded.Id);

	res, err := httpPut(url, itemToBeAdded)

	if !isHTTPStatus(http.StatusCreated, res, err){
		debug("PUT", url, res.StatusCode, http.StatusCreated)
		t.Fail()
	}

	//GET Item
	testingServerGET := httptest.NewServer(http.HandlerFunc(service.HandleGetItem))
	defer testingServerGET.Close()
	url = getURLToBeTested(testingServerGET.URL, itemToBeAdded.Id);

	res, err = http.Get(url)

	if !isHTTPStatus(http.StatusOK, res, err){
		debug("GET", url, res.StatusCode, http.StatusOK)
		t.Fail()
	}

	if !areItemsEquals(itemToBeAdded, createItemFromJson(res.Body)){
		log.Printf("Error when GETting item to contrast it with the saved one")
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

func Test_returns_an_item_just_saved (t *testing.T) {

	itemToBeAdded := createItemDto()

	service := NewService(db);
	service.PutItem(itemToBeAdded)

	item := service.GetItem(itemToBeAdded.Id)

	if item.Id != itemToBeAdded.Id || item.Desc != itemToBeAdded.Desc || item.Price != itemToBeAdded.Price {
		t.Fail()
	}
}

func Test_returns_an_empty_item_if_it_does_not_exist(t *testing.T){

	service := NewService(db);
	item := service.GetItem("non_existing_item")

	if item.Id != "" {
		t.Fail()
	}
}

func Test_returns_no_error_when_adding_an_item(t *testing.T){

	service := NewService(db);
	item := createItemDto()
	err := service.PutItem(item)

	if err != 0 {
		t.Fail()
	}
}

//Test auxiliary functions
func debug(method string, url string, expectedStatusCode int, receivedStatusCode int){

	var buf bytes.Buffer
	logger := log.New(&buf, "logger: ", log.Lshortfile)
	logger.Printf("%s URL: %s StatusCode %d different from what was expected %d", method, url, expectedStatusCode, receivedStatusCode)
	fmt.Print(&buf)
}

func createItemDto() dto.Item {

	id := "2"
	price := float32(10)
	descr := "milk 100 cm3"

	return dto.Item{id, descr, price}
}

func getRequestBody(item dto.Item) string {
	body := "{\"desc\":\"" + item.Desc + "\", \"price\":" + fmt.Sprintf("%.2f", item.Price) + "}"
	return body
}

func returnItemIdFromURL(item_id string) func (r *http.Request) map[string]string {

	f := func (r *http.Request) map[string]string {
		req_vars := make(map[string]string)
		req_vars["id"]= item_id
		log.Printf("Returning itemid:%s from a testing callback.", req_vars["id"])
		return req_vars
	}

	return f
}

//API to be tested
func getURLToBeTested(base_url, item_id string) string {

	catalog_api := "/catalog/productestingServer/"
	return base_url + catalog_api + item_id;
}

func httpPut(url string, item dto.Item) (resp * http.Response, err error) {

	bodyAsString := getRequestBody(item)
	body := strings.NewReader(bodyAsString)
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		log.Printf("Error when creating PUT request %d.", err)
		return nil, err
	}
	resp, err = http.DefaultClient.Do(req)
	return resp, err
}


func isHTTPStatus(httpStatus int, res *http.Response, err error ) bool {
	return !( (err != nil) || (res.StatusCode != httpStatus) )
}

func areItemsEquals(item, item2 dto.Item) bool{

	return !((item.Id != item2.Id) || (item.Desc != item2.Desc) || (item.Price != item2.Price))
}

func createItemFromJson(itemAsJson io.ReadCloser) dto.Item {

	item := new(dto.Item)
	response, err := ioutil.ReadAll(itemAsJson)

	if err != nil {
		log.Printf("Error when reading Json from response")
	}

	if err := json.Unmarshal(response, item); err != nil {
		log.Printf("Error when unmarshaling Json to dto.Item")
	}

	return *item
}