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
	MONTH = "month"
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

	var purchasesByMonth map[time.Month][]Purchase


	purchasesByMonth = service.getPurchasesByMonth(user, year)

	pByMonthContainer := PurchasesByMonthContainer{make([]PurchasesByMonth, 0)}
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

	purchases := new(PurchaseContainer)

	if err := json.Unmarshal(body, purchases); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("POST items. The request contains a wrong format %s", err)
		return
	}

	service.savePurchases(purchases.Purchases, user)
	w.WriteHeader(http.StatusCreated)
}


func (service PurchaseService) handleDeletePurchase (w http.ResponseWriter, r *http.Request) {

	service.db.DeletePurchase()
}

func (service PurchaseService) getPurchases(userId string) []Purchase {
	log.Printf("Getting items from DB")
	purchases := service.db.GetPurchases(userId)
	return  purchases;
}

func (service PurchaseService) getPurchasesByMonth(user string, year int) map[time.Month][]Purchase {

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

type ItemsService struct {
	GetRequestParameters GetPathParams
	error                string
	name                 string
	db                   DB
	productIdsHandler    map[string] func(http.ResponseWriter,*http.Request)
	productsHandler	     map[string] func(http.ResponseWriter,*http.Request)
}

func NewItemService(db DB) *ItemsService {

	service := new(ItemsService)
	service.GetRequestParameters = getPathParams
	service.db = db
	service.error = "ERROR"

	service.productIdsHandler = make(map[string] func(http.ResponseWriter,*http.Request))
	service.productIdsHandler[http.MethodGet] = service.handleGetItem
	service.productIdsHandler[http.MethodPut] = service.handlePutItem
	service.productIdsHandler[service.error]  = service.handleError

	service.productsHandler = make(map[string] func(http.ResponseWriter,*http.Request))
	service.productsHandler[http.MethodPost] = service.handlePostItem
	service.productsHandler[http.MethodGet] = service.handleGetItems
	service.productsHandler[service.error]  = service.handleError

	return service
}

func (service ItemsService) ConfigureRouter(router Router) {

	router.HandleFunc("/products/{id}", service.handleRequestProductId).Name("products+id")
	router.HandleFunc("/products", service.handleRequestProducts).Name("products no slash")
	router.HandleFunc("/products/", service.handleRequestProducts).Name("products + slash")
}

func (service ItemsService) handleError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprint(w, "The request method is not supported for the requested resource")
}

//Handle request of type GET and PUT against /catalog/products/{id}
//This method derives to another different handler according to the HTTP method.
func (service ItemsService) handleRequestProductId(w http.ResponseWriter, r *http.Request){

	handler := service.productIdsHandler[r.Method]
	if handler == nil {
		service.productIdsHandler[service.error] (w, r)
	}else {
		handler(w, r)
	}
}


//Handle request of type GET and PUT against /catalog/products/{id}
//This method derives to another different handler according to the HTTP method.
func (service ItemsService) handleRequestProducts(w http.ResponseWriter, r *http.Request){

	handler := service.productsHandler[r.Method]

	if handler == nil {
		service.productsHandler[service.error] (w, r)
	}else {
		handler(w, r)
	}
}


//URL catalog/products/{id}
func (service ItemsService) handleGetItem(w http.ResponseWriter, r *http.Request) {

	vars := service.GetRequestParameters(r)

	prodId := vars["id"]
	item := service.getItem(prodId)

	if item.Id == "" {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("GET item_id: %s not found", prodId)
		return
	}

	strB, _ := json.Marshal(item)

	fmt.Fprintf(w, "%s", strB)
	log.Printf("GET item_id: %s returned OK", item.Id)

}

func (service ItemsService) handleGetItems(w http.ResponseWriter, r *http.Request) {

	items := service.getItems()

	container := NewItemContainer()

	for _, item := range items {
		container.Add(item)
	}

	itemsAsJson, _ := json.Marshal(container)

	fmt.Fprintf(w, "%s", itemsAsJson)
	log.Printf("GET items returned OK %s", itemsAsJson)

}

func (service ItemsService) handlePutItem(w http.ResponseWriter, r *http.Request){

	vars := service.GetRequestParameters(r)
	itemId := vars["id"]

	if service.getItem(itemId).IsEmpty(){
		w.WriteHeader(http.StatusNotFound)
		log.Printf("PUT itemId: %s Not found", itemId)
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatalf("PUT itemId %s : Problem while reading body: %s Body: %s",itemId, err, body)
		return
	}

	item := new(Item)
	if err := json.Unmarshal(body, item); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("PUT itemId %s. The request contains a wrong format: %s Body: %s", itemId, err, body)
		return
	}
	item.Id = itemId
	service.addUpdateItem(*item)
	w.WriteHeader(http.StatusOK)

}

func (service ItemsService) handlePostItem(w http.ResponseWriter, r *http.Request){

	body, _ := ioutil.ReadAll(r.Body)

	items := new(ItemContainer)

	if err := json.Unmarshal(body, items); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("POST items. The request contains a wrong format %s", err)
		return
	}

	for _, item := range items.GetItems() {

		if item.Id == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("id must not be empty."))
			return
		}

		if service.getItem(item.Id).IsNOTEmpty() {
			w.WriteHeader(http.StatusForbidden)
			log.Printf("POST itemId: %s Already exists", item.Id)
			w.Write([]byte("Id already exists"))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		service.addUpdateItem(item)
	}

	w.WriteHeader(http.StatusCreated)

}

func (service ItemsService) getItem(id string) Item {
	log.Printf("Getting item_id: %s from DB", id)
	item := service.db.GetItem(id)
	return  item;
}

func (service ItemsService) getItems() []Item {
	log.Printf("Getting items from DB")
	items := service.db.GetItems()
	return  items;
}

func (service ItemsService) addUpdateItem(item Item) int {

	if item.Id == "" {
		log.Printf("Error at trying to save an empty item.")
		return -1
	}

	service.db.SaveItem(item)
	log.Printf("PUT item_id: %s returned OK", item.Id)

	return 0
}