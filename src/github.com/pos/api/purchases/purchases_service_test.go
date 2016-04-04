package main

import (
	"testing"

	"net/http"
	"log"

	"github.com/pos/dto"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}


const STATUS_ERROR_MESSAGE string = "%s %s Received status code: %d different from what was expected: %d"

var (
	setOfPurchases = []dto.Purchase{

		{
			Time: time.Now(),
			Items: dto.NewContainer(),
		},

	}

	//postPurchases = dto.PurchasesContainer{Purchases:setOfPurchases}
)

func Test_GET_Purchases_Returns_A_List_Of_Purchases(t *testing.T) {

/*	service := NewService(infrastructure.NewMemDb())
	server := httptest.NewServer(http.HandlerFunc(service.HandleGetPurchases))
	defer server.Close()

	res, err := http.Get(server.URL)
	if !isHTTPStatus(http.StatusOK, res, err){
		log.Printf(STATUS_ERROR_MESSAGE, http.MethodGet, server.URL, res.StatusCode, http.StatusOK)
		t.FailNow()
	}

	purchases := new(dto.Purchase)

	if err := json.Unmarshal(res.Body, purchases); err != nil {

		log.Printf("Error when reading response %s", err)
		t.FailNow()
	}

	if len(purchases) != len(p){
		log.Printf("Error: Expected items quantity is different from the received one")
		t.FailNow()

	}*/
}

func isHTTPStatus(httpStatus int, res *http.Response, err error ) bool {
	return !( (err != nil) || (res.StatusCode != httpStatus) )
}


func httpPOST(service Service) error{
	/*

	server := httptest.NewServer(http.HandlerFunc(service.HandlePostPurchases))
	defer server.Close()

	//POST Items for later being retrieved.
	url := getURLToBeTested(server.URL);

	res, err := httpPost(url, postItems)

	if !isHTTPStatus(http.StatusCreated, res, err){
		debug(http.MethodPost, url, res.StatusCode, http.StatusCreated)
		return err
	}
	*/

	return nil
}