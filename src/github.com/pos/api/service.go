package main

import (
	"fmt"
	//"html"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/pos/infrastructure"
	"encoding/json"
	"github.com/pos/dto"
	"io/ioutil"
)

type Service struct {
	error string
	name string
	db infrastructure.DB
	header_handler map[string] func(http.ResponseWriter,*http.Request)
}

func NewService(db infrastructure.DB) *Service{

	service := new(Service)
	service.header_handler = make(map[string] func(http.ResponseWriter,*http.Request))
	service.db = db
	service.error = "ERROR"
	service.header_handler[http.MethodGet] = service.HandleGetItem
	service.header_handler[http.MethodPut] = service.HandlePutItem
	service.header_handler[service.error] = service.HandleError

	return service
}

func (service Service) HandleError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprint(w, "The request method is not supported for the requested resource")
}

//URL catalog/prod/{id}
func (service Service) HandleRequest(w http.ResponseWriter, r *http.Request){

	handler := service.header_handler[r.Method]
	if handler == nil {
		service.header_handler[service.error] (w, r)
	}else {
		handler(w, r)
	}
}

func (service Service) HandleGetItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	prodId := vars["id"]
	item := service.GetItem(prodId)
	strB, _ := json.Marshal(item)

	fmt.Fprintf(w, "%s", strB)
}
//PUT catalog/prod/{id}
func (service Service) HandlePutItem(w http.ResponseWriter, r *http.Request){

	body, err := ioutil.ReadAll(r.Body)

	fmt.Fprintf(w, "%s %q", body, err)

	//	r.Method
	//vars := mux.Vars(r)
	//prodId := vars["id"]
	//item := service.PutItem(prodId)
	//strB, _ := json.Marshal(item)
	//fmt.Printf("ProductId: %s %s", item.Id, string(strB))
	//fmt.Fprintf(w, "%s", strB)

}

func (service Service) GetItem(id string) dto.Item {
	item := service.db.GetItem(id)
	return  item;
}

func (service Service) PutItem(id string, desc string, price float32) dto.Item {

	itemDto := dto.Item{id, desc, price};
	service.db.SaveItem(itemDto)
	return itemDto
}
