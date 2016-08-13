package main

import (
	"net/http"
	"log"
	"strings"
	"crypto/sha1"
	"io"
	"fmt"
	"time"
	"testing"
	"io/ioutil"
	"encoding/json"
)


var user1 string
var user2 string


const (
	DYNAMODB = 1
	MEMORYDB = 2

	TESTDB = MEMORYDB    //Change here to test services by using either mem or dynamo db
)

func init() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	sha := sha1.New()
	io.WriteString(sha, "mayname:password@gmail.com.ar")
	user1 = fmt.Sprintf("%x", sha.Sum(nil))

	sha.Reset()
	io.WriteString(sha, "mayname2:password@gmail.com.ar")
	user2 = fmt.Sprintf("%x", sha.Sum(nil))

}

func NewDB(dbType int) DB{

	switch dbType {

		case DYNAMODB:
			db, _ := NewDynamoDB("http://localhost:8000", "us-west-2")
			return db

		case MEMORYDB:
			return NewMemDb()
	}

	return NewMemDb()
}

func getDB(dbType int) DB{
	return NewDB(dbType)
}

const STATUS_ERROR_MESSAGE string = "%s %s Received status code: %d different from what was expected: %d"

const tt = "2016-01-12T00:01:23Z"

var (
	itemsPurchaseA = []Item{

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
			Description: "fourth product",
			Price: 22.0,
		},
	}

	itemsPurchaseB = []Item{

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
			Description: "fourth product",
			Price: 212.0,
		},
	}

	timeToTest,_ = time.Parse(time.RFC3339, tt)

	setOfPurchases = []Purchase{

		{
			Time: timeToTest,
			Location:NewPoint(-31.4165791, -64.1855098),
			Shop:"Libertad",
			Items:itemsPurchaseA,
		},
		{
			Time: timeToTest.AddDate(0,0,1),
			Items:itemsPurchaseB,
			Shop:"Libertad",
		},
		{
			Time: timeToTest.AddDate(0,1,1),
			Items:itemsPurchaseB,
			Shop:"Libertad",
		},
		{
			Time: timeToTest.AddDate(0,2,1),
			Items:itemsPurchaseA,
			Shop:"Libertad",
		},
		{
			Time: timeToTest.AddDate(0,3,1),
			Items:itemsPurchaseA,
			Shop:"Libertad",
		},

	}

	postPurchases = PurchaseContainer{Purchases:setOfPurchases}
)

func Test_GET_Purchases_WITH_NO_TOKEN_Returns_An_Error(t *testing.T) {
	server := getServer(NewPurchaseService(getDB(TESTDB)))
	defer server.Close()

	res, err := http.Get(getURL(server.URL))

	if !isHTTPStatus(http.StatusForbidden, res, err){
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusForbidden)
		t.FailNow()
	}
}


func Test_GET_Purchases_Returns_A_List_Of_Purchases_By_User(t *testing.T) {


	server := getServer(NewPurchaseService(getDB(TESTDB)))
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

	purchases := new(PurchaseContainer)

	if err := json.Unmarshal(body, purchases); err != nil {

		log.Printf("Error when reading response %s", err)
		t.FailNow()
	}

	if len(purchases.Purchases) != len(setOfPurchases){
		log.Printf("Error: Expected items quantity is different from the received one")
		t.FailNow()

	}

	for _, purchase := range postPurchases.Purchases {
		p := purchases.GetPurchase(purchase.Time.Unix())
		if p == nil {
			log.Print("Error, purchases saved not found")
			t.FailNow();
		}
	}
}


