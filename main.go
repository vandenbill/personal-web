package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	"github.com/vandenbill/personal-web/internal"
)

func main() {
	r := chi.NewRouter()
	db, err := sql.Open("sqlite3", "./web.db")
	if err != nil {
		panic(err)
	}

	if err := internal.Migrate(db); err != nil {
		panic(err)
	}

	internal.NewHandler(internal.NewRepo(db), r)

	if err := http.ListenAndServe(":3000", r); err != nil {
		fmt.Println("Error: ", err)
	}
}
