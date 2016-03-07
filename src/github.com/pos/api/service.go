package main

import (
	"fmt"
	"html"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/pos/infrastructure"
	"github.com/pos/domain"
	"encoding/json"

	"github.com/pos/dto"
)

type Service struct {
	name string
	db infrastructure.DB
}

//GET catalog/prod/{id}
func (service Service) HandleProducts(w http.ResponseWriter, r *http.Request){

	vars := mux.Vars(r)
	prodId := vars["id"]
	service.GetItem(prodId)
	fmt.Fprintf(w, "ProductId: %q %q", service.name, html.EscapeString(r.URL.Path))

}

func (service Service) HandleRoot(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Hello, %q %q", service.name, html.EscapeString(r.URL.Path))
}

func (service Service) GetItem(id string) domain.Item {

	item := service.db.GetItem(id)
	strB, _ := json.Marshal(dto.Item{}.GetDto(item))
	fmt.Printf("ProductId: %s %s", item.GetId(), string(strB))
	return  item;
}

func (service Service) PutItem(id string, desc string, price float32) domain.Item {

	item := domain.NewItem(id)
	item.SetDescription(desc)
	item.SetPrice(price)

	itemDto := dto.Item{}.GetDto(*item)

	service.db.Save(itemDto)

}
