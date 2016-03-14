package main

import (
	"os"
	"fmt"
	//"html"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/pos/infrastructure"
	"encoding/json"
	"github.com/pos/dto"
	"io/ioutil"
	log "github.com/Sirupsen/logrus"
)

type Service struct {
	error string
	name string
	db infrastructure.DB
	header_handler map[string] func(http.ResponseWriter,*http.Request)
}

func NewService(db infrastructure.DB) *Service{

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stderr)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)

	service := new(Service)
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

//GET or PUT over /catalog/products/{id}
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

	vars := mux.Vars(r)
	prodId := vars["id"]
	item := service.GetItem(prodId)

	if item.Id == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	strB, _ := json.Marshal(item)

	fmt.Fprintf(w, "%s", strB)
	log.WithFields(log.Fields{"itemid":item.Id, "descripton":item.Desc, "price":item.Price}).Debugf("GET item:")

}
//PUT catalog/products/{id}
func (service Service) HandlePutItem(w http.ResponseWriter, r *http.Request){

	body, _ := ioutil.ReadAll(r.Body)

	item := dto.Item{}

	if err := json.Unmarshal(body, &item); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "The request contains a wrong format: %s ", err)
		return
	}

	fmt.Fprintf(w, "%s", item)


}

func (service Service) GetItem(id string) dto.Item {
	log.WithFields(log.Fields{"itemid":id}).Debugf("Get item:")
	item := service.db.GetItem(id)
	return  item;
}

func (service Service) PutItem(id string, desc string, price float32) dto.Item {

	itemDto := dto.Item{id, desc, price};
	service.db.SaveItem(itemDto)
	log.WithFields(log.Fields{"itemid":id, "descripton": desc, "price": price}).Debugf("Saving item:")

	return itemDto
}
