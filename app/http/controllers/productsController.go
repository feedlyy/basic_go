package controllers

import (
	"basic_go/app/models"
	"basic_go/database"
	"encoding/json"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	en_translations "gopkg.in/go-playground/validator.v10/translations/en"
	"log"
	"net/http"
	"strconv"
)

func Index(w http.ResponseWriter, r *http.Request) {
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

func Show(w http.ResponseWriter, r *http.Request) {
	var products models.Products
	var arr_products []models.Products
	var response models.Response

	//get the id from url
	params := mux.Vars(r)
	Id := params["id"]

	db := database.Connect()

	rows, err := db.Query("select * from products where id = ?", Id)
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
	defer db.Close()

	//cek datanya
	if arr_products == nil {
		response.Message = "Data tidak ditemukan"
		response.Data = nil
		w.WriteHeader(404)
	} else {
		response.Message = "Success"
		response.Data = arr_products
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func Create(w http.ResponseWriter, r *http.Request) {
	var products models.Products
	var arr_products []models.Products
	var response models.Response

	//form value input
	item := r.FormValue("item")
	quantity := r.FormValue("quantity")

	//translator for validator input
	translator := en.New()
	uni := ut.New(translator, translator)

	// this is usually known or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, found := uni.GetTranslator("en")
	if !found {
		log.Fatal("translator not found")
	}

	v := validator.New()

	if err := en_translations.RegisterDefaultTranslations(v, trans); err != nil {
		log.Fatal(err)
	}

	_ = v.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	parse, _ := strconv.Atoi(quantity)

	check := models.Products{
		Id:       0,
		Item:     item,
		Quantity: parse,
	}

	err := v.Struct(check)

	if err != nil {
		for _, errors := range err.(validator.ValidationErrors) {
			w.WriteHeader(400)
			_ = json.NewEncoder(w).Encode(errors.Translate(trans))
		}
	} else {
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
}

func Update(w http.ResponseWriter, r *http.Request) {
	var products models.Products
	var arr_products []models.Products
	var response models.Response

	//get the id from url
	params := mux.Vars(r)
	Id := params["id"]

	//get form value
	item := r.FormValue("item")
	quantity := r.FormValue("quantity")

	db := database.Connect()

	rows, err := db.Query("select * from products where id = ?", Id)
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

	//validation if id exist or not
	if arr_products == nil {
		response.Message = "Id not found"
		response.Data = nil
		w.WriteHeader(404)
	} else {
		_, err := db.Exec("update products set item = ?, quantity = ? where id = ?",
			item, quantity, Id)
		if err != nil {
			log.Fatal(err)
		}

		rows, err := db.Query("select * from products where id = ?", Id)
		if err != nil {
			log.Fatal(err)
		}

		for rows.Next() {
			if err := rows.Scan(&products.Id, &products.Item, &products.Quantity); err != nil {
				log.Fatal(err.Error())
			} else {
				arr_products = nil //kosongkan dulu dari pengecekan data sebelumnya
				arr_products = append(arr_products, products)
			}
		}
		response.Message = "Success"
		response.Data = arr_products
	}
	defer db.Close()

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	var products models.Products
	var arr_products []models.Products
	var response models.Response

	//get the id from url
	params := mux.Vars(r)
	Id := params["id"]

	db := database.Connect()

	//search for display deleted item
	rows, err := db.Query("select * from products where id = ?", Id)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		if err := rows.Scan(&products.Id, &products.Item, &products.Quantity); err != nil {
			log.Fatal(err)
		} else {
			arr_products = append(arr_products, products)
		}
	}

	//check id given id is exist
	if arr_products == nil {
		response.Message = "Id not found"
		response.Data = nil
		w.WriteHeader(404)
	} else {
		_, err := db.Exec("delete from products where id = ?", Id)
		if err != nil {
			log.Fatal(err)
		}

		response.Message = "Delete Success"
		response.Data = arr_products
	}

	defer db.Close()

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
