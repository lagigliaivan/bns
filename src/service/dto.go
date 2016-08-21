package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)


type Container interface {

	Add(value interface{})
	IsEmpty() bool
}

type Item struct {
	Id          string  `json:"id"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Category    string  `json:"category"`
}


func (item Item) IsEmpty() bool{
	return item.Id == ""
}

func (item Item) IsNOTEmpty() bool{
	return !(item.Id == "")
}

func (item Item) ToJsonString() string{

	itemAsJson, _ := json.Marshal(item)
	return string(itemAsJson)
}


type ItemContainer struct {

	Items []Item  `json:"items"`

}

func NewItemContainer() ItemContainer {

	return ItemContainer{Items: make([]Item, 0)}

}

func (items *ItemContainer) Add(item Item) {

	items.Items = append(items.Items, item)
}

func (items ItemContainer) GetItems() []Item{

	return items.Items

}

func (items ItemContainer) ToJsonString() string{

	itemsAsJson, _ := json.Marshal(items)
	return string(itemsAsJson)

}

func (items ItemContainer) IsEmpty() bool{

	if len(items.GetItems()) < 1 {

		return true
	}

	return false
}


type Purchase struct {
	Id       string `json:"id"`
	Time     time.Time `json:"time"`//DateTime when this purchase was done.
	Items    [] Item `json:"items"`
	Location Point `json:"location"`
	Shop     string `json:"shop"`
}

type PurchaseContainer struct {

	Purchases []Purchase  `json:"purchases"`

}

func NewPurchaseContainer() PurchaseContainer {

	return PurchaseContainer{Purchases: make([]Purchase, 0)}

}

func (container *PurchaseContainer) Add(purchase Purchase) {

	container.Purchases = append(container.Purchases, purchase)
}

func (container PurchaseContainer) GetPurchases() []Purchase{

	return container.Purchases

}

func (container PurchaseContainer) ToJsonString() string{

	itemsAsJson, _ := json.Marshal(container)
	return string(itemsAsJson)

}

func (container PurchaseContainer) IsEmpty() bool{

	return false
}

func (container PurchaseContainer) GetPurchase(id int64) *Purchase {
	for _, p := range container.Purchases{
		if reflect.DeepEqual(p.Time.UTC().Unix(), id) {
			return &p
		}
	}

	return nil
}

type PurchasesByMonth struct {
	Month     string  `json:"month"`
	Purchases []Purchase  `json:"purchases"`
}

type PurchasesByMonthContainer struct {

	PurchasesByMonth []PurchasesByMonth `json:"purchasesByMonth"`
}

type Point struct {
	Lat float64
	Long float64
}

func NewPoint(lat, long float64) Point {
	point := Point{Lat:lat, Long:long}
	return point
}

func (point Point) toString() string {
	return fmt.Sprintf("%s %s", point.Lat, point.Long)

}


type ItemDescription struct {

	ItemId string `json:"itemid"`
	Description string `json:"description"`
}