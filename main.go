package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"github.com/warroyo/tap-go-sample/pkg/database"
	"github.com/warroyo/tap-go-sample/pkg/handler"
	"github.com/warroyo/tap-go-sample/pkg/listing"
	_ "modernc.org/sqlite"
)

func main() {

	database.SeedDB()

	router := httprouter.New()

	router.GET("/companies", getCompanies())
	router.ServeFiles("/docs/*filepath", http.Dir("docs"))
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getCompanies() func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	return handler.Writer(listing.GetAllCompanies())
}