func Test_GET_Purchases_Returns_A_Purchase_With_Latitude_and_Long(t *testing.T) {

	server := getServer(NewPurchaseService(getDB(TESTDB)))
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

	purchases := new(PurchaseContainer)

	if err := json.Unmarshal(body, purchases); err != nil {

		log.Printf("Error when reading response %s", err)
		t.FailNow()
	}

	if len(purchases.Purchases) != len(setOfPurchases){
		log.Printf("Error: Expected items quantity is different from the received one")
		t.FailNow()

	}
	purchaseWithLatAndLong := postPurchases.Purchases[0]
	p := purchases.GetPurchase(purchaseWithLatAndLong.Time.UTC().Unix())

	if p == nil ||
	p.Location.Lat != purchaseWithLatAndLong.Location.Lat ||
	p.Location.Long != purchaseWithLatAndLong.Location.Long ||
	p.Shop != purchaseWithLatAndLong.Shop {
		log.Print("Error, purchases saved not found")
		t.FailNow();
	}
}

func Test_GET_Purchases_Grouped_By_Month_Returns_A_List_Of_Purchases_Groups(t *testing.T) {


	server := getServer(NewPurchaseService(getDB(TESTDB)))
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

	purchases := new(PurchasesByMonthContainer)

	if err := json.Unmarshal(body, purchases); err != nil {

		log.Printf("Error when reading response %s", err)
		t.FailNow()
	}

	if len(purchases.PurchasesByMonth) != 4{
		log.Printf("Error: Expected items quantity is different from the received one: %d", len(purchases.PurchasesByMonth))
		t.FailNow()

	}



}

func Test_GET_Purchases_Grouped_By_ANYTHING_Returns_A_List_Of_Purchases_Grouped_By_Month(t *testing.T) {

	server := getServer(NewPurchaseService(getDB(TESTDB)))
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

	purchases := new(PurchasesByMonthContainer)

	if err := json.Unmarshal(body, purchases); err != nil {

		log.Printf("Error when reading response %s", err)
		t.FailNow()
	}

	if len(purchases.PurchasesByMonth) != 4{
		log.Printf("Error: Expected items quantity is different from the received one")
		t.FailNow()

	}
}

func Test_GET_Purchases_From_Other_User_Responds_different_purchases(t *testing.T) {

	server := getServer(NewPurchaseService(getDB(TESTDB)))
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

	purchases := new(PurchasesByMonthContainer)

	if err := json.Unmarshal(body, purchases); err != nil {

		log.Printf("Error when reading response %s", err)
		t.FailNow()
	}

	if len(purchases.PurchasesByMonth) == 3{
		log.Printf("Error: Expected items quantity is different from the received one")
		t.FailNow()

	}
}

func Test_DELETE_A_Purchase(t *testing.T) {

	server := getServer(NewPurchaseService(getDB(TESTDB)))
	defer server.Close()

	res, err := httpPost(user1, getURL(server.URL), postPurchases)

	if !isHTTPStatus(http.StatusCreated, res, err){
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusCreated)
		t.FailNow()
	}

	purchaseToDelete := getURL(server.URL) + "/" + fmt.Sprintf("%d", timeToTest.Unix());

	res, err = httpDelete(user1, purchaseToDelete)

	if !isHTTPStatus(http.StatusOK, res, err){
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusOK)
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

	purchases := new(PurchasesByMonthContainer)

	if err := json.Unmarshal(body, purchases); err != nil {

		log.Printf("Error when reading response %s", err)
		t.FailNow()
	}

	var purchasesByMonth []PurchasesByMonth  = purchases.PurchasesByMonth;


	for _, purchases := range purchasesByMonth {

		for _, p := range purchases.Purchases {
			if strings.Compare(p.Id, fmt.Sprintf("%d", timeToTest.Unix()) ) == 0 {
				log.Printf("%s %s", p.Time, timeToTest)
				t.FailNow()
			}
		}
	}

	if len(purchases.PurchasesByMonth) != 4 {
		log.Printf("Error: Expected items quantity is different from the received one: %d", len(purchases.PurchasesByMonth))
		t.FailNow()

	}
}

func getURL(url string) string{
	return url + "/catalog/purchases"
}

