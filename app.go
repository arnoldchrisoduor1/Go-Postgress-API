// a simple application to test the connections to the databases.

package main

import (
	"database/sql"
	"log"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize() { 
	connectionString :=
		"host=localhost port=5432 user=postgres password=12345 dbname=go_api sslmode=disable"

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	a.Router = mux.NewRouter()
 }

func (a *App) Run(addr string) { }