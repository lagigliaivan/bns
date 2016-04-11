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
	"time"
	"sort"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

const (
	GROUP_BY = "groupBy"
	MONTH = "month"
)
type GetPathParams func (*http.Request) map[string]string

type Service struct {
	GetRequestParameters GetPathParams
	error                string
	name                 string
	db                   infrastructure.DB
	productIdsHandler    map[string] func(http.ResponseWriter,*http.Request)
	purchasesHandler     map[string] func(http.ResponseWriter,*http.Request)
}

func NewService(db infrastructure.DB) *Service{

	service := new(Service)
	service.GetRequestParameters = getPathParams
	service.db = db
	service.error = "ERROR"

	service.purchasesHandler = make(map[string] func(http.ResponseWriter,*http.Request))
	service.purchasesHandler[http.MethodGet] = service.HandleGetPurchases
	service.purchasesHandler[http.MethodPost] = service.HandlePostPurchases
	service.purchasesHandler[service.error]  = service.HandleError

	return service
}

func (service Service) ConfigureService(router *mux.Router) {
	router.HandleFunc("/catalog/purchases", service.HandleRequestPurchases)
}

func (service Service) HandleError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprint(w, "The request method is not supported for the requested resource")
}

func (service Service) HandleRequestPurchases(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	log.Printf("len == 0")
	if len(params) != 0 {
		log.Printf("len != 0")
		for key, _ := range params {
			log.Printf("key: %s", key)
			if key == GROUP_BY {
				service.HandleGetPurchasesGroupByMonth(w, r)
				break
			}
		}
	} else {
		handler := service.purchasesHandler[r.Method]
		if handler == nil {
			service.purchasesHandler[service.error](w, r)
		}else {
			handler(w, r)
		}
	}
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
}

func (service Service) HandleGetPurchasesGroupByMonth(w http.ResponseWriter, r *http.Request) {

	var purchasesByMonth map[time.Month][]purchase.Purchase

	purchasesByMonth = service.getPurchasesGroupedBy(MONTH)

	pByMonthContainer := purchase.PurchasesByMonthContainer{make([]purchase.PurchasesByMonth, 0)}
	pByMonth := purchase.PurchasesByMonth{}

	for month, purchases := range purchasesByMonth {
		pByMonth.Month = month.String()
		pByMonth.Purchases = purchases
		pByMonthContainer.PurchasesByMonth = append(pByMonthContainer.PurchasesByMonth,pByMonth)
	}

	purchasesAsJson, err := json.Marshal(pByMonthContainer)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error %s", err)
		return
	}

	fmt.Fprintf(w, "%s", purchasesAsJson)

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

func (service Service) getPurchasesGroupedBy(period string) map[time.Month][]purchase.Purchase {
	log.Printf("Getting items from DB")
	purchases := service.db.GetPurchasesGroupedByMonth()
	keys := make([]int, 0, len(purchases))

	for key := range purchases {
		keys = append(keys, int(key))
	}
	sort.Ints(keys)

	sortedPurchases := make(map[time.Month][]purchase.Purchase, len(keys))

	for _,key := range keys {
		sortedPurchases[time.Month(key)] = purchases[time.Month(key)];
	}
	log.Printf("a verrr: %s", sortedPurchases);
	return  sortedPurchases;
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