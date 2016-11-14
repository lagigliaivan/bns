package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
	"github.com/gorilla/mux"
	"crypto/sha1"
	"io"
	"net/url"
        "strconv"
)

type Service interface {
	ConfigureRouter(router *mux.Router)
}

func getPathParams(r *http.Request) map[string]string {
	return mux.Vars(r)
}

const (
        GROUP_BY = "groupBy"
        ORDER_ASC = "orderAsc"
        ORDER_DESC = "orderDesc"
        DATE_FROM = "from"
        DATE_TO = "to"

        //Do not change the order.
        NV = iota
        GB
        OA
        OD
        DF
        DT

)

var queryParams map[string]int = make(map[string]int)

func init() {

        log.SetFlags(log.LstdFlags | log.Lshortfile)

        queryParams[GROUP_BY] = GB
        queryParams[ORDER_ASC] = OA
        queryParams[ORDER_DESC] = OD
        queryParams[DATE_FROM] = DF
        queryParams[DATE_TO] = DT

}


type GetPathParams func(*http.Request) map[string]string

type PurchaseService struct {
	getRequestParameters GetPathParams
	error                string
	name                 string
	db                   DB
	purchasesHandler     map[string]func(http.ResponseWriter, *http.Request)

}

func NewPurchaseService(db DB) *PurchaseService {

	service := new(PurchaseService)
	service.getRequestParameters = getPathParams
	service.db = db
	service.error = "ERROR"

        return service
}

//This method sets what resources are going to be managed by the router
func (service PurchaseService) ConfigureRouter(router *mux.Router) {

	routes := Routes{

		Route{
			"get_purchases",
			"GET",
			"/users/{userid}/purchases",
			service.handleGetPurchases,
		},
                Route{
                        "get_purchases",
                        "GET",
                        "/users/{userid}/purchases/metrics",
                        service.handleGetMetrics,
                },
		Route{
			"get_purchase_by_id",
			"GET",
			"/users/{userid}/purchases/{id}",
			service.handleGetPurchaseById,
		},
		Route{
			"post_purchases",
			"POST",
			"/users/{userid}/purchases",
			service.handlePostPurchases,
		},
		Route{
			"delete_purchase",
			"DELETE",
			"/users/{userid}/purchases/{id}",
			service.handleDeletePurchase,
		},
		Route{
			"get_items_description",
			"GET",
			"/users/{userid}/items",
			service.handleGetItemsDescription,
		},
	}

	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}
}

func (service PurchaseService) handleDefaultError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	fmt.Fprint(w, "The request method is not supported for the requested resource")
}

func (service PurchaseService) handleGetPurchases(w http.ResponseWriter, r *http.Request) {

	user := r.Header.Get(USER_ID)
	params := r.URL.Query()

	year := time.Now().Year()

        /*
	handlerId := isPresent(params, ORDER_ASC) | isPresent(params, ORDER_DESC) | isPresent(params, DATE_FROM) | isPresent(params, DATE_TO) | isPresent(params, GROUP_BY)

        handler := getHandler(handlerId)

        handler()
        log.Printf("Handler:%d", handler)
        */

	if paramIsPresent(params, GROUP_BY) {

		from := getDefaultDateFrom(year)
		to := getDefaultDateTo(year)

		purchases := service.getPurchases(user, from, to)

		purchasesSortedByMonth := sortPurchasesByMonth(purchases)

		pByMonthContainer, _ := getPurchasesByMonthContainer(purchasesSortedByMonth)
		purchasesAsJson, err := json.Marshal(pByMonthContainer)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Error %s", err)
			return
		}

		fmt.Fprintf(w, "%s", purchasesAsJson)

	} else if paramIsPresent(params, DATE_FROM) && dateFormatValid(getParam(params, DATE_FROM) + "T00:00:00Z"){

		from := getParam(params, DATE_FROM) + "T00:00:00Z"
		var to string

		if paramIsPresent(params, DATE_TO) && dateFormatValid(getParam(params, DATE_TO) + "T23:59:00Z"){
			to = getParam(params, DATE_TO) + "T23:59:00Z"
		}else {
			to = getDefaultDateTo(year)
		}

		purchases := service.getPurchases(user, from, to)

		purchasesSortedByMonth := sortPurchasesByMonth(purchases)

		pByMonthContainer, _ := getPurchasesByMonthContainer(purchasesSortedByMonth)
		purchasesAsJson, err := json.Marshal(pByMonthContainer)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Error %s", err)
			return
		}

		fmt.Fprintf(w, "%s", purchasesAsJson)

	} else {

		from := getDefaultDateFrom(year)
		to := getDefaultDateTo(year)

		purchases := service.getPurchases(user, from, to)

		container := getPurchaseContainer(purchases)

		purchasesAsJson, err := json.Marshal(container)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("Error")
			return
		}
		fmt.Fprintf(w, "%s", purchasesAsJson)
	}
}



