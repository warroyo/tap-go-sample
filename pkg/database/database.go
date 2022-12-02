package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func ConnectToDB(init bool) *sql.DB {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")

	if !init {
		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s", dbUser, dbPass)+"@tcp("+dbHost+":3306)/"+dbName)
		if err != nil {
			log.Fatal(err)
		}
		return db

	} else {
		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s", dbUser, dbPass)+"@tcp("+dbHost+":3306)/")
		if err != nil {
			log.Fatal(err)
		}
		return db

	}

}
