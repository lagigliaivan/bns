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
	"io"
	"reflect"

)


func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

var (
	setOfItems = []dto.Item{

		{
			Id: "1",
			Description: "first product",
			Price: 2.0,
		},
		{
			Id: "2",
			Description: "second product",
			Price: 34.0,
		},
		{
			Id: "3",
			Description: "third product",
			Price: 332.0,
		},
		{
			Id: "4",
			Description: "forth product",
			Price: 22.0,
		},
	}

	postItems = dto.ItemContainer{Items:setOfItems}
)

//Testing service to check GET /catalog/product/{id}
func Test_GET_item_returns_404_when_it_does_not_exist(t *testing.T){

	testingServer := getServer(NewItemService(infrastructure.NewMemDb()))
	defer testingServer.Close()

	itemToBeAdded := createItemDto()

	//GETting URL
	url := getURLToBeTested(testingServer.URL, itemToBeAdded.Id);

	res, err := httpGet("lagigliaivan@gmail.com", url)
	if !isHTTPStatus(http.StatusNotFound, res, err){
		deb("GET", url, res.StatusCode, http.StatusOK)
		t.FailNow()
	}
}

func Test_GET_item_returns_200_when_it_exists(t *testing.T){

	itemToBeAdded := createItemDto()
	service := NewItemService(infrastructure.NewMemDb())

	//Adding ITEM without calling RESTapi. Calling a service function directly
	service.addUpdateItem(itemToBeAdded)

	router := NewRouter();
	service.ConfigureRouter(router)

	testingServer := httptest.NewServer(router)
	defer testingServer.Close()

	//GETting URL
	url := getURLToBeTested(testingServer.URL, itemToBeAdded.Id);
	log.Printf("url:%s\n", url)
	res, err := httpGet("lagigliaivan@gmail.com", url)
	if !isHTTPStatus(http.StatusOK, res, err){
		deb("GET", url, res.StatusCode, http.StatusOK)
		t.FailNow()
	}
}

//Testing server to check POST /catalog/product/{id}
func Test_POST_item_returns_201_when_it_is_successfully_created (t *testing.T) {

	testingServer := getServer(NewItemService(infrastructure.NewMemDb()))

	defer testingServer.Close()

	itemToBeAdded := createItemDto()
	items := dto.NewItemContainer()
	items.Add(itemToBeAdded)

	url := getURLToBeTested(testingServer.URL);
	res, err := httpPost("121aseda2123123", strings.TrimSuffix(url, "/"), items)

	if !isHTTPStatus(http.StatusCreated, res, err){
		deb(http.MethodPost, url, res.StatusCode, http.StatusCreated)
		t.FailNow()
	}
}

func Test_POST_GET_returns_the_same_item_after_it_is_created(t *testing.T){

	itemToBeAdded := createItemDto()

	server := getServer(NewItemService(infrastructure.NewMemDb()))
	defer server.Close()

	//POST Item
	url := getURLToBeTested(server.URL);


	items := dto.NewItemContainer()
	items.Add(itemToBeAdded)
	res, err := httpPost("abafadfaf9a9fa0fa", strings.TrimSuffix(url, "/"), items)

	if !isHTTPStatus(http.StatusCreated, res, err){
		deb(http.MethodPost, url, res.StatusCode, http.StatusCreated)
		t.FailNow()
	}

	//GET Item

	url = getURLToBeTested(server.URL, itemToBeAdded.Id);

	res, err = httpGet("abafadfaf9a9fa0fa", url)

	if !isHTTPStatus(http.StatusOK, res, err){
		deb(http.MethodGet, url, res.StatusCode, http.StatusOK)
		t.FailNow()
	}

	if !reflect.DeepEqual(itemToBeAdded, createItemFromJson(res.Body)) {
		log.Printf("Error when GETting item to contrast it with the saved one")
		t.FailNow()
	}

}

