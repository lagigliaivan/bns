package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/pos/infrastructure"

	"github.com/pos/dto"

	"log"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type GetPathParams func (*http.Request) map[string]string

type Service struct {
	GetRequestParameters GetPathParams
	error                string
	name                 string
	db                   infrastructure.DB
	productIdsHandler    map[string] func(http.ResponseWriter,*http.Request)
	productsHandler	     map[string] func(http.ResponseWriter,*http.Request)
}

func NewService(db infrastructure.DB) *Service{

	service := new(Service)
	service.GetRequestParameters = getPathParams
	service.db = db
	service.error = "ERROR"

	service.productIdsHandler = make(map[string] func(http.ResponseWriter,*http.Request))
	service.productIdsHandler[http.MethodGet] = service.HandleGetItem
	service.productIdsHandler[service.error]  = service.HandleError

	service.productsHandler = make(map[string] func(http.ResponseWriter,*http.Request))
	service.productsHandler[http.MethodGet] = service.HandleGetPurchases
	service.productsHandler[service.error]  = service.HandleError

	return service
}

func (service Service) HandleError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprint(w, "The request method is not supported for the requested resource")
}

//Handle request of type GET and PUT against /catalog/products/{id}
//This method derives to another different handler according to the HTTP method.
func (service Service) HandleRequestProductId(w http.ResponseWriter, r *http.Request){

	handler := service.productIdsHandler[r.Method]
	if handler == nil {
		service.productIdsHandler[service.error] (w, r)
	}else {
		handler(w, r)
	}
}


//Handle request of type GET and PUT against /catalog/products/{id}
//This method derives to another different handler according to the HTTP method.
func (service Service) HandleRequestProducts(w http.ResponseWriter, r *http.Request){

	handler := service.productsHandler[r.Method]
	if handler == nil {
		service.productsHandler[service.error] (w, r)
	}else {
		handler(w, r)
	}
}


//URL catalog/products/{id}
func (service Service) HandleGetItem(w http.ResponseWriter, r *http.Request) {

	/*vars := service.GetRequestParameters(r)

	prodId := vars["id"]
	item := service.getItem(prodId)

	if item.Id == "" {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("GET item_id: %s not found", prodId)
		return
	}

	strB, _ := json.Marshal(item)

	fmt.Fprintf(w, "%s", strB)
	log.Printf("GET item_id: %s returned OK", item.Id)*/

}

func (service Service) HandleGetPurchases(w http.ResponseWriter, r *http.Request) {
	/*

	purchases := service.getPurchases()

	for _, purchase := range purchases {
		container.Add(purchase)
	}

	purchasesAsJson, _ := json.Marshal(container)

	fmt.Fprintf(w, "%s", pruchasesAsJson)
	log.Printf("GET items returned OK %s", pruchasesAsJson)
	*/

}


func (service Service) getPurchases() []dto.Purchase {
	log.Printf("Getting items from DB")
	purchases := service.db.GetPurchases()
	return  purchases;
}

func (service Service) addUpdateItem(item dto.Item) int {

	if item.Id == "" {
		log.Printf("Error at trying to save an empty item.")
		return -1
	}

	service.db.SaveItem(item)
	log.Printf("PUT item_id: %s returned OK", item.Id)

	return 0
}

//This function returns a map containing all the path params contained in the request URL.
//In this case, the implementation uses mux.
//This function is used by default, but can be overwritten for testing purposes or any other one.
func getPathParams(r *http.Request) map[string]string {
	return mux.Vars(r)
}