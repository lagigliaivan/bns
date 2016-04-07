package purchases

import (
	"testing"
	"net/http"
	"log"
	"time"
	"github.com/pos/infrastructure"
	"net/http/httptest"
	"github.com/pos/dto/purchase"
	"encoding/json"
	"github.com/pos/dto/item"
	"io/ioutil"
	"strings"

)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}


const STATUS_ERROR_MESSAGE string = "%s %s Received status code: %d different from what was expected: %d"

var (
	itemsPurchaseA = []item.Item{

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

	itemsPurchaseB = []item.Item{

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

	setOfPurchases = []purchase.Purchase{

		{
			Time: time.Now(),
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

	postPurchases = purchase.Container{Purchases:setOfPurchases}
)

func Test_GET_Purchases_Returns_A_List_Of_Purchases(t *testing.T) {

	service := NewService(infrastructure.NewMemDb())
	server := httptest.NewServer(http.HandlerFunc(service.HandlePostPurchases))
	defer server.Close()

	res, err := httpPost(server.URL, postPurchases)

	if !isHTTPStatus(http.StatusCreated, res, err){
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusCreated)
		t.FailNow()
	}


	server = httptest.NewServer(http.HandlerFunc(service.HandleGetPurchases))
	res, err = http.Get(server.URL)

	if !isHTTPStatus(http.StatusOK, res, err){
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusOK)
		t.FailNow()
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal("Error")
		t.FailNow()
	}

	purchases := new(purchase.Container)

	if err := json.Unmarshal(body, purchases); err != nil {

		log.Printf("Error when reading response %s", err)
		t.FailNow()
	}

	if len(purchases.Purchases) != len(setOfPurchases){
		log.Printf("Error: Expected items quantity is different from the received one")
		t.FailNow()

	}
}

func Test_GET_Purchases_Grouped_By_Month_Returns_A_List_Of_Purchases_Groups(t *testing.T) {

	service := NewService(infrastructure.NewMemDb())
	server := httptest.NewServer(http.HandlerFunc(service.HandlePostPurchases))
	defer server.Close()

	res, err := httpPost(server.URL, postPurchases)

	if !isHTTPStatus(http.StatusCreated, res, err){
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusCreated)
		t.FailNow()
	}

	server = httptest.NewServer(http.HandlerFunc(service.HandleGetPurchasesGroupByMonth))
	res, err = http.Get(server.URL + "?groupBy=month")

	if !isHTTPStatus(http.StatusOK, res, err){
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusOK)
		t.FailNow()
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal("Error")
		t.FailNow()
	}

	purchases := new(purchase.PurchasesByMonthContainer)

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

	service := NewService(infrastructure.NewMemDb())
	server := httptest.NewServer(http.HandlerFunc(service.HandlePostPurchases))
	//server := httptest.NewServer(http.HandlerFunc(main.GetRouter()))
	defer server.Close()

	res, err := httpPost(server.URL, postPurchases)

	if !isHTTPStatus(http.StatusCreated, res, err){
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusCreated)
		t.FailNow()
	}


	server = httptest.NewServer(http.HandlerFunc(service.HandleGetPurchasesGroupByMonth))
	res, err = http.Get(server.URL + "?groupBy=ANYTHING")

	if !isHTTPStatus(http.StatusOK, res, err){
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusOK)
		t.FailNow()
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal("Error")
		t.FailNow()
	}

	purchases := new(purchase.PurchasesByMonthContainer)

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


func isHTTPStatus(httpStatus int, res *http.Response, err error ) bool {
	return !( (err != nil) || (res.StatusCode != httpStatus) )
}

func httpPost(url string, purchases purchase.Container) (*http.Response, error){

	body := strings.NewReader(purchases.ToJsonString())
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		log.Printf("Error when creating POST request %d.", err)
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)

	return resp, err
}