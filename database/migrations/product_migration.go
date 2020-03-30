package migrations

import (
	"basic_go/database"
	"log"
)

func Up() {
	db := database.Connect() //open db connection

	_, err := db.Exec("drop table if exists products")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE products (
	id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	item VARCHAR(100), quantity INT)`)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer db.Close() //close the db connection
}
