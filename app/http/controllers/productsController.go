package controllers

import (
	"basic_go/app/models"
	"basic_go/database"
	"encoding/json"
	"log"
	"net/http"
)

func Products(w http.ResponseWriter, r *http.Request) {
	var products models.Products
	var arr_products []models.Products
	var response models.Response

	db := database.Connect()

	rows, err := db.Query("SELECT * from products")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		if err := rows.Scan(&products.Id, &products.Item, &products.Quantity); err != nil {
			log.Fatal(err.Error())
		} else {
			arr_products = append(arr_products, products)
		}
	}

	response.Message = "Success"
	response.Data = arr_products

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
