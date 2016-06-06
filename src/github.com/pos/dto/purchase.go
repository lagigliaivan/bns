package dto

import (
	"time"
	"encoding/json"
	"reflect"
)

type Purchase struct {
	Time time.Time `json:"time"`//DateTime when this purchase was acquired.
	Items [] Item `json:"items"`
	Point Point `json:"location"`
	Shop string `json:"shop"`
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

func (container PurchaseContainer) GetPurchase(time time.Time) *Purchase {
	for _, p := range container.Purchases{
		if reflect.DeepEqual(p.Time, time) {
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
