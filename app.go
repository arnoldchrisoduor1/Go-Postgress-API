// a simple application to test the connections to the databases.

package main

import (
	"databse/sql"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialze(postgres, password, go_api string) { }

func (a *App) Run(addr string) { }