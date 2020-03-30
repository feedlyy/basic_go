package main

import (
	"basic_go/app/http/controllers"
	"basic_go/database/migrations"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	//init router
	r := mux.NewRouter()

	//router list
	r.HandleFunc("/api/products", controllers.Products).Methods("GET")

	//serve a server
	_ = http.ListenAndServe(":8000", r)

	/*to run this migration, hit the go build && basic_go(project name folder) in the terminal
	and run this func main / run the application / ctrl + shift + f10*/
	migrations.Up()
}
