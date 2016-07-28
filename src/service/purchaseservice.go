package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"sort"
	"time"
	"strings"
)

type Service interface {
	ConfigureRouter(router Router)
}

func getPathParams(r *http.Request) map[string]string {
	return mux.Vars(r)
}


func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

const (
	GROUP_BY = "groupBy"
	//MONTH = "month"
)
type GetPathParams func (*http.Request) map[string]string

type PurchaseService struct {
	getRequestParameters GetPathParams
	error                string
	name                 string
	db                   DB
	productIdsHandler    map[string] func(http.ResponseWriter,*http.Request)
	purchasesHandler     map[string] func(http.ResponseWriter,*http.Request)
}

func NewPurchaseService(db DB) *PurchaseService {

	service := new(PurchaseService)
	service.getRequestParameters = getPathParams
	service.db = db
	service.error = "ERROR"

	service.purchasesHandler = make(map[string] func(http.ResponseWriter,*http.Request))
	service.purchasesHandler[http.MethodGet] = service.handleGetPurchases
	service.purchasesHandler[http.MethodPost] = service.handlePostPurchases
	service.purchasesHandler[http.MethodDelete]  = service.handleDeletePurchase
	service.purchasesHandler[service.error]  = service.handleDefaultError


	return service
}
//This method sets what resources are going to be managed by the router
func (service PurchaseService) ConfigureRouter(router Router) {
	router.HandleFunc("/purchases", service.handleRequestPurchases)
	router.HandleFunc("/purchases/{id}", service.handleRequestPurchases)
}

func (service PurchaseService) handleDefaultError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprint(w, "The request method is not supported for the requested resource")
}

func (service PurchaseService) handleRequestPurchases(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()

	if len(params) != 0 {
		if params[GROUP_BY] != nil {
			service.handleGetPurchasesByMonth(w, r)
		}
	}else {
		handler := service.purchasesHandler[r.Method]
		if handler == nil {
			service.purchasesHandler[service.error](w, r)
		}else {
			handler(w, r)
		}
	}
}

func (service PurchaseService) handleGetPurchases(w http.ResponseWriter, r *http.Request) {
	user := r.Header.Get(HEADER)

	container := NewPurchaseContainer()
	purchases := service.getPurchases(user)

	for _, purchase := range purchases {
		container.Add(purchase)
	}

	purchasesAsJson, err := json.Marshal(container)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error")
		return
	}

	fmt.Fprintf(w, "%s", purchasesAsJson)
}

func (service PurchaseService) handleGetPurchasesByMonth(w http.ResponseWriter, r *http.Request) {

	user := r.Header.Get(HEADER)
	//params := r.URL.Query()

	year := time.Now().Year()
/*
	if params["year"] != nil {
		year = params["year"]
	}*/

	var purchasesByMonth map[time.Month] []Purchase


	purchasesByMonth = service.getPurchasesByMonth(user, year)

	pByMonthContainer := PurchasesByMonthContainer{PurchasesByMonth: make([]PurchasesByMonth, 0)}
	pByMonth := PurchasesByMonth{}

	for month, purchases := range purchasesByMonth {
		pByMonth.Month = month.String()
		pByMonth.Purchases = purchases
		pByMonthContainer.PurchasesByMonth = append(pByMonthContainer.PurchasesByMonth,pByMonth)
	}

	purchasesAsJson, err := json.Marshal(pByMonthContainer)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error %s", err)
		return
	}

	fmt.Fprintf(w, "%s", purchasesAsJson)

}

func (service PurchaseService) handlePostPurchases(w http.ResponseWriter, r *http.Request) {

	user := r.Header.Get(HEADER)
	body, _ := ioutil.ReadAll(r.Body)

	purchasesContainer := new(PurchaseContainer)

	if err := json.Unmarshal(body, purchasesContainer); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("POST items. The request contains a wrong format %s", err)
		return
	}

	for k, purchases := range purchasesContainer.Purchases {
		if strings.Compare(purchases.Id, "") == 0 {
			purchasesContainer.Purchases[k].Id = fmt.Sprintf("%d", purchases.Time.Unix())
		}
	}
	service.savePurchases(purchasesContainer.Purchases, user)
	w.WriteHeader(http.StatusCreated)
}


func (service PurchaseService) handleDeletePurchase (w http.ResponseWriter, r *http.Request) {

	user := r.Header.Get(HEADER)
	vars := getPathParams(r)

	itemId := vars["id"]

	log.Printf("Deleting item %s", itemId)
	service.db.DeletePurchase(user, itemId)
}

func (service PurchaseService) getPurchases(userId string) []Purchase {
	log.Printf("Getting items from DB")
	purchases := service.db.GetPurchases(userId)
	return  purchases;
}

func (service PurchaseService) getPurchasesByMonth(user string, year int) map[time.Month] []Purchase {

	log.Printf("Getting purchases from DB")

	purchases := service.db.GetPurchasesByMonth(user, year)
	keys := make([]int, 0, len(purchases))

	for key := range purchases {
		keys = append(keys, int(key))
	}
	sort.Ints(keys)

	sortedPurchases := make(map[time.Month][]Purchase, len(keys))

	for _,key := range keys {
		sortedPurchases[time.Month(key)] = purchases[time.Month(key)];
	}

	return  sortedPurchases;
}

func (service PurchaseService) savePurchases( purchases []Purchase, userId string)  {
	log.Printf("Saving items in  DB")

	for _, purchase := range purchases {
		service.db.SavePurchase(purchase, userId)
	}
}

func (service PurchaseService) addUpdateItem(item Item) int {

	if item.Id == "" {
		log.Printf("Error at trying to save an empty item.")
		return -1
	}

	service.db.SaveItem(item)
	log.Printf("PUT item_id: %s returned OK", item.Id)

	return 0
}