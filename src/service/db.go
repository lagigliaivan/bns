package main

import (
	//"fmt"
	"time"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"fmt"
)

type DB interface {
	GetItem(string) Item
	SaveItem(Item) int
	GetItems() []Item

	//GetPurchase(time.Time) Purchase
	SavePurchase(Purchase, string) error
	GetPurchases(string) []Purchase
	GetPurchasesByMonth(string, int) map[time.Month][]Purchase
	//GetPurchasesByUser(user string)
}

type CatalogDB struct {
	endpoint string
	svc *dynamodb.DynamoDB
}

func NewCatalogDB() *CatalogDB {

	catalogDB := new(CatalogDB)
	catalogDB.endpoint = "http://localhost:8000"
	catalogDB.svc = dynamodb.New(session.New(&aws.Config{Region: aws.String("us-west-2"), Endpoint:&catalogDB.endpoint}))

	return catalogDB
}

func (catDb CatalogDB) GetItem(id string) (Item){
	item := Item{}
	item.Id = id
	return item
}

func (catDb CatalogDB) GetItems() []Item{
	return nil
}

func (catDb CatalogDB) SaveItem(Item) int {

	return 0
}

func (catDb CatalogDB) GetPurchase(time time.Time) Purchase  {

	return Purchase{}
}

func (catDb CatalogDB) SavePurchase( p Purchase, userId string) error {


	return nil
}

func (catDb CatalogDB) GetPurchases() []Purchase  {


	return []Purchase{}
}

func (catDb CatalogDB) GetPurchasesByMonth(user string, year int) map[time.Month][]Purchase  {

	table := "purchases"

	from := fmt.Sprintf("%d%s", year, "-01-00T00:00:00Z")
	to := fmt.Sprintf("%d%s", year, "-12-31T23:59:00Z")

 	params := &dynamodb.QueryInput{
		TableName: aws.String(table),
		ConsistentRead:      aws.Bool(true),
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
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return make(map[time.Month][]Purchase)
	}

	// Pretty-print the response data.
	fmt.Println(resp)

	/*body, err := ioutil.ReadAll(resp)

	if err != nil {
		log.Fatal("Error")

	}


	purchases := new(PurchasesByMonthContainer)

	if err := json.Unmarshal(body, purchases); err != nil {

		log.Printf("Error when reading response %s", err)
		//t.FailNow()
	}
*/

	return make(map[time.Month][]Purchase)
}

func (catDb CatalogDB) GetPurchasesByUser(user string) []Purchase  {
	return []Purchase{}
}