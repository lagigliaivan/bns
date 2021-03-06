package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

)

var user1 string
var user2 string

const (
	DYNAMODB = 1
	MEMORYDB = 2

	TESTDB = MEMORYDB //Change here to test services by using either mem or dynamo db
)

func init() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	sha := sha1.New()
	io.WriteString(sha, "mayname@gmail.com.ar")
	user1 = fmt.Sprintf("%x", sha.Sum(nil))

	sha.Reset()
	io.WriteString(sha, "mayname2@gmail.com.ar")
	user2 = fmt.Sprintf("%x", sha.Sum(nil))

}

func NewDB(dbType int) DB {

	switch dbType {

	case DYNAMODB:
		db, _ := NewDynamoDB("http://localhost:8000", "us-west-2")
		return db

	case MEMORYDB:
		return NewMemDb()
	}

	return NewMemDb()
}

func getDB(dbType int) DB {
	return NewDB(dbType)
}

const STATUS_ERROR_MESSAGE string = "%s %s Received status code: %d different from what was expected: %d"

const tt = "2016-01-12T00:01:23Z"

var (
	itemsPurchaseA = []Item{

		{
			Description: "first product",
			Price:       2.0,
		},
		{
			Description: "second product",
			Price:       34.0,
		},
		{
			Description: "third product",
			Price:       332.50,
		},
		{
			Description: "fourth product",
			Price:       22.0,
		},
	}

	itemsPurchaseB = []Item{

		{
			Description: "first product",
			Price:       122.0,
		},
		{
			Description: "second product",
			Price:       314.0,
		},
		{
			Description: "third product",
			Price:       3212.22,
		},
		{
			Description: "fourth product",
			Price:       212.0,
		},
	}

	timeToTest, _ = time.Parse(time.RFC3339, tt)

	setOfPurchases = []Purchase{

		{
			Time:     timeToTest,   //12-01
			Location: NewPoint(-31.4165791, -64.1855098),
			Shop:     "Libertad",
			Items:    itemsPurchaseA,
		},
		{
			Time:  timeToTest.AddDate(0, 0, 1), //13-01
			Items: itemsPurchaseB,
			Shop:  "Libertad",
		},
		{
			Time:  timeToTest.AddDate(0, 1, 1), //14-02
			Items: itemsPurchaseB,
			Shop:  "Libertad",
		},
		{
			Time:  timeToTest.AddDate(0, 2, 1), //15-04
			Items: itemsPurchaseA,
			Shop:  "Libertad",
		},
		{
			Time:  timeToTest.AddDate(0, 3, 1), //16-07
			Items: itemsPurchaseA,
			Shop:  "Libertad",
		},
	}

	postPurchases = PurchaseContainer{Purchases: setOfPurchases}
)

func Test_GET_Purchases_WITH_NO_TOKEN_Returns_An_Error(t *testing.T) {
	server := getServer(NewPurchaseService(getDB(TESTDB)))
	defer server.Close()

	res, err := http.Get(getURL(server.URL, user1))

	if !isHTTPStatus(http.StatusForbidden, res, err) {
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusForbidden)
		t.FailNow()
	}
}

func Test_GET_Purchases_Returns_A_List_Of_Purchases_By_User(t *testing.T) {

	server := getServer(NewPurchaseService(getDB(TESTDB)))
	defer server.Close()

	res, err := httpPost(user1, getURL(server.URL, user1), postPurchases)

	if !isHTTPStatus(http.StatusCreated, res, err) {
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusCreated)
		t.FailNow()
	}

	res, err = httpGet(user1, getURL(server.URL, user1))

	if !isHTTPStatus(http.StatusOK, res, err) {
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

	if len(purchases.Purchases) != len(setOfPurchases) {
		log.Printf("Error: Expected items quantity is different from the received one")
		t.FailNow()

	}

	for _, purchase := range postPurchases.Purchases {
		p := purchases.GetPurchase(purchase.Time.Unix())
		if p == nil {
			log.Print("Error, purchases saved not found")
			t.FailNow()
		}
	}
}

func Test_GET_Purchases_Returns_A_Purchase_With_Latitude_and_Long(t *testing.T) {

	server := getServer(NewPurchaseService(getDB(TESTDB)))
	defer server.Close()

	res, err := httpPost(user1, getURL(server.URL, user1), postPurchases)

	if !isHTTPStatus(http.StatusCreated, res, err) {
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusCreated)
		t.FailNow()
	}

	res, err = httpGet(user1, getURL(server.URL, user1))

	if !isHTTPStatus(http.StatusOK, res, err) {
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

	if len(purchases.Purchases) != len(setOfPurchases) {
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
		t.FailNow()
	}
}