/*func isPresent(params url.Values, value string) int {

        intRepresentation := NV

        param := params[value]

        if param != nil && len(param) != 0 {
                if strings.Compare(param[0], "") != 0 {
                        intRepresentation = queryParams[value]
                }
        }
        return intRepresentation
}*/

func getPurchaseContainer(purchases []Purchase) PurchaseContainer {
        container := NewPurchaseContainer()

        for _, purchase := range purchases {
                container.Add(purchase)
        }

        return container
}

func (service PurchaseService) handleGetMetrics(w http.ResponseWriter, r *http.Request) {

        user := r.Header.Get(USER_ID)

        from := fmt.Sprintf("%d%s", time.Now().Year(), "-01-00T00:00:00Z")
        to := time.Now().Format(time.RFC3339)

        metrics := service.getMetrics(user,from, to)

        metricsAsJson, err := json.Marshal(metrics)

        if err != nil {
                w.WriteHeader(http.StatusBadRequest)
                log.Printf("Error while marshalling GetPurchase response")
                return
        }

        fmt.Fprintf(w, "%s", metricsAsJson)

}

func (service PurchaseService) handleGetPurchaseById(w http.ResponseWriter, r *http.Request) {

	user := r.Header.Get(USER_ID)

	vars := mux.Vars(r)
	id := vars["id"]

	purchase := service.getPurchase(user, id)

	if purchase.Id == "" {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("Purchase id: %s not found", id)
		return
	}

	purchaseAsJson, err := json.Marshal(purchase)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error while marshalling GetPurchase response")
		return
	}

	fmt.Fprintf(w, "%s", purchaseAsJson)
}

func  getPurchasesByMonthContainer(purchasesSortedByMonth  map[time.Month][]Purchase) (PurchasesByMonthContainer, error){

	pByMonthContainer := PurchasesByMonthContainer{PurchasesByMonth: make([]PurchasesByMonth, 0)}
	pByMonth := PurchasesByMonth{}

	for month, purchases := range purchasesSortedByMonth {
		pByMonth.Month = month.String()
		pByMonth.Purchases = purchases
		pByMonthContainer.PurchasesByMonth = append(pByMonthContainer.PurchasesByMonth, pByMonth)
	}

	return pByMonthContainer, nil

}

func (service PurchaseService) handlePostPurchases(w http.ResponseWriter, r *http.Request) {

	user := r.Header.Get(USER_ID)
	body, _ := ioutil.ReadAll(r.Body)

	purchasesContainer := new(PurchaseContainer)

	if err := json.Unmarshal(body, purchasesContainer); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("POST items. The request contains a wrong format %s", err)
		return
	}

	purchases := addPurchasesIds(purchasesContainer.Purchases)

	//TODO: What if savePurchases fails? Where are we handling the error?
	service.savePurchases(user, purchases)

	go service.saveItemsDescriptions(user, purchases)

	w.WriteHeader(http.StatusCreated)
}

func (service PurchaseService) handleDeletePurchase(w http.ResponseWriter, r *http.Request) {

	user := r.Header.Get(USER_ID)
	vars := getPathParams(r)

	itemId := vars["id"]

	log.Printf("Deleting item %s", itemId)
	service.db.DeletePurchase(user, itemId)
}

func (service PurchaseService) getPurchases(userId string, from string, to string) []Purchase {
	log.Printf("Getting items from DB")
	purchases := service.db.GetPurchases(userId, from, to)
	return purchases
}

func (service PurchaseService) getPurchase(userId string, purchaseId string) Purchase {

	log.Printf("Getting purchase from DB")
	purchase := service.db.GetPurchase(userId, purchaseId)
	return purchase
}

func sortPurchasesByMonth(purchases []Purchase) map[time.Month][]Purchase {

	log.Printf("Getting purchases from DB")

	purchasesByMonth := groupPurchasesByMonth(purchases)

	keys := make([]int, 0, len(purchasesByMonth))

	for key := range purchasesByMonth {
		keys = append(keys, int(key))
	}
	sort.Ints(keys)

	sortedPurchases := make(map[time.Month][]Purchase, len(keys))

	for _, key := range keys {
		sortedPurchases[time.Month(key)] = purchasesByMonth[time.Month(key)]
	}

	return sortedPurchases
}

func (service PurchaseService) savePurchases(userId string, purchases []Purchase) {

	log.Printf("Saving items in  DB")

	for _, purchase := range purchases {

                purchase.TotalAmount = calculateTotalAmount(purchase)
		service.db.SavePurchase(purchase, userId)
	}
}

