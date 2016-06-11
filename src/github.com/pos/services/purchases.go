package services

import (
	"fmt"
	"net/http"
	"github.com/pos/infrastructure"
	"log"
	"github.com/pos/dto"
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

type PurchaseService struct {
	getRequestParameters GetPathParams
	error                string
	name                 string
	db                   infrastructure.DB
	productIdsHandler    map[string] func(http.ResponseWriter,*http.Request)
	purchasesHandler     map[string] func(http.ResponseWriter,*http.Request)
}

func NewPurchaseService(db infrastructure.DB) *PurchaseService {

	service := new(PurchaseService)
	service.getRequestParameters = getPathParams
	service.db = db
	service.error = "ERROR"

	service.purchasesHandler = make(map[string] func(http.ResponseWriter,*http.Request))
	service.purchasesHandler[http.MethodGet] = service.handleGetPurchases
	service.purchasesHandler[http.MethodPost] = service.handlePostPurchases
	service.purchasesHandler[service.error]  = service.handleDefaultError

	return service
}
//This method sets what resources are going to be managed by the router
func (service PurchaseService) ConfigureRouter(router Router) {
	router.HandleFunc("/purchases", service.handleRequestPurchases)
}

func (service PurchaseService) handleDefaultError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprint(w, "The request method is not supported for the requested resource")
}

func (service PurchaseService) handleForbiddenError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
	fmt.Fprint(w, "The request method is not supported for the requested resource")
}

func (service PurchaseService) handleRequestPurchases(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()
	/*if (len(params) == 0 || params["token"] == nil) {
		service.handleForbiddenError(w, r)
		return
	}*/

	if len(params) != 0 {
		if params[GROUP_BY] != nil {
			service.handleGetPurchasesGroupByMonth(w, r)
		}
	}else {
		handler := service.purchasesHandler[r.Method]
		if handler == nil {
			service.purchasesHandler[service.error](w, r)
		}else {
			handler(w, r)
		}
	}
	/*log.Printf("len == 0")
	if len(params) != 0 {
		log.Printf("len != 0")
		for key, _ := range params {
			log.Printf("key: %s", key)
			if key == GROUP_BY {
				service.handleGetPurchasesGroupByMonth(w, r)
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
	}*/
}

func (service PurchaseService) handleGetPurchases(w http.ResponseWriter, r *http.Request) {

	container := dto.NewPurchaseContainer()
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

func (service PurchaseService) handleGetPurchasesGroupByMonth(w http.ResponseWriter, r *http.Request) {

	user := getPathParams(r)["user"]

	var purchasesByMonth map[time.Month][]dto.Purchase

	purchasesByMonth = service.getPurchasesGroupedBy(user, MONTH)

	pByMonthContainer := dto.PurchasesByMonthContainer{make([]dto.PurchasesByMonth, 0)}
	pByMonth := dto.PurchasesByMonth{}

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

func (service PurchaseService) handlePostPurchases(w http.ResponseWriter, r *http.Request) {

	body, _ := ioutil.ReadAll(r.Body)

	purchases := new(dto.PurchaseContainer)

	if err := json.Unmarshal(body, purchases); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("POST items. The request contains a wrong format %s", err)
		return
	}

	service.savePurchases(purchases.Purchases)
	w.WriteHeader(http.StatusCreated)
}

func (service PurchaseService) getPurchases() []dto.Purchase {
	log.Printf("Getting items from DB")
	purchases := service.db.GetPurchases()
	return  purchases;
}

func (service PurchaseService) getPurchasesGroupedBy(user, period string) map[time.Month][]dto.Purchase {

	log.Printf("Getting purchases from DB")
	purchases := service.db.GetPurchasesGroupedByMonth()
	keys := make([]int, 0, len(purchases))

	for key := range purchases {
		keys = append(keys, int(key))
	}
	sort.Ints(keys)

	sortedPurchases := make(map[time.Month][]dto.Purchase, len(keys))

	for _,key := range keys {
		sortedPurchases[time.Month(key)] = purchases[time.Month(key)];
	}
	log.Printf("a verrr: %s", sortedPurchases);
	return  sortedPurchases;
}

func (service PurchaseService) savePurchases( purchases []dto.Purchase)  {
	log.Printf("Saving items in  DB")

	for _, purchase := range purchases {
		service.db.SavePurchase(purchase)
	}
}

func (service PurchaseService) addUpdateItem(item dto.Item) int {

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
/*


func getUser(r *http.Request) string {
	return getPathParams(r)["user"]
}

func getUserId(header http.Header) string{
	s := strings.SplitN(header.Get("Authorization"), " ", 2)
	if len(s) != 2 { return "" }

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil { return "" }

	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 { return "" }

	//return pair[0] == "user" && pair[1] == "pass"
	return pair[0]
}*/
