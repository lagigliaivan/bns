package main

import (
	"time"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"fmt"
	"encoding/json"
	"log"
	"strings"
	"errors"
)
const (
	PURCHASES = "Purchases"
)

type DB interface {
	GetItem(string) Item
	SaveItem(Item) int
	GetItems() []Item

	SavePurchase(Purchase, string) error
	GetPurchases(string) []Purchase
	GetPurchasesByMonth(string, int) map[time.Month][]Purchase

}

type DynamoDB struct {
	endpoint string
	svc *dynamodb.DynamoDB
}

func NewDynamoDB(endpoint, region string) (*DynamoDB, error) {


	var config *aws.Config

	if strings.Compare(region, "") == 0 {
		return nil, errors.New("region cannot be nil")
	}

	if strings.Compare(endpoint, "") == 0 {
		config = &aws.Config{Region: aws.String(region)}
	}else{
		config = &aws.Config{Region: aws.String(region), Endpoint:&endpoint}
	}

	catalogDB := new(DynamoDB)
	catalogDB.endpoint = endpoint
	catalogDB.svc = dynamodb.New(session.New(config))

	return catalogDB, nil
}

func (catDb DynamoDB) GetItem(id string) (Item){
	item := Item{}
	item.Id = id
	return item
}

func (catDb DynamoDB) GetItems() []Item{
	return nil
}

func (catDb DynamoDB) SaveItem(Item) int {

	return 0
}

func (catDb DynamoDB) GetPurchase(time time.Time) Purchase  {

	return Purchase{}
}

func (catDb DynamoDB) SavePurchase( p Purchase, userId string) error {

	tableName := PURCHASES

	it := buildDynamoItem(p, userId)

	putItem := dynamodb.PutItemInput{Item:it, TableName:&tableName}

	result, err := catDb.svc.PutItem(&putItem)


	if err != nil {
		log.Println(err)
		return err
	}

	log.Println(result)
	return nil
}

func (catDb DynamoDB) GetPurchases(user string) []Purchase  {

	resp, err := catDb.getPurchasesFromAWS(user, time.Now().Year())

	if err != nil {
		log.Printf("Error while querying DB %s\n", err)
		return []Purchase{}
	}

	purchases := []Purchase{}

	for _, p := range resp.Items{

		t, err := time.Parse(time.RFC3339, *(p["dt"].S))

		if err != nil {
			fmt.Printf("Error while parsing Purchase date: %s \n", err)
			return []Purchase{}
		}

		itemsContainer := new(ItemContainer)
		if err := json.Unmarshal([]byte(*(p["items"].S)), itemsContainer); err != nil {

			log.Printf("Error when reading response %s", err)
			return []Purchase{}
		}

		purchase := Purchase{Time:t, Shop:*(p["shop"].S), Items:itemsContainer.Items}

		purchases = append(purchases, purchase)
		fmt.Println(purchase)
	}

	return purchases
}

func (catDb DynamoDB) GetPurchasesByMonth(user string, year int) map[time.Month][]Purchase  {


	resp, err := catDb.getPurchasesFromAWS(user, year)

	if err != nil {
		log.Printf("Error while querying DB %s\n", err)
		return make(map[time.Month][]Purchase)
	}

	purchasesByMonth := make(map[time.Month][]Purchase)

	for _, p := range resp.Items{

		t, err := time.Parse(time.RFC3339, *(p["dt"].S))

		if err != nil {
			fmt.Printf("Error while parsing Purchase date: %s \n", err)
			return make(map[time.Month][]Purchase)
		}

		itemsContainer := new(ItemContainer)
		if err := json.Unmarshal([]byte(*(p["items"].S)), itemsContainer); err != nil {

			log.Printf("Error when reading response %s", err)
			return make(map[time.Month][]Purchase)
		}

		purchase := Purchase{Time:t, Shop:*(p["shop"].S), Items:itemsContainer.Items}

		if purchasesByMonth[t.Month()] == nil {
			purchasesByMonth[t.Month()] = make([]Purchase,0)
		}
		purchasesByMonth[t.Month()] = append(purchasesByMonth[t.Month()], purchase)

	}

	return purchasesByMonth
}

func (catDb DynamoDB) getPurchasesFromAWS(user string, year int) ( *dynamodb.QueryOutput, error) {

	log.Println("Querying AWS Dynamodb")

	from := fmt.Sprintf("%d%s", year, "-01-00T00:00:00Z")
	to := fmt.Sprintf("%d%s", year, "-12-31T23:59:00Z")

	params := &dynamodb.QueryInput{
		TableName: aws.String(PURCHASES),
		ConsistentRead: aws.Bool(true),
		ExpressionAttributeValues: map[string] *dynamodb.AttributeValue {
			":v1": {
				S:    aws.String(user),
			},
			":v2": {
				S:    aws.String(from),
			},
			":v3": {
				S:    aws.String(to),
			},
		},
		KeyConditionExpression: aws.String("id = :v1 AND dt BETWEEN :v2 AND :v3 "),
	}

	resp, err := catDb.svc.Query(params)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return resp, nil
}


//func buildDynamoItem(id, dt, user, shop, items string) map[string]* dynamodb.AttributeValue {
func buildDynamoItem(purchase Purchase, user string) map[string]* dynamodb.AttributeValue {


	shop := purchase.Shop

	if shop == "" {
		shop = "-"
	}

	itemsContainer := ItemContainer{}

	for _, item := range purchase.Items {
		itemsContainer.Add(item)
	}

	it := map[string]* dynamodb.AttributeValue {
		"id": {
			S: aws.String(user),
		},
		"dt": {
			S: aws.String(purchase.Time.Format(time.RFC3339)),
		},
		"user":{
			S: aws.String(user),
		},
		"shop":{
			S: aws.String(shop),
		},
		"items":{
			S: aws.String(itemsContainer.ToJsonString()),
		},
	}

	return it
}