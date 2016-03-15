package main

import (

	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/pos/infrastructure"
	"encoding/json"
	"github.com/pos/dto"
	"io/ioutil"
	"log"
)

type GetPathParams func (*http.Request) map[string]string

type Service struct {
	GetRequestParameters GetPathParams
	error string
	name string
	db infrastructure.DB
	header_handler map[string] func(http.ResponseWriter,*http.Request)
}

func NewService(db infrastructure.DB) *Service{

	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})

	// Output to stderr instead of stdout, could also be a file.
	//log.SetOutput(os.Stderr)

	// Only log the warning severity or above.
	//log.SetLevel(log.DebugLevel)

	service := new(Service)
	service.GetRequestParameters = getPathParams
	service.header_handler = make(map[string] func(http.ResponseWriter,*http.Request))
	service.db = db
	service.error = "ERROR"
	service.header_handler[http.MethodGet] = service.HandleGetItem
	service.header_handler[http.MethodPut] = service.HandlePutItem
	service.header_handler[service.error]  = service.HandleError

	return service
}

func (service Service) HandleError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprint(w, "The request method is not supported for the requested resource")
}

//Handle request of type GET and PUT against /catalog/products/{id}
//This method derives to another different handler according to the HTTP method.
func (service Service) HandleRequest(w http.ResponseWriter, r *http.Request){

	handler := service.header_handler[r.Method]
	if handler == nil {
		service.header_handler[service.error] (w, r)
	}else {
		handler(w, r)
	}
}

//URL catalog/products/{id}
func (service Service) HandleGetItem(w http.ResponseWriter, r *http.Request) {

	vars := service.GetRequestParameters(r)

	prodId := vars["id"]
	item := service.GetItem(prodId)

	if item.Id == "" {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("GET item_id: %s not found", prodId)
		return
	}

	strB, _ := json.Marshal(item)

	fmt.Fprintf(w, "%s", strB)
	log.Printf("GET item_id: %s returned OK", item.Id)

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
	item_id := vars["id"]
	body, _ := ioutil.ReadAll(r.Body)
	item := dto.Item{item_id, "", 0}

	if err := json.Unmarshal(body, &item); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "The request contains a wrong format: %s ", err)
		log.Printf("PUT item_id: %s .The request contains a wrong format %s", item.Id, err)
		return
	}

	service.PutItem(item)
	w.WriteHeader(http.StatusCreated)

}

func (service Service) GetItem(id string) dto.Item {
	log.Printf("Getting item_id: %s from DB", id)
	item := service.db.GetItem(id)
	return  item;
}

func (service Service) PutItem(item dto.Item) int {

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

