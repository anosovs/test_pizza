package models

type Order struct {
	Order_id string `json:"order_id"`
	Items []int `json:"items"`
	Done bool `json:"done"`
}