func Test_GET_Purchases_Grouped_By_Month_Returns_A_List_Of_Purchases_Groups(t *testing.T) {

	server := getServer(NewPurchaseService(getDB(TESTDB)))
	defer server.Close()

	res, err := httpPost(user1, getURL(server.URL, user1), postPurchases)

	if !isHTTPStatus(http.StatusCreated, res, err) {
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusCreated)
		t.FailNow()
	}

	res, err = httpGet(user1, getURL(server.URL, user1)+"?groupBy=month")

	if !isHTTPStatus(http.StatusOK, res, err) {
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

	if len(purchases.PurchasesByMonth) != 4 {
		log.Printf("Error: Expected items quantity is different from the received one: %d", len(purchases.PurchasesByMonth))
		t.FailNow()

	}
}

func Test_GET_A_Purchase_By_Id_Returns_It_If_It_Exists(t *testing.T) {

	server := getServer(NewPurchaseService(getDB(TESTDB)))
	defer server.Close()

	res, err := httpPost(user1, getURL(server.URL, user1), postPurchases)

	if !isHTTPStatus(http.StatusCreated, res, err) {
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusCreated)
		t.FailNow()
	}

	res, err = httpGet(user1, getURL(server.URL, user1))

	if !isHTTPStatus(http.StatusOK, res, err) {
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

	for _, v := range (*purchases).Purchases {

		res, err = httpGet(user1, getURL(server.URL, user1)+"/"+v.Id)
		if err != nil {
			log.Fatal("Error")
			t.FailNow()
		}

		if !isHTTPStatus(http.StatusOK, res, err) {
			log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusOK)
			t.FailNow()
		}

		body, err := ioutil.ReadAll(res.Body)

		purchase := new(Purchase)

		if err := json.Unmarshal(body, purchase); err != nil {

			log.Printf("Error when reading response %s", err)
			t.FailNow()
		}

		if err != nil {
			log.Fatal("Error")
			t.FailNow()
		}

		if strings.Compare(purchase.Id, v.Id) != 0 {
			log.Printf("Error, returned purchase id: %s does not match the expected one %s.", purchase.Id, v.Id)
			t.FailNow()
		}

	}
}

func Test_GET_Purchases_Grouped_By_ANYTHING_Returns_A_List_Of_Purchases_Grouped_By_Month(t *testing.T) {

	server := getServer(NewPurchaseService(getDB(TESTDB)))
	defer server.Close()

	res, err := httpPost(user1, getURL(server.URL, user1), postPurchases)

	if !isHTTPStatus(http.StatusCreated, res, err) {
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusCreated)
		t.FailNow()
	}

	res, err = httpGet(user1, getURL(server.URL, user1)+"?groupBy=ANYTHING")

	if !isHTTPStatus(http.StatusOK, res, err) {
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

	if len(purchases.PurchasesByMonth) != 4 {
		log.Printf("Error: Expected items quantity is different from the received one")
		t.FailNow()

	}
}

func Test_GET_Purchases_From_Other_User_Responds_different_purchases(t *testing.T) {

	server := getServer(NewPurchaseService(getDB(TESTDB)))
	defer server.Close()

	res, err := httpPost(user1, getURL(server.URL, user1), postPurchases)

	if !isHTTPStatus(http.StatusCreated, res, err) {
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusCreated)
		t.FailNow()
	}

	res, err = httpGet(user2, getURL(server.URL, user2)+"?groupBy=month")

	if !isHTTPStatus(http.StatusOK, res, err) {
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

	if len(purchases.PurchasesByMonth) == 3 {
		log.Printf("Error: Expected items quantity is different from the received one")
		t.FailNow()

	}
}

