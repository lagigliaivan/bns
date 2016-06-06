package services

import (
	"testing"
	"net/http"
	"log"
	"time"
	"github.com/pos/infrastructure"
	"encoding/json"
	"github.com/pos/dto"
	"io/ioutil"

/*	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"*/
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
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
/*
func Test_GET_Purchases_WITH_NO_TOKEN_Returns_An_Error(t *testing.T) {
	server := getServer(NewPurchaseService(infrastructure.NewMemDb()))
	defer server.Close()

	res, err := http.Get(getURL(server.URL))

	if !isHTTPStatus(http.StatusForbidden, res, err){
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusForbidden)
		t.FailNow()
	}
}
*/

func Test_GET_Purchases_Returns_A_List_Of_Purchases(t *testing.T) {


	server := getServer(NewPurchaseService(infrastructure.NewMemDb()))
	defer server.Close()

	res, err := httpPost(getURL(server.URL), postPurchases)

	if !isHTTPStatus(http.StatusCreated, res, err){
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusCreated)
		t.FailNow()
	}


	res, err = http.Get(getURL(server.URL))

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


	res, err := httpPost(getURL(server.URL), postPurchases)

	if !isHTTPStatus(http.StatusCreated, res, err){
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusCreated)
		t.FailNow()
	}


	res, err = http.Get(getURL(server.URL))

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

	res, err := httpPost(getURL(server.URL), postPurchases)

	if !isHTTPStatus(http.StatusCreated, res, err){
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusCreated)
		t.FailNow()
	}

	res, err = http.Get(getURL(server.URL) + "?groupBy=month")

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

	res, err := httpPost(getURL(server.URL), postPurchases)

	if !isHTTPStatus(http.StatusCreated, res, err){
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusCreated)
		t.FailNow()
	}


	res, err = http.Get(getURL(server.URL) + "?groupBy=ANYTHING")

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