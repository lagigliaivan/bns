package purchases

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/pos/infrastructure"

	"log"
	"github.com/pos/dto/purchase"
	"github.com/pos/dto/item"
	"encoding/json"
	"io/ioutil"
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

	service.productsHandler = make(map[string] func(http.ResponseWriter,*http.Request))
	service.productsHandler[http.MethodGet] = service.HandleGetPurchases
	service.productsHandler[service.error]  = service.HandleError

	return service
}

func (service Service) HandleError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprint(w, "The request method is not supported for the requested resource")
}

func (service Service) HandleGetPurchases(w http.ResponseWriter, r *http.Request) {

	container := purchase.NewContainer()
	purchases := service.getPurchases()

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
	log.Printf("GET items returned OK %s", purchasesAsJson)

}

func (service Service) HandlePostPurchases (w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)

	purchases := new(purchase.Container)

	if err := json.Unmarshal(body, purchases); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("POST items. The request contains a wrong format %s", err)
		return
	}

	service.savePurchases(purchases.Purchases)
	w.WriteHeader(http.StatusCreated)

}

func (service Service) getPurchases() []purchase.Purchase {
	log.Printf("Getting items from DB")
	purchases := service.db.GetPurchases()
	return  purchases;
}

func (service Service) savePurchases( purchases []purchase.Purchase)  {
	log.Printf("Saving items in  DB")

	for _, purchase := range purchases {
		service.db.SavePurchase(purchase)
	}
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