func Test_PUT_item_returns_200_when_it_is_successfully_updated (t *testing.T) {

	//POST Item
	itemToBeAdded := createItemDto()

	testingServer := getServer(NewItemService(infrastructure.NewMemDb()))
	defer testingServer.Close()
	url := getURLToBeTested(testingServer.URL);

	items := dto.NewItemContainer()
	items.Add(itemToBeAdded)

	res, err := httpPost("lagigliaiv@gmail.com.ar", strings.TrimSuffix(url, "/"), items)

	if !isHTTPStatus(http.StatusCreated, res, err){
		deb(http.MethodPut, url, res.StatusCode, http.StatusCreated)
		t.FailNow()
	}

	//PUT Item
	url = getURLToBeTested(testingServer.URL, itemToBeAdded.Id);

	itemToBeAdded.Description = "Description updated"
	itemToBeAdded.Price = float32(21)

	res, err = httpPut("lagigliaiv@gmail.com.ar", url, itemToBeAdded)

	if !isHTTPStatus(http.StatusOK, res, err){
		deb(http.MethodPut, url, res.StatusCode, http.StatusOK)
		t.FailNow()
	}

	//GET Item
	res, err = httpGet("lagigliaiv@gmail.com.ar", url)

	if !isHTTPStatus(http.StatusOK, res, err){
		deb(http.MethodGet, url, res.StatusCode, http.StatusOK)
		t.FailNow()
	}

	if !reflect.DeepEqual(itemToBeAdded, createItemFromJson(res.Body)) {
		log.Printf("Error when GETting item to contrast it with the saved one")
		t.FailNow()
	}
}

func Test_POST_item_returns_400_when_body_is_sent_without_item_id(t *testing.T){

	itemToBeAdded := createItemDto()
	itemToBeAdded.Id = ""

	server := getServer(NewItemService(infrastructure.NewMemDb()))
	defer server.Close()

	//POST Item
	url := getURLToBeTested(server.URL);

	items := dto.NewItemContainer()
	items.Add(itemToBeAdded)

	res, err := httpPost("lagigliaiv@gmail.com.ar", strings.TrimSuffix(url, "/"), items)

	if !isHTTPStatus(http.StatusBadRequest, res, err){
		deb(http.MethodPost, url, res.StatusCode, http.StatusBadRequest)
		t.FailNow()
	}
}

func Test_GET_items_returns_a_list_of_items(t *testing.T){

	server := getServer(NewItemService(infrastructure.NewMemDb()))
	defer server.Close()

	if httpPOST("lagigliaiv@gmail.com.ar", *server) != nil {
		t.FailNow();
	}

	url := getURLToBeTested(server.URL);

	defer server.Close()
	log.Printf("server.URL:%s", url)

	res, err := httpGet("lagigliaiv@gmail.com.ar", url)

	if err != nil {
		log.Printf("ERROR")
		t.FailNow()
	}

	items := dto.NewItemContainer()

	body, err := ioutil.ReadAll(res.Body)

	if err := json.Unmarshal(body, &items); err != nil {

		log.Printf("Error when reading response %s", err)
		t.FailNow()
	}

	if len(items.GetItems()) != len(setOfItems){
		log.Printf("Error: Expected items quantity is different from the received one")
		t.FailNow()

	}

	itemsFound := 0

	for _, i := range items.GetItems() {

		for _, x := range setOfItems {
			if x.Id == i.Id && reflect.DeepEqual(x, i){
				itemsFound++
				break
			}
		}
	}

	if itemsFound != len(setOfItems) {
		log.Printf("Error: Some items are missing")
		t.FailNow()
	}

}
/*
func Test_POSTing_Date_Is_Stored_When_Item_Is_Saved(t *testing.T) {

	server := getServer(NewItemService(infrastructure.NewMemDb()))
	defer server.Close()

	if httpPOST(*server) != nil {
		t.FailNow();
	}

}*/

//Testing functions that are not exposed as REST services.
func Test_returns_an_error_when_item_does_NOT_exist (t *testing.T) {

	service := NewItemService(infrastructure.NewMemDb());
	item := service.getItem("1021")

	if item.Id == "1021" {
		t.FailNow()
	}
}

func Test_returns_an_item_just_saved (t *testing.T) {

	itemToBeAdded := createItemDto()

	service := NewItemService(infrastructure.NewMemDb());
	service.addUpdateItem(itemToBeAdded)

	item := service.getItem(itemToBeAdded.Id)

	if item.Id != itemToBeAdded.Id || item.Description != itemToBeAdded.Description || item.Price != itemToBeAdded.Price {
		t.FailNow()
	}
}

func Test_returns_an_empty_item_if_it_does_not_exist(t *testing.T){

	service := NewItemService(infrastructure.NewMemDb());
	item := service.getItem("non_existing_item")

	if item.Id != "" {
		t.FailNow()
	}
}

