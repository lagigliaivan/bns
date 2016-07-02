package main



type AWSQueryPurchases struct {

	Count int  `json:"Count"`
	ScannedCount int  `json: ScannedCount`
	Items []AWSPurchase `json: Items`
}

type AWSPurchase struct {

	Date string `json:"dt"`
	Shop string `json: "shop"`
	Items string `json: "items"`
}
