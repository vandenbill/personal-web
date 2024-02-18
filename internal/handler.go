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
	Repo   *Repo
	Router *chi.Mux
}

func NewHandler(Repo *Repo, Router *chi.Mux) *Handler {
	h := Handler{Repo, Router}

	Router.Get("/", h.RootView)
	Router.Get("/public/{name}", h.StaticFile)
	Router.Get("/articles/{article_id}", h.GetArticle)

	Router.Get("/admin", h.AdminView)
	Router.Get("/admin/articles/{article_id}", h.CreateArticleView)
	Router.Post("/admin/articles/{article_id}/save", h.SaveArticle)
	Router.Post("/admin/tags/{article_id}/{tag_id}", h.AddTag)

	return &h
}

func (h Handler) AddTag(w http.ResponseWriter, r *http.Request) {
	articleID := chi.URLParam(r, "article_id")
	tagID := chi.URLParam(r, "tag_id")

	err := h.Repo.AddTag(articleID, tagID)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tags, err := h.Repo.GetTagsByArticleID(articleID)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("./public/create.html")
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(w, "list-tags", map[string]any{"usedTags":tags})
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h Handler) RootView(w http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.ParseFiles("./public/index.html")
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tags, err := h.Repo.GetAllTags()
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	articles, err := h.Repo.GetArticles()
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	halfLen := len(articles) / 2
	leftCol := make([]Article, 0, 100)
	rightCol := make([]Article, 0, 100)

	for i := 0; i < len(articles); i++ {
		if i > halfLen-1 {
			rightCol = append(rightCol, articles[i])
			continue
		}
		leftCol = append(leftCol, articles[i])
	}

	data := map[string]any{
		"leftArticles":  leftCol,
		"rightArticles": rightCol,
		"tags":          tags,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h Handler) AdminView(w http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.ParseFiles("./public/admin.html")
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	articles, err := h.Repo.GetArticles()
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	leftCol := make([]Article, 0, 100)
	rightCol := make([]Article, 0, 100)

	halfLen := len(articles) / 2
	if halfLen%2 != 0 {
		halfLen++
	}

	for i := 0; i < len(articles); i++ {
		if i%2 == 0 {
			leftCol = append(leftCol, articles[i])
		} else {
			rightCol = append(rightCol, articles[i])
		}
	}

	data := map[string]any{
		"leftArticles":  leftCol,
		"rightArticles": rightCol,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h Handler) StaticFile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, fmt.Sprintf("./public/%s", chi.URLParam(r, "name")))
}

func (h Handler) GetArticle(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./public/article.html")
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	articleID := chi.URLParam(r, "article_id")

	article, err := h.Repo.GetArticleByID(articleID)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	html := pkg.MdToHTML([]byte(article.Markdown))

	data := map[string]any{
		"title":   article.Title,
		"article": html,
		"tags":    article.Tags,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h Handler) CreateArticleView(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./public/create.html")
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	articleID := chi.URLParam(r, "article_id")

	allTags, err := h.Repo.GetAllTags()
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	type tempTag struct {
		Tag
		ArticleID string
	}

	tempTags := make([]tempTag, 0, 100)
	for i, _ := range allTags {
		tempTags = append(tempTags, tempTag{Tag: allTags[i], ArticleID: articleID})
	}

	data := map[string]any{
		"title":       "",
		"markdown":    "",
		"article":     "",
		"description": "",
		"articleID":   "new",
		"usedTags":    []Tag{},
		"allTags":     tempTags,
	}

	if articleID != "new" {
		article, err := h.Repo.GetArticleByID(articleID)
		if err != nil {
			slog.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		html := pkg.MdToHTML([]byte(fmt.Sprintf("<h1>%s</h1>\n%s", article.Title, article.Markdown)))
		data = map[string]any{
			"title":       article.Title,
			"markdown":    article.Markdown,
			"article":     html,
			"description": article.Description,
			"usedTags":    article.Tags,
			"allTags":     tempTags,
			"articleID":   article.ID,
		}
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h Handler) SaveArticle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "article_id")
	md := r.FormValue("markdown")
	title := r.FormValue("title")
	desc := r.FormValue("description")

	if id == "new" {
		id = ""
	}

	err := h.Repo.SaveArticle(id, title, string(md), desc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if id == "" {
		w.Header().Set("HX-Redirect", "/admin")
	}
	w.WriteHeader(http.StatusOK)
}
