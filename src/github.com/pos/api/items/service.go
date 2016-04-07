package items

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/pos/infrastructure"
	"encoding/json"
	"github.com/pos/dto/item"
	"io/ioutil"
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
	service.productIdsHandler[http.MethodPut] = service.HandlePutItem
	service.productIdsHandler[service.error]  = service.HandleError

	service.productsHandler = make(map[string] func(http.ResponseWriter,*http.Request))
	service.productsHandler[http.MethodPost] = service.HandlePostItem
	service.productsHandler[http.MethodGet] = service.HandleGetItems
	service.productsHandler[service.error]  = service.HandleError

	return service
}

func (service Service) ConfigureRouter(router *mux.Router) {

	router.HandleFunc("/catalog/products/{id}", service.HandleRequestProductId)
	router.HandleFunc("/catalog/products", service.HandleRequestProducts)
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

func (service Service) HandleGetItems(w http.ResponseWriter, r *http.Request) {

	items := service.getItems()

	container := item.NewContainer()

	for _, item := range items {
		container.Add(item)
	}

	itemsAsJson, _ := json.Marshal(container)

	fmt.Fprintf(w, "%s", itemsAsJson)
	log.Printf("GET items returned OK %s", itemsAsJson)

}
// @Title Get Users Information
// @Description Get Users Information
// @Accept json
// @Param userId path int true "User ID"
// @Success 200 {object} string "Success"
// @Failure 401 {object} string "Access denied"
// @Failure 404 {object} string "Not Found"
// @Resource /users
// @Router /v1/users/:userId.json [get]//PUT catalog/products/{id}
func (service Service) HandlePutItem(w http.ResponseWriter, r *http.Request){

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

	item := new(item.Item)
	if err := json.Unmarshal(body, item); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("PUT itemId %s. The request contains a wrong format: %s Body: %s", itemId, err, body)
		return
	}
	item.Id = itemId
	service.addUpdateItem(*item)
	w.WriteHeader(http.StatusOK)

}

func (service Service) HandlePostItem(w http.ResponseWriter, r *http.Request){

	body, _ := ioutil.ReadAll(r.Body)

	items := new(item.Container)

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

func (service Service) getItem(id string) item.Item {
	log.Printf("Getting item_id: %s from DB", id)
	item := service.db.GetItem(id)
	return  item;
}

func (service Service) getItems() []item.Item {
	log.Printf("Getting items from DB")
	items := service.db.GetItems()
	return  items;
}

func (service Service) addUpdateItem(item item.Item) int {

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