package main

import (
	"fmt"
	"html"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/pos/infrastructure"
	"encoding/json"
	"github.com/pos/dto"
)

type Service struct {
	name string
	db infrastructure.DB
}

//GET catalog/prod/{id}
func (service Service) HandleGetItem(w http.ResponseWriter, r *http.Request){

	vars := mux.Vars(r)
	prodId := vars["id"]
	item := service.GetItem(prodId)
	strB, _ := json.Marshal(item)
	//fmt.Printf("ProductId: %s %s", item.Id, string(strB))
	fmt.Fprintf(w, "%s", strB)

}

//PUT catalog/prod/{id}
func (service Service) HandlePutItem(w http.ResponseWriter, r *http.Request){

	body := r.Body
	p := make([]byte, 1000)
	body.Read(p)
	//vars := mux.Vars(r)
	//prodId := vars["id"]
	//item := service.PutItem(prodId)
	//strB, _ := json.Marshal(item)
	//fmt.Printf("ProductId: %s %s", item.Id, string(strB))
	//fmt.Fprintf(w, "%s", strB)

}

func (service Service) HandleRoot(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Hello, %q %q", service.name, html.EscapeString(r.URL.Path))
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
