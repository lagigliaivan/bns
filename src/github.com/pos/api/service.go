package api

import (
	"fmt"
	"html"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/pos/infrastructure"
	"encoding/json"
	"github.com/pos/domain"
)

type Service struct {
	name string
	db infrastructure.CatalogDB
}

//GET catalog/prod/{id}
func (service Service) HandleProducts(w http.ResponseWriter, r *http.Request){

	vars := mux.Vars(r)
	prodId := vars["id"]
	service.process_get_id(prodId)
	fmt.Fprintf(w, "ProductId: %q %q", service.name, html.EscapeString(r.URL.Path))

}

func (service Service) HandleRoot(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Hello, %q %q", service.name, html.EscapeString(r.URL.Path))
}

func (service Service) process_get_id(prodId string) domain.Item{
	item := service.db.GetItem(prodId)
	strB, _ := json.Marshal(item)
	fmt.Printf("ProductId: %s %s", item.GetId(), string(strB))
	return  item;
}