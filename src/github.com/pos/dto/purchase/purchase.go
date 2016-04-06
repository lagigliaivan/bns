package purchase

import (
	"time"
	"github.com/pos/dto/item"
	"encoding/json"
)

type Purchase struct {
	Time time.Time //DateTime when this purchase was acquired.
	Items []item.Item
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

type PurchasesByMonth struct {
	Month     string  `json:"month"`
	Purchases []Purchase  `json:"purchases"`
}

type PurchasesByMonthContainer struct {

	PurchasesByMonth []PurchasesByMonth `json:"purchasesByMonth"`
}
