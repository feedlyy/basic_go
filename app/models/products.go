package models

type Products struct {
	Id       int    `form:"id" json:"id"`
	Item     string `form:"item" json:"item"`
	Quantity int    `form:"quantity" json:"quantity"`
}

type Response struct {
	Message string `json:"message"`
	Data    []Products
}
