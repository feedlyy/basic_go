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
	defer db.Close()

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

func Create(w http.ResponseWriter, r *http.Request) {
	var products models.Products
	var arr_products []models.Products
	var response models.Response

	item := r.FormValue("item")
	quantity := r.FormValue("quantity")

	db := database.Connect()

	rows, err := db.Exec("insert into products(item, quantity) values (?, ?)", item, quantity)
	if err != nil {
		log.Fatal(err)
	}

	lastId, _ := rows.LastInsertId()

	getLast, err := db.Query("select * from products where id = ?", lastId)
	if err != nil {
		log.Fatal(err)
	}

	for getLast.Next() {
		if err := getLast.Scan(&products.Id, &products.Item, &products.Quantity); err != nil {
			log.Fatal(err.Error())
		} else {
			arr_products = append(arr_products, products)
		}
	}
	defer db.Close()

	response.Message = "Success"
	response.Data = arr_products

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
