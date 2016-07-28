package main

import (
	"net/http"
	"log"
	"fmt"
	"encoding/json"
	"io/ioutil"
)


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