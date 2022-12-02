package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"github.com/warroyo/tap-go-sample/pkg/database"
	"github.com/warroyo/tap-go-sample/pkg/handler"
	"github.com/warroyo/tap-go-sample/pkg/listing"
)

func main() {

	seedDB()

	router := httprouter.New()

	router.GET("/companies", getCompanies())

	log.Fatal(http.ListenAndServe(":8080", router))
}

func getCompanies() func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return handler.Writer(listing.GetAllCompanies())
}

func seedDB() {
	dbname := os.Getenv("DB_NAME")
	db := database.ConnectToDB(true)
	defer db.Close()

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+dbname)
	if err != nil {
		log.Printf("Error %s when creating DB\n", err)
		return
	}
	no, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when fetching rows", err)
		return
	}
	log.Printf("rows affected %d\n", no)

	_, err = db.Exec("USE " + dbname)
	if err != nil {
		log.Printf("Error %s", err)
		return
	}

	_, table_check := db.Query("select * from companies;")

	if table_check == nil {
		log.Printf("table exists not seeding")
	} else {
		res, err = db.ExecContext(ctx, "CREATE TABLE companies ( id integer, name varchar(32) )")
		if err != nil {
			log.Printf("Error %s when creating Table\n", err)
			return
		}

		company := listing.Company{
			Id:   1,
			Name: "acme",
		}

		query := "INSERT INTO `companies` (`id`, `name`) VALUES ( ?, ?)"
		res, err = db.ExecContext(ctx, query, company.Id, company.Name)
		if err != nil {
			log.Printf("Error %s when seeding Table\n", err)
			return
		}
	}

	db.Close()
	db = database.ConnectToDB(true)
	defer db.Close()

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB", err)
		return
	}
	log.Printf("Connected to DB %s successfully\n", dbname)

}
