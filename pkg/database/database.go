package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/nebhale/client-go/bindings"
)

func ConnectToDB(init bool) *sql.DB {
	const file string = "companies.db"
	var dbUser = ""
	var dbPass = ""
	var dbName = ""
	var dbHost = ""
	var ok = false

	b := bindings.FromServiceBindingRoot()
	b = bindings.Filter(b, "mysql")
	if len(b) != 1 {
		_, _ = fmt.Fprintf(os.Stderr, "Incorrect number of MYSQL drivers: %d\n", len(b))
		log.Println("no service bindings present falling back to ENV vars")
		dbUser = os.Getenv("DB_USER")
		dbPass = os.Getenv("DB_PASSWORD")
		dbName = os.Getenv("DB_NAME")
		dbHost = os.Getenv("DB_HOST")
	} else {
		log.Println("service bindings present")
		dbName, ok = bindings.Get(b[0], "database")
		if !ok {
			log.Fatal("No database name in binding")
		}
		dbHost, ok = bindings.Get(b[0], "host")
		if !ok {
			log.Fatal("No host in binding")
		}
		dbUser, ok = bindings.Get(b[0], "username")
		if !ok {
			log.Fatal("No username in binding")
		}
		dbPass, ok = bindings.Get(b[0], "password")
		if !ok {
			log.Fatal("No password in binding")
		}

	}

	//check if DB params are sent, if not create sqllite
	if dbHost == "" {
		log.Print("using sqllite")
		db, err := sql.Open("sqlite3", file)
		if err != nil {
			log.Fatal(err)
		}
		return db
	} else if !init {

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
