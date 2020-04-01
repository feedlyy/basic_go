package models

type Products struct {
	Id       int    `form:"id" json:"id"`
	Item     string `form:"item" json:"item" validate:"required"`
	Quantity int    `form:"quantity" json:"quantity" validate:"required"`
}

type Response struct {
	Message string `json:"message"`
	Data    []Products
}
