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
	TABLE_PURCHASES = "Purchases"
	TABLE_ITEMS_DESCRIPTIONS = "ItemsDescriptions"
)

type DB interface {
	GetItem(string) Item
	SaveItem(Item) int
	SaveItemsDescriptions(string, []ItemDescription) error
	GetItemsDescriptions(string) ([]ItemDescription, error)

	GetItems() []Item

	SavePurchase(Purchase, string) error
	GetPurchases(string) []Purchase
	GetPurchasesByMonth(string, int) map[time.Month] []Purchase

	DeletePurchase(string, string)

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

	dynamoDb := new(DynamoDB)
	dynamoDb.endpoint = endpoint
	dynamoDb.svc = dynamodb.New(session.New(config))

	return dynamoDb, nil
}

func (db DynamoDB) GetItem(id string) (Item){
	item := Item{}
	item.Id = id
	return item
}

func (db DynamoDB) GetItems() []Item{
	return nil
}

func (db DynamoDB) SaveItem(Item) int {

	return 0
}

func (db DynamoDB) GetPurchase(time time.Time) Purchase  {

	return Purchase{}
}

func (db DynamoDB) SavePurchase( p Purchase, userId string) error {

	tableName := TABLE_PURCHASES

	it := buildDynamoPurchaseItem(p, userId)

	putItem := dynamodb.PutItemInput{Item:it, TableName:&tableName}

	_, err := db.svc.PutItem(&putItem)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (db DynamoDB) GetPurchases(user string) []Purchase  {

	resp, err := db.getPurchasesFromAWS(user, time.Now().Year())

	if err != nil {
		log.Printf("Error while querying DB %s\n", err)
		return []Purchase{}
	}

	purchases := getPurchases(resp)

	return purchases
}

func (db DynamoDB) GetPurchasesByMonth(user string, year int) map[time.Month][]Purchase  {


	resp, err := db.getPurchasesFromAWS(user, year)

	if err != nil {
		log.Printf("Error while querying DB %s\n", err)
		return make(map[time.Month][]Purchase)
	}

	purchases := getPurchases(resp)

	purchasesByMonth := make(map[time.Month][]Purchase)


	for _, purchase := range purchases {

		if purchasesByMonth[purchase.Time.Month()] == nil {
			purchasesByMonth[purchase.Time.Month()] = make([]Purchase,0)
		}
		purchasesByMonth[purchase.Time.Month()] = append(purchasesByMonth[purchase.Time.Month()], purchase)
	}

	return purchasesByMonth
}

func (db DynamoDB) DeletePurchase(user string, id string)  {

	params := &dynamodb.DeleteItemInput{

		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S:    aws.String(user),
			},
			"dt": {
				N:    aws.String(id),
			},
		},
		TableName:aws.String(TABLE_PURCHASES),
	}

	_, err := db.svc.DeleteItem(params)

	if err != nil {
		// Print the error, cast err to aws err.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

}


func (db DynamoDB) SaveItemsDescriptions(userId string, itemsDescriptions []ItemDescription )  error {

	tableName := TABLE_ITEMS_DESCRIPTIONS

	for _, itemDescription := range itemsDescriptions {

		it := buildDynamoItemDescriptionItem(itemDescription.ItemId, itemDescription.Description, userId)

		putItem := dynamodb.PutItemInput{Item:it, TableName:&tableName}

		_, err := db.svc.PutItem(&putItem)

		if err != nil {
			log.Println(err)
			return err
		}
		log.Printf("Saving item for user:%s itemid:%s description:%s", userId, itemDescription.ItemId, itemDescription.Description)
	}


	return nil
}