func Test_returns_no_error_when_adding_an_item(t *testing.T){

	service := NewItemService(infrastructure.NewMemDb());
	item := createItemDto()
	err := service.addUpdateItem(item)

	if err != 0 {
		t.FailNow()
	}
}

//Tests auxiliary functions

func deb(method string, url string, expectedStatusCode int, receivedStatusCode int){

	var buf bytes.Buffer
	logger := log.New(&buf, "logger: ", log.Lshortfile)
	logger.Printf("%s URL: %s StatusCode %d different from what was expected %d", method, url, expectedStatusCode, receivedStatusCode)
	fmt.Print(&buf)
}

func createItemDto() dto.Item {

	id := "12345"
	price := float32(10.1)
	descr := "milk 100 cm3"

	return dto.Item{id, descr, price}
}

//API to be tested
func getURLToBeTested(base_url string, params ... string) string {

	var p string

	for _, v := range params {
		p = v + "/"
	}

	catalog_api := "/catalog/products/"

	p = strings.TrimSuffix(p, "/")

	return base_url + catalog_api + p;
}

func httpPOST(user string, server httptest.Server) error{

	//POST Items for later being retrieved.
	url := getURLToBeTested(server.URL);

	res, err := httpPost(user, strings.TrimSuffix(url, "/"), postItems)

	if !isHTTPStatus(http.StatusCreated, res, err){
		deb(http.MethodPost, url, res.StatusCode, http.StatusCreated)
		return err
	}

	return nil
}

func httpPut(user, url string, item dto.Stringifiable) (resp * http.Response, err error) {

	bodyAsString := item.ToJsonString()
	log.Printf("body: %s", bodyAsString)
	body := strings.NewReader(bodyAsString)

	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		log.Printf("Error when creating PUT request %d.", err)
		return nil, err
	}
	req.Header.Add(HEADER, user)
	resp, err = http.DefaultClient.Do(req)
	return resp, err
}

func httpGet(user, url string) (*http.Response, error){

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("Error when creating PUT request %d.", err)
		return nil, err
	}
	req.Header.Add(HEADER, user)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error when creating PUT request %d.", err)
		return nil, err
	}
	return resp, err
}

func httpPost(user, url string, values dto.Stringifiable) (*http.Response, error){

	body := strings.NewReader(values.ToJsonString())
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		log.Printf("Error when creating POST request %d.", err)
		return nil, err
	}
	req.Header.Add(HEADER, user)
	resp, err := http.DefaultClient.Do(req)

	return resp, err
}

func getServer(service Service) *httptest.Server {

	router := NewRouter()
	service.ConfigureRouter(router)
	server := httptest.NewServer(router)

	return server
}

func isHTTPStatus(httpStatus int, res *http.Response, err error ) bool {
	return !( (err != nil) || (res.StatusCode != httpStatus) )
}

func createItemFromJson(itemAsJson io.ReadCloser) dto.Item {

	item := new(dto.Item)
	response, err := ioutil.ReadAll(itemAsJson)

	if err != nil {
		log.Printf("Error when reading Json from response")
	}

	if err := json.Unmarshal(response, item); err != nil {
		log.Printf("Error when unmarshaling Json to item.Item")
	}

	return *item
}



var user1 string
var user2 string

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	sha := sha1.New()
	io.WriteString(sha, "mayname:password@gmail.com.ar")
	user1 = fmt.Sprintf("%x", sha.Sum(nil))

	sha.Reset()
	io.WriteString(sha, "mayname2:password@gmail.com.ar")
	user2 = fmt.Sprintf("%x", sha.Sum(nil))
}


const STATUS_ERROR_MESSAGE string = "%s %s Received status code: %d different from what was expected: %d"

