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

	err = internal.Migrate(db)
	if err != nil {
		panic(err)
	}

	repo := internal.Repo{DB : db}
	h := internal.Handler{R : repo}

	r.Get("/", h.RootHandler)
	r.Get("/{name}", h.StaticFile)
	r.Get("/articles", h.CreateArticleView)

	r.Post("/articles/draft", h.CreateDraftArticle)
	r.Post("/articles", h.CreateArticle)

	if err := http.ListenAndServe(":3000", r); err != nil {
		fmt.Println("Error: ", err)
	}
}
