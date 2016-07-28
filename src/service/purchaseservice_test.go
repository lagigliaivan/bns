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

	TESTDB = DYNAMODB    //Change here to test services by using either mem or dynamo db
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
	purchaseTime = time.Now()

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
			Description: "forth product",
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
			Description: "forth product",
			Price: 212.0,
		},
	}

	timeToTest,_ = time.Parse(time.RFC3339, tt)

	setOfPurchases = []Purchase{

		{
			Time: timeToTest,
			Point:NewPoint(-31.4165791, -64.1855098),
			Shop:"Libertad",
			Items:itemsPurchaseA,
		},
		{
			Time: timeToTest.AddDate(0,0,1),
			Items:itemsPurchaseB,
			Shop:"Libertad",
		},
		{
			Time: purchaseTime.AddDate(0,1,1),
			Items:itemsPurchaseB,
			Shop:"Libertad",
		},
		{
			Time: purchaseTime.AddDate(0,2,1),
			Items:itemsPurchaseA,
			Shop:"Libertad",
		},
		{
			Time: purchaseTime.AddDate(0,3,1),
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
		p := purchases.GetPurchase(purchase.Time)
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

	log.Printf("GET items returned OK %s", body)
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

	fmt.Printf("purchases %s", purchases.PurchasesByMonth)
	if len(purchases.PurchasesByMonth) != 4 {
		log.Printf("Error: Expected items quantity is different from the received one: %d", len(purchases.PurchasesByMonth))
		t.FailNow()

	}
}

func getURL(url string) string{
	return url + "/catalog/purchases"
}

/*func getDynamoDBItem(id string, dt string, user string, shop string, items ItemContainer) map[string]* dynamodb.AttributeValue {

	i, err := strconv.ParseInt(dt, 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(i, 0)

	it := map[string]* dynamodb.AttributeValue {
		"id": {
			S: aws.String(id),
		},
		"dt": {
			N: aws.String(dt),
		},
		"date": {
			S: aws.String(tm.UTC().Format(time.RFC3339)),
		},
		"user":{
			S: aws.String(user),
		},
		"shop":{
			S: aws.String(shop),
		},
		"items":{
			S: aws.String(items.ToJsonString()),
		},
	}

	return it
}

var endpoint = "http://localhost:8000"
var tname = "Purchases"

/*

var dts [10]int64
var now int64 = 1469652314 //time.Unix()

func init(){

	for i :=int64(0); i<10; i++ {
		dts[i] = (now + i)
	}
}

func Test_aws_purchases_creation(t *testing.T) {

	count := 0

	svc := dynamodb.New(session.New(&aws.Config{Region: aws.String("us-west-2"), Endpoint:&endpoint}))

	for _, dt := range dts {


		items := []Item { {Id:"12312313", Description:"Cafe la morenita", Price:10},  {Id:"3332", Description:"Jabon de tocador", Price:13.4}}
		itemsContainer := new (ItemContainer)
		itemsContainer.Items = items

		user := "mayuser:password@gmail.com.ar"
		buildDynamoItem(user1, user)

		it := getDynamoDBItem(user1, fmt.Sprintf("%d", dt), user, "carrefour", *itemsContainer)
		putItem := dynamodb.PutItemInput{Item:it, TableName:&tname}

		result, err := svc.PutItem(&putItem)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("Result :%s ", result)

		user = "mayuser2:password@gmail.com.ar"
		it = getDynamoDBItem(user2, fmt.Sprintf("%d", dt), user, "carrefour", *itemsContainer)
		putItem = dynamodb.PutItemInput{Item:it, TableName:&tname}

		count = count+2

		//log.Println("Result :%s ", result)
	}


	log.Printf("%d items were inserted", count)


	key := map[string]* dynamodb.AttributeValue {

		"id": {
			S: aws.String(user1),
		},
		"dt": {
			N: aws.String(fmt.Sprintf("%d",dts[1])),
		},
	}

	item := dynamodb.GetItemInput{Key:key, TableName:&tname}
	itemResult, err := svc.GetItem(&item)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Result:%s ", itemResult)

}

func Test_aws_get_items(t *testing.T) {


	svc := dynamodb.New(session.New(&aws.Config{Region: aws.String("us-west-2"), Endpoint:&endpoint}))

	for _, d := range dts {
		key := map[string]*dynamodb.AttributeValue{

			"id": {
				S: aws.String(user1),
			},
			"dt": {
				N: aws.String(fmt.Sprintf("%d", d)),
			},
		}

		item := dynamodb.GetItemInput{Key:key, TableName:&tname}
		itemResult, err := svc.GetItem(&item)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("Result:%s ", itemResult)
	}
}


func Test_aws_items_query(t *testing.T){



	svc := dynamodb.New(session.New(&aws.Config{Region: aws.String("us-west-2"), Endpoint:&endpoint}))


	id := user1

	params := &dynamodb.QueryInput{
		TableName: aws.String(tname),
		ConsistentRead:      aws.Bool(true),
		ExpressionAttributeNames: map[string]*string{
			"#s": aws.String("shop"),
			"#i": aws.String("items"),
			"#d": aws.String("dt"),
			"#t": aws.String("date"),
		},
		ProjectionExpression: aws.String("#s, #i, #d, #t"),
		ExpressionAttributeValues: map[string] *dynamodb.AttributeValue {
			":v1": {
				S:    aws.String(id),
			},
			":v2": {
				N:    aws.String(fmt.Sprintf("%d",now)),
			},
			":v3": {
				N:    aws.String(fmt.Sprintf("%d",now + 5)),
			},
		},
		KeyConditionExpression: aws.String("id = :v1 AND dt BETWEEN :v2 AND :v3 "),

	}
	resp, err := svc.Query(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	parseQueryResponse(resp.Items)

}


func parseQueryResponse (items []map[string]*dynamodb.AttributeValue) {

		for _, p := range items{
			//fmt.Println(p)
			t, err := time.Parse(time.RFC3339, *(p["date"].S))

			if err != nil {
				fmt.Println("Error %s", err)
			}

			id := *(p["dt"].N)

			itemsContainer := new(ItemContainer)
			if err := json.Unmarshal([]byte(*(p["items"].S)), itemsContainer); err != nil {

				log.Printf("Error when reading response %s", err)
				//t.FailNow()
			}

			purchase := Purchase{Id:id, Time:t, Shop:*(p["shop"].S), Items:itemsContainer.Items}
			fmt.Println(purchase)
		}

}




func Test_DeletePurchase(t *testing.T)  {

	svc := dynamodb.New(session.New(&aws.Config{Region: aws.String("us-west-2"), Endpoint:&endpoint}))
	params := &dynamodb.DeleteItemInput{

		Key: map[string]*dynamodb.AttributeValue{ // Required
			"id": {
				S:    aws.String(user1),
			},
			"dt": {
				N:    aws.String(fmt.Sprintf("%d", now)),
			},
		},
		TableName:           aws.String(TABLE_PURCHASES), // Required
	}

	resp, err := svc.DeleteItem(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
}

*/



