package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"sort"
	"time"
	"strings"
//	"github.com/aws/aws-sdk-go/service"
	"crypto/sha1"
	"io"
)

type Service interface {
	ConfigureRouter(router *mux.Router)
}

func getPathParams(r *http.Request) map[string]string {
	return mux.Vars(r)
}


func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

const (
	GROUP_BY = "groupBy"
	//MONTH = "month"
)
type GetPathParams func (*http.Request) map[string]string

type PurchaseService struct {
	getRequestParameters GetPathParams
	error                string
	name                 string
	db                   DB
	purchasesHandler     map[string] func(http.ResponseWriter,*http.Request)
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

	if len(params) != 0 && params[GROUP_BY] != nil {

		service.handleGetPurchasesByMonth(w, r)

	}else {
		container := NewPurchaseContainer()
		purchases := service.getPurchases(user)

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

func (service PurchaseService) handleGetPurchasesByMonth(w http.ResponseWriter, r *http.Request) {

	user := r.Header.Get(USER_ID)
	year := time.Now().Year()

	/*if params["year"] != nil {
		year = params["year"]
	}*/

	var purchasesSortedByMonth map[time.Month] []Purchase


	purchasesSortedByMonth = service.getPurchasesSortedByMonth(user, year)

	pByMonthContainer := PurchasesByMonthContainer{PurchasesByMonth: make([]PurchasesByMonth, 0)}
	pByMonth := PurchasesByMonth{}

	for month, purchases := range purchasesSortedByMonth {
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


func (service PurchaseService) handleDeletePurchase (w http.ResponseWriter, r *http.Request) {

	user := r.Header.Get(USER_ID)
	vars := getPathParams(r)

	itemId := vars["id"]

	log.Printf("Deleting item %s", itemId)
	service.db.DeletePurchase(user, itemId)
}

func (service PurchaseService) getPurchases(userId string) []Purchase {
	log.Printf("Getting items from DB")
	purchases := service.db.GetPurchases(userId)
	return  purchases;
}


func (service PurchaseService) getPurchase(userId string, purchaseId string) Purchase {

	log.Printf("Getting purchase from DB")
	purchase := service.db.GetPurchase(userId, purchaseId)
	return  purchase;
}


func (service PurchaseService) getPurchasesSortedByMonth(user string, year int) map[time.Month] []Purchase {

	log.Printf("Getting purchases from DB")

	purchases := service.db.GetPurchasesByMonth(user, year)
	keys := make([]int, 0, len(purchases))

	for key := range purchases {
		keys = append(keys, int(key))
	}
	sort.Ints(keys)

	sortedPurchases := make(map[time.Month][]Purchase, len(keys))

	for _,key := range keys {
		sortedPurchases[time.Month(key)] = purchases[time.Month(key)];
	}

	return  sortedPurchases;
}

func (service PurchaseService) savePurchases( userId string, purchases []Purchase)  {

	log.Printf("Saving items in  DB")

	for _, purchase := range purchases {
		service.db.SavePurchase(purchase, userId)
	}
}

func (service PurchaseService) saveItemsDescriptions(userId string, purchases []Purchase){

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

func getItemsDescriptions( purchases []Purchase) []ItemDescription {

	items_descriptions := make(map[string] string)
	itemsDescriptions := []ItemDescription{}

	for _, purchase := range purchases {

		for _, item := range purchase.Items{
			items_descriptions[item.Id] = strings.ToLower(item.Description)
		}
	}

	for k, v := range items_descriptions{
		itemsDescriptions = append(itemsDescriptions, ItemDescription{ItemId:k, Description:v})
	}


	return itemsDescriptions
}


func addPurchasesIds(purchases []Purchase) []Purchase{

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
			purchase.Items[k].Id = trimAndSha(item.Description);
		}
	}

	return identifiable
}

func trimAndSha(value string) string {

	sha := sha1.New()
	defer sha.Reset()

	//trim and remove spaces
	trimmedAndLowDescription :=  strings.Replace(strings.TrimSpace(value), " ", "", -1)
	// convert to lower case
	trimmedAndLowDescription = strings.ToLower(trimmedAndLowDescription)
	io.WriteString(sha, trimmedAndLowDescription)

	return fmt.Sprintf("%x", sha.Sum(nil))
}