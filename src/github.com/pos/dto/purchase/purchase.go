package purchase

import (
	"time"
	"github.com/pos/dto/item"
	"encoding/json"
	"github.com/pos/dto"
	"reflect"
)

type Purchase struct {
	Time time.Time `json:"time"`//DateTime when this purchase was acquired.
	Items []item.Item `json:"items"`
	Point dto.Point `json:"location"`
	Shop	string `json:"shop"`
}

type Container struct {

	Purchases []Purchase  `json:"purchases"`

}

func NewContainer() Container {

	return Container{Purchases: make([]Purchase, 0)}

}

func (container *Container) Add(purchase Purchase) {

	container.Purchases = append(container.Purchases, purchase)
}

func (container Container) GetPurchases() []Purchase{

	return container.Purchases

}

func (container Container) ToJsonString() string{

	itemsAsJson, _ := json.Marshal(container)
	return string(itemsAsJson)

}

func (container Container) IsEmpty() bool{

	return false
}

func (container Container) GetPurchase(time time.Time) *Purchase {
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
