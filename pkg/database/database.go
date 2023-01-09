package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/nebhale/client-go/bindings"
)

type dbBindings struct {
	dbUser string
	dbName string
	dbPass string
	dbHost string
}

func ConnectToDB(init bool) *sql.DB {
	const file string = "companies.db"
	dbinfo := getBindings()
	//check if DB params are sent, if not create sqllite
	if dbinfo.dbHost == "" {
		log.Print("using sqllite")
		db, err := sql.Open("sqlite", file)
		if err != nil {
			log.Fatal(err)
		}
		return db
	} else if !init {

		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s", dbinfo.dbUser, dbinfo.dbPass)+"@tcp("+dbinfo.dbHost+":3306)/"+dbinfo.dbName)
		if err != nil {
			log.Fatal(err)
		}
		return db

	} else {
		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s", dbinfo.dbUser, dbinfo.dbPass)+"@tcp("+dbinfo.dbHost+":3306)/")
		if err != nil {
			log.Fatal(err)
		}
		return db

	}

}

func getBindings() dbBindings {
	var dbinfo dbBindings
	var ok = false

	b := bindings.FromServiceBindingRoot()
	b = bindings.Filter(b, "mysql")
	if len(b) != 1 {
		_, _ = fmt.Fprintf(os.Stderr, "Incorrect number of MYSQL drivers: %d\n", len(b))
		log.Println("no service bindings present falling back to ENV vars")
		dbinfo.dbUser = os.Getenv("DB_USER")
		dbinfo.dbPass = os.Getenv("DB_PASSWORD")
		dbinfo.dbName = os.Getenv("DB_NAME")
		dbinfo.dbHost = os.Getenv("DB_HOST")
	} else {
		log.Println("service bindings present")
		dbinfo.dbName, ok = bindings.Get(b[0], "database")
		if !ok {
			log.Fatal("No database name in binding")
		}
		dbinfo.dbHost, ok = bindings.Get(b[0], "host")
		if !ok {
			log.Fatal("No host in binding")
		}
		dbinfo.dbUser, ok = bindings.Get(b[0], "username")
		if !ok {
			log.Fatal("No username in binding")
		}
		dbinfo.dbPass, ok = bindings.Get(b[0], "password")
		if !ok {
			log.Fatal("No password in binding")
		}

	}

	return dbinfo
}

func SeedDB() {
	dbinfo := getBindings()
	dbname := dbinfo.dbName
	db := ConnectToDB(true)
	defer db.Close()

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	log.Print(reflect.TypeOf(db.Driver()).String())
	if reflect.TypeOf(db.Driver()).String() != "*sqlite.Driver" {
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
	}

	_, table_check := db.Query("select * from companies;")

	if table_check == nil {
		log.Printf("table exists not seeding")
	} else {
		_, err := db.ExecContext(ctx, "CREATE TABLE companies ( id integer, name varchar(32) )")
		if err != nil {
			log.Printf("Error %s when creating Table\n", err)
			return
		}

		query := "INSERT INTO `companies` (`id`, `name`) VALUES ( ?, ?)"
		_, err = db.ExecContext(ctx, query, 1, "acme")
		if err != nil {
			log.Printf("Error %s when seeding Table\n", err)
			return
		}
	}

	db.Close()
	db = ConnectToDB(true)
	defer db.Close()

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err := db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB", err)
		return
	}
	log.Printf("Connected to DB %s successfully\n", dbname)

}
