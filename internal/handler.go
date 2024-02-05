package internal

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vandenbill/personal-web/pkg"
)

type Handler struct {
	R Repo
}

func (h Handler) RootHandler(w http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.ParseFiles("./public/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(nil)
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}
}

func (h Handler) StaticFile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, fmt.Sprintf("./public/%s", chi.URLParam(r, "name")))
}

func (h Handler) CreateArticleView(w http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.ParseFiles("./public/create.html")
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	draft, err := h.R.GetDraft()
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	article := pkg.MdToHTML([]byte(fmt.Sprintf("<h1>%s</h1>\n%s", draft.Title, draft.Draft)))

	data := map[string]any{
		"title":   draft.Title,
		"draft":   draft.Draft,
		"article": string(article),
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h Handler) CreateDraftArticle(w http.ResponseWriter, r *http.Request) {
	md := r.FormValue("md")
	title := r.FormValue("title")

	err := h.R.SaveDraft(title, string(md))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	article := pkg.MdToHTML([]byte(md))

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("<h1>%s</h1>\n%s", title, article)))
}

func (h Handler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	md := r.FormValue("md")
	title := r.FormValue("title")

	html := pkg.MdToHTML([]byte(md))
	html = []byte(fmt.Sprintf("<h1>%s</h1>\n%s", title, html))

	w.WriteHeader(http.StatusOK)
	w.Write(html)
}