func (service PurchaseService) saveItemsDescriptions(userId string, purchases []Purchase) {

	items_descriptions := getItemsDescriptions(purchases)

	service.db.SaveItemsDescriptions(userId, items_descriptions)
}

func (service PurchaseService) handleGetItemsDescription(w http.ResponseWriter, r *http.Request) {

	user := r.Header.Get(USER_ID)

	itemsDescriptions, _ := service.db.GetItemsDescriptions(user)

	itemsDescriptionsAsJson, err := json.Marshal(itemsDescriptions)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error when creating the reponse, I'm sorry. %s", err)
	}

	fmt.Fprintf(w, "%s", itemsDescriptionsAsJson)
}


func (service PurchaseService) getMetrics(userId string, from string, to string) Metrics{


        purchases := service.getPurchases(userId, from, to)

        sortedPurchases := sortPurchasesByMonth(purchases)

        accumulated := calculateTotalAccumulated(purchases)

        avg := accumulated / float64(len(sortedPurchases))

        avgTrunc, _ := strconv.ParseFloat(fmt.Sprintf("%0.02f", avg), 32)
        accumulatedTrunc, _ := strconv.ParseFloat(fmt.Sprintf("%0.02f", accumulated), 32)

        metrics := Metrics{Month_avg: avgTrunc, Accumulated: accumulatedTrunc}


        return metrics
}

func calculateTotalAccumulated(purchases []Purchase) float64 {

        total := 0.0

        for _, purchase := range purchases{
                total += purchase.TotalAmount

        }

        return total
}

func calculateTotalAmount(purchase Purchase) float64 {

        total := 0.0

        for _, item := range purchase.Items {
                total+=item.Price
        }

        return total
}

func getItemsDescriptions(purchases []Purchase) []ItemDescription {

	items_descriptions := make(map[string]string)
	itemsDescriptions := []ItemDescription{}

	for _, purchase := range purchases {

		for _, item := range purchase.Items {
			items_descriptions[item.Id] = strings.ToLower(item.Description)
		}
	}

	for k, v := range items_descriptions {
		itemsDescriptions = append(itemsDescriptions, ItemDescription{ItemId: k, Description: v})
	}

	return itemsDescriptions
}

func addPurchasesIds(purchases []Purchase) []Purchase {

	identifiable := purchases

	for k, purchase := range identifiable {

		if strings.Compare(purchase.Id, "") == 0 {

			if purchase.Time.IsZero() {
				purchase.Time = time.Now()
			}

			identifiable[k].Id = fmt.Sprintf("%d", purchase.Time.UTC().Unix())
		}

		identifiable[k].Time = purchase.Time.UTC()

		for k, item := range purchase.Items {
			purchase.Items[k].Id = trimAndSha(item.Description)
		}
	}

	return identifiable
}

func trimAndSha(value string) string {

	sha := sha1.New()
	defer sha.Reset()

	//trim and remove spaces
	trimmedAndLowDescription := strings.Replace(strings.TrimSpace(value), " ", "", -1)
	// convert to lower case
	trimmedAndLowDescription = strings.ToLower(trimmedAndLowDescription)
	io.WriteString(sha, trimmedAndLowDescription)

	return fmt.Sprintf("%x", sha.Sum(nil))
}

func paramIsPresent(params url.Values, param string) bool {

	if params != nil && len(params) != 0 {
		p := params[param]

		if len(p) != 0 {
			return true
		}
	}
	return false
}

func getParam(params url.Values, param string) string {

	var value string = ""

	if params != nil && len(params) != 0 {
		p := params[param]
		if len(p) != 0 {
			value = p[0]
		}
	}

	return value
}

func dateFormatValid(date string) bool{

	if _, err := time.Parse(time.RFC3339, date); err != nil {

		return false
	}

	return true
}

func getDefaultDateFrom(year int) string{
	return fmt.Sprintf("%d%s", year , "-01-00T00:00:00Z")

}

func getDefaultDateTo(year int) string {
	return fmt.Sprintf("%d%s", year, "-12-31T23:59:00Z")
}
/**
Given a group or purchases, then they are grouped by month and returned, but not sorted.
 */
func groupPurchasesByMonth( purchases []Purchase ) map[time.Month][]Purchase{

	purchasesByMonth := make(map[time.Month][]Purchase)


	for _, purchase := range purchases {

		if purchasesByMonth[purchase.Time.Month()] == nil {
			purchasesByMonth[purchase.Time.Month()] = make([]Purchase,0)
		}
		purchasesByMonth[purchase.Time.Month()] = append(purchasesByMonth[purchase.Time.Month()], purchase)
	}

	return purchasesByMonth
}
