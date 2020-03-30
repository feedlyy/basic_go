package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func Connect() *sql.DB {
	db, err := sql.Open("mysql",
		"root@tcp(localhost:3306)/basic_go")

	if err != nil {
		log.Fatal(err)
	}

	return db
}