func Test_DELETE_A_Purchase(t *testing.T) {

	server := getServer(NewPurchaseService(getDB(TESTDB)))
	defer server.Close()

	res, err := httpPost(user1, getURL(server.URL, user1), postPurchases)

	if !isHTTPStatus(http.StatusCreated, res, err) {
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusCreated)
		t.FailNow()
	}

	purchaseToDelete := getURL(server.URL, user1) + "/" + fmt.Sprintf("%d", timeToTest.Unix())

	res, err = httpDelete(user1, purchaseToDelete)

	if !isHTTPStatus(http.StatusOK, res, err) {
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusOK)
		t.FailNow()
	}

	res, err = httpGet(user1, getURL(server.URL, user1)+"?groupBy=month")

	if !isHTTPStatus(http.StatusOK, res, err) {
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

	var purchasesByMonth []PurchasesByMonth = purchases.PurchasesByMonth

	for _, purchases := range purchasesByMonth {

		for _, p := range purchases.Purchases {
			if strings.Compare(p.Id, fmt.Sprintf("%d", timeToTest.Unix())) == 0 {
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

/*
func Test_items_ids_are_generated_from_their_trimmed_and_lower_case_description(t *testing.T){


	 itemsIds := [...]string { "83c379262dd8fc10dea3ebf7097e12ae7a8dff06",
		              	   "71714c4009f19de8e20a6df8f7a201bdf989af5f",
	                           "b70ba6c3070f343131c1f646c41b1aca0c2ea11f",
	 	                   "7d45da213c946480619093d1eea4e7bd402a77b9"}



	containsId := func (itemsDescriptions []ItemDescription, valueToFind string) bool{
		for _, v := range itemsDescriptions{
			if strings.Compare(v.ItemId, valueToFind) == 0 {
				return true
			}
		}

		return false
	}

	purchases := addPurchasesIds(postPurchases.Purchases)
	items := getItemsDescriptions(purchases)


	log.Printf("%s",items)
	if len(items) != 4 {
		log.Printf("Items expected 3, obtained: %d", len(items))
		t.FailNow()
	}

	for _, id := range itemsIds {

		if !containsId(items, id) {
			log.Printf("Val expected but not prsesent")
			t.FailNow()
		}
	}

}*/

func Test_that_items_descriptions_are_being_saved(t *testing.T) {

	service := NewPurchaseService(getDB(TESTDB))
	server := getServer(service)
	defer server.Close()

	res, err := httpPost(user1, getURL(server.URL, user1), postPurchases)

	if !isHTTPStatus(http.StatusCreated, res, err) {
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusCreated)
		t.FailNow()
	}

	time.Sleep(1000 * time.Millisecond)

	res, err = httpGet(user1, server.URL+"/catalog/users/"+user1+"/items")

	if err != nil {
		log.Printf("Error %s", err.Error())
		t.FailNow()
	}

	itemsDescriptions := new([]ItemDescription)

	body, err := ioutil.ReadAll(res.Body)
	if err := json.Unmarshal(body, itemsDescriptions); err != nil {

		log.Printf("Error when reading response %s", err)
		t.FailNow()
	}

	count := 0

	for range *itemsDescriptions {
		count++
	}

	if count != 4 {
		log.Printf("Expected size 4 but %d", count)
	}
	containsDescriptions := func(itemsDescriptions []ItemDescription, valueToFind string) bool {
		for _, v := range itemsDescriptions {
			if strings.Compare(v.Description, valueToFind) == 0 {
				return true
			}
		}

		return false
	}

	for _, purchase := range postPurchases.Purchases {
		for _, item := range purchase.Items {
			if !containsDescriptions(*itemsDescriptions, item.Description) {
				log.Printf("Items descriptions was not saved")
				t.FailNow()
			}
		}
	}

	log.Printf("%s", body)
}

func Test_that_purchases_are_returned_between_dates(t *testing.T) {

	service := NewPurchaseService(getDB(TESTDB))
	server := getServer(service)
	defer server.Close()


	res, err := httpPost(user1, getURL(server.URL, user1), postPurchases)

	if !isHTTPStatus(http.StatusCreated, res, err) {
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusCreated)
		t.FailNow()
	}

        dateToTest, _ := time.Parse(time.RFC3339, tt)

        from := getDateParam("from", dateToTest)

        dateToTest = dateToTest.AddDate(0,1,0)

        to := getDateParam("to", dateToTest)

        url := fmt.Sprintf("%s?%s&%s", getURL(server.URL, user1), from, to)

	res, err = httpGet(user1, url)

        if !isHTTPStatus(http.StatusOK, res, err) {
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
        log.Printf("purchases :%s", purchases.PurchasesByMonth)
        if len(purchases.PurchasesByMonth) != 1 {
                log.Printf("Error: Expected items quantity is different from the received one: %d", len(purchases.PurchasesByMonth))
                t.FailNow()
        }
        count := 0
        for _, ps := range purchases.PurchasesByMonth {

                for range ps.Purchases {
                      count++
              }
        }

        if count != 2 {
                log.Printf("Error: Expected items quantity is different from the received one: %d", count)
                t.FailNow()
        }
}


func getDateParam(paramName string, date time.Time) string {
        return fmt.Sprintf(paramName + "=%d-%0.2d-%0.2d", date.Year(),date.Month(), date.Day())
}



func Test_that_amount_month_avg_is_returned_by_user(t *testing.T) {

	service := NewPurchaseService(getDB(TESTDB))
	server := getServer(service)
	defer server.Close()



	res, err := httpPost(user1, getURL(server.URL, user1), postPurchases)

	if !isHTTPStatus(http.StatusCreated, res, err) {
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusCreated)
		t.FailNow()
	}


	res, err = httpGet(user1, getURL(server.URL, user1) + "/metrics")

	if !isHTTPStatus(http.StatusOK, res, err) {
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusOK)
		t.FailNow()
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal("Error")
		t.FailNow()
	}

	metrics := new(Metrics)

	if err := json.Unmarshal(body, metrics); err != nil {

		log.Printf("Error when reading response %s", err)
		t.FailNow()
	}

	if err !=  nil || metrics.Accumulated < 8891.94 || metrics.Accumulated >= 8891.95 || metrics.Month_avg < 2222.97 || metrics.Month_avg >= 2222.99 {
		log.Printf("Err: %s or metrics seem not to be the expected ones: avg:%f acc:%f",err, metrics.Month_avg, metrics.Accumulated)
		t.FailNow()
	}

	log.Printf("metrics avg:%0.2f accum:%0.2f", metrics.Month_avg, metrics.Accumulated)
}


//For the moment there is not a more practical way to use, later,
//the user email as ID in DB. So, what I'm doing is to add it in a http header :(

func getURL(url string, userid string) string {
	return url + "/catalog/users/" + userid + "/purchases"
}