func (db DynamoDB) getPurchasesFromAWS(user string, year int) ( *dynamodb.QueryOutput, error) {

	log.Println("Querying AWS Dynamodb")

	from := fmt.Sprintf("%d%s", year, "-01-00T00:00:00Z")
	to := fmt.Sprintf("%d%s", year, "-12-31T23:59:00Z")

	fromInMillis, err := time.Parse(time.RFC3339, from)

	if err != nil {
		log.Printf("Error while parsing year from -- this error should not happen: %s", err.Error())
		return nil, err
	}

	toInMillis, err := time.Parse(time.RFC3339, to)


	if err != nil {
		log.Printf("Error while parsing year to -- this error should not happen: %s", err.Error())
		return nil, err
	}

	params := &dynamodb.QueryInput{
		TableName: aws.String(TABLE_PURCHASES),
		ConsistentRead: aws.Bool(true),
		ExpressionAttributeValues: map[string] *dynamodb.AttributeValue {
			":v1": {
				S:    aws.String(user),
			},
			":v2": {
				N:    aws.String(fmt.Sprintf("%d", fromInMillis.Unix())),
			},
			":v3": {
				N:    aws.String(fmt.Sprintf("%d", toInMillis.Unix())),
			},
		},
		KeyConditionExpression: aws.String("id = :v1 AND dt BETWEEN :v2 AND :v3 "),
	}

	resp, err := db.svc.Query(params)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return resp, nil
}

func (db DynamoDB) GetItemsDescriptions (user string) ( []ItemDescription, error) {


	params := &dynamodb.QueryInput{
		TableName: aws.String(TABLE_ITEMS_DESCRIPTIONS),
		ConsistentRead: aws.Bool(true),
		ExpressionAttributeValues: map[string] *dynamodb.AttributeValue {
			":v1": {
				S:    aws.String(user),
			},
		},
		KeyConditionExpression: aws.String("userid = :v1"),
	}

	log.Printf("Quering items for user: %s", user)

	resp, err := db.svc.Query(params)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	items := getItemsDescription(resp)
	return items, nil
}

func buildDynamoPurchaseItem(purchase Purchase, user string) map[string]* dynamodb.AttributeValue {

	shop := purchase.Shop

	itemsContainer := ItemContainer{}

	for _, item := range purchase.Items {
		itemsContainer.Add(item)
	}

	it := map[string]* dynamodb.AttributeValue {

		"id": {
			S: aws.String(user),
		},
		"dt": {
			N: aws.String(purchase.Id),
		},
		"date": {
			S: aws.String(purchase.Time.Format(time.RFC3339)),
		},
		"user":{
			S: aws.String(user),
		},
		"shop":{
			S: aws.String(shop),
		},
		"location":{
			S: aws.String(purchase.Location.toString()),
		},
		"items":{
			S: aws.String(itemsContainer.ToJsonString()),
		},
	}

	return it
}

func buildDynamoItemDescriptionItem(itemId string, description string, user string) map[string]* dynamodb.AttributeValue {

	it := map[string]* dynamodb.AttributeValue {

		"userid": {
			S: aws.String(user),
		},
		"itemid": {
			S: aws.String(itemId),
		},
		"description": {
			S: aws.String(description),
		},
	}

	return it
}


func getPurchases(awsResponse *dynamodb.QueryOutput) []Purchase {

	purchases := []Purchase{}

	for _, p := range awsResponse.Items{

		t, err := time.Parse(time.RFC3339, *(p["date"].S))

		if err != nil {
			fmt.Println("Error %s", err)
		}

		itemsContainer := new(ItemContainer)
		if err := json.Unmarshal([]byte(*(p["items"].S)), itemsContainer); err != nil {

			log.Printf("Error when reading response %s", err)
			return []Purchase{}
		}

		purchase := Purchase{Id: *(p["dt"].N), Time:t, Shop:*(p["shop"].S), Items:itemsContainer.Items}

		purchases = append(purchases, purchase)
	}

	return purchases
}

func getItemsDescription(awsResponse *dynamodb.QueryOutput) []ItemDescription {

	var itemsDescriptions []ItemDescription

	for _, item := range awsResponse.Items{

		itemDesc := ItemDescription{ItemId:*(item["itemid"].S), Description:*(item["description"].S)}
		itemsDescriptions = append(itemsDescriptions, itemDesc)

	}

	return itemsDescriptions
}