var (
	itemsPurchaseA = []dto.Item{

		{
			Id: "1",
			Description: "first product",
			Price: 2.0,
		},
		{
			Id: "2",
			Description: "second product",
			Price: 34.0,
		},
		{
			Id: "3",
			Description: "third product",
			Price: 332.0,
		},
		{
			Id: "4",
			Description: "forth product",
			Price: 22.0,
		},
	}

	itemsPurchaseB = []dto.Item{

		{
			Id: "100",
			Description: "first product",
			Price: 122.0,
		},
		{
			Id: "200",
			Description: "second product",
			Price: 314.0,
		},
		{
			Id: "300",
			Description: "third product",
			Price: 3212.0,
		},
		{
			Id: "400",
			Description: "forth product",
			Price: 212.0,
		},
	}

	setOfPurchases = []dto.Purchase{

		{
			Time: time.Now(),
			Point:dto.NewPoint(-31.4165791, -64.1855098),
			Shop:"Libertad",
			Items:itemsPurchaseA,
		},
		{
			Time: time.Now().AddDate(0,0,1),
			Items:itemsPurchaseB,
		},
		{
			Time: time.Now().AddDate(0,1,1),
			Items:itemsPurchaseB,
		},
		{
			Time: time.Now().AddDate(0,2,1),
			Items:itemsPurchaseA,
		},

	}

	postPurchases = dto.PurchaseContainer{Purchases:setOfPurchases}
)

func Test_GET_Purchases_WITH_NO_TOKEN_Returns_An_Error(t *testing.T) {
	server := getServer(NewPurchaseService(infrastructure.NewMemDb()))
	defer server.Close()

	res, err := http.Get(getURL(server.URL))

	if !isHTTPStatus(http.StatusForbidden, res, err){
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusForbidden)
		t.FailNow()
	}
}


func Test_GET_Purchases_Returns_A_List_Of_Purchases_By_User(t *testing.T) {


	server := getServer(NewPurchaseService(infrastructure.NewMemDb()))
	defer server.Close()

	res, err := httpPost(user1, getURL(server.URL), postPurchases)

	if !isHTTPStatus(http.StatusCreated, res, err){
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusCreated)
		t.FailNow()
	}

	res, err = httpGet(user1, getURL(server.URL))

	if !isHTTPStatus(http.StatusOK, res, err){
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusOK)
		t.FailNow()
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal("Error")
		t.FailNow()
	}

	purchases := new(dto.PurchaseContainer)

	if err := json.Unmarshal(body, purchases); err != nil {

		log.Printf("Error when reading response %s", err)
		t.FailNow()
	}

	if len(purchases.Purchases) != len(setOfPurchases){
		log.Printf("Error: Expected items quantity is different from the received one")
		t.FailNow()

	}

	for _, purchase := range postPurchases.Purchases {
		p := purchases.GetPurchase(purchase.Time)
		if p == nil {
			log.Print("Error, purchases saved not found")
			t.FailNow();
		}
	}
}


func Test_GET_Purchases_Returns_A_Purchase_With_Latitude_and_Long(t *testing.T) {

	server := getServer(NewPurchaseService(infrastructure.NewMemDb()))
	defer server.Close()


	res, err := httpPost(user1, getURL(server.URL), postPurchases)

	if !isHTTPStatus(http.StatusCreated, res, err){
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusCreated)
		t.FailNow()
	}


	res, err = httpGet(user1, getURL(server.URL))

	if !isHTTPStatus(http.StatusOK, res, err){
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusOK)
		t.FailNow()
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal("Error")
		t.FailNow()
	}

	purchases := new(dto.PurchaseContainer)

	if err := json.Unmarshal(body, purchases); err != nil {

		log.Printf("Error when reading response %s", err)
		t.FailNow()
	}

	if len(purchases.Purchases) != len(setOfPurchases){
		log.Printf("Error: Expected items quantity is different from the received one")
		t.FailNow()

	}
	purchaseWithLatAndLong := postPurchases.Purchases[0]
	p := purchases.GetPurchase(purchaseWithLatAndLong.Time)

	if p == nil ||
	p.Point.Lat != purchaseWithLatAndLong.Point.Lat ||
	p.Point.Long != purchaseWithLatAndLong.Point.Long ||
	p.Shop != purchaseWithLatAndLong.Shop {

		log.Print("Error, purchases saved not found")
		t.FailNow();
	}

}

func Test_GET_Purchases_Grouped_By_Month_Returns_A_List_Of_Purchases_Groups(t *testing.T) {


	server := getServer(NewPurchaseService(infrastructure.NewMemDb()))
	defer server.Close()

	res, err := httpPost(user1, getURL(server.URL), postPurchases)

	if !isHTTPStatus(http.StatusCreated, res, err){
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusCreated)
		t.FailNow()
	}

	res, err = httpGet(user1, getURL(server.URL) + "?groupBy=month")

	if !isHTTPStatus(http.StatusOK, res, err){
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusOK)
		t.FailNow()
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal("Error")
		t.FailNow()
	}

	purchases := new(dto.PurchasesByMonthContainer)

	if err := json.Unmarshal(body, purchases); err != nil {

		log.Printf("Error when reading response %s", err)
		t.FailNow()
	}

	if len(purchases.PurchasesByMonth) != 3{
		log.Printf("Error: Expected items quantity is different from the received one")
		t.FailNow()

	}

	log.Printf("GET items returned OK %s", body)
}

func Test_GET_Purchases_Grouped_By_ANYTHING_Returns_A_List_Of_Purchases_Grouped_By_Month(t *testing.T) {

	server := getServer(NewPurchaseService(infrastructure.NewMemDb()))
	defer server.Close()

	res, err := httpPost(user1, getURL(server.URL), postPurchases)

	if !isHTTPStatus(http.StatusCreated, res, err){
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusCreated)
		t.FailNow()
	}


	res, err = httpGet(user1, getURL(server.URL) + "?groupBy=ANYTHING")

	if !isHTTPStatus(http.StatusOK, res, err){
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusOK)
		t.FailNow()
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal("Error")
		t.FailNow()
	}

	purchases := new(dto.PurchasesByMonthContainer)

	if err := json.Unmarshal(body, purchases); err != nil {

		log.Printf("Error when reading response %s", err)
		t.FailNow()
	}

	if len(purchases.PurchasesByMonth) != 3{
		log.Printf("Error: Expected items quantity is different from the received one")
		t.FailNow()

	}

	log.Printf("GET items returned OK %s", body)
}

func Test_GET_Purchases_From_Other_User_Responds_different_purchases(t *testing.T) {

	server := getServer(NewPurchaseService(infrastructure.NewMemDb()))
	defer server.Close()

	res, err := httpPost(user1, getURL(server.URL), postPurchases)

	if !isHTTPStatus(http.StatusCreated, res, err){
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusCreated)
		t.FailNow()
	}

	res, err = httpGet(user2, getURL(server.URL) + "?groupBy=month")

	if !isHTTPStatus(http.StatusOK, res, err){
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusOK)
		t.FailNow()
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal("Error")
		t.FailNow()
	}

	purchases := new(dto.PurchasesByMonthContainer)

	if err := json.Unmarshal(body, purchases); err != nil {

		log.Printf("Error when reading response %s", err)
		t.FailNow()
	}

	if len(purchases.PurchasesByMonth) == 3{
		log.Printf("Error: Expected items quantity is different from the received one")
		t.FailNow()

	}

	log.Printf("GET items returned OK %s", body)
}

/*func Test_aws(t *testing.T){

	endpoint := "http://172.17.0.2:8000"
	svc := dynamodb.New(session.New(&aws.Config{Region: aws.String("us-west-2"), Endpoint:&endpoint}))

	it := map[string]* dynamodb.AttributeValue {
		"user": {
			S: aws.String("lagigliaivan"),
		},
		"purchase": {
			S: aws.String("2016-04-12T00:06:22.364Z"),
		},
		"location": {
			S: aws.String("-31.4165791, -64.1855098"),
		},
		"shop":{
			S: aws.String("Carrefour"),
		},
	}
	tname := "Purchases"
	putItem := dynamodb.PutItemInput{Item:it, TableName:&tname}


	result, err := svc.PutItem(&putItem)
	//result, err := svc.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Result :%s ", result)
	//tname := "Purchases"
	key := map[string]* dynamodb.AttributeValue {
		 "user": {
			 S: aws.String("lagigliaivan"),
		 },
		 "purchase": {
			 S: aws.String("2016-04-12T00:06:22.364Z"),
		 },
	 }

	item := dynamodb.GetItemInput{Key:key, TableName:&tname}
	itemResult, err := svc.GetItem(&item)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Result:%s ", itemResult)
	*//*for _, table := range result.TableNames {
		log.Println(*table)
	}*//*
}*/

/*func isHTTPStatus(httpStatus int, res *http.Response, err error ) bool {
	return !( (err != nil) || (res.StatusCode != httpStatus) )
}*/

func getURL(url string) string{
	return url + "/catalog/purchases"
}