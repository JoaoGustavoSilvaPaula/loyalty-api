package models

type ProductInvoice struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	Quantity string `json:"quantity"`
	Unit     string `json:"unit"`
	Value    string `json:"value"`
}
