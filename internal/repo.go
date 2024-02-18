package internal

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/vandenbill/personal-web/pkg"
)

type Repo struct {
	DB *sql.DB
}

func NewRepo(DB *sql.DB) *Repo {
	return &Repo{DB}
}

func (r Repo) GetTagsByArticleID(articleID string) ([]Tag, error) {
	rows, err := r.DB.Query(`select tag_id from articles_tags where article_id = ?`, articleID)
	if err != nil {
		return nil, err
	}
	tagsID := make([]string, 0, 100)
	for rows.Next() {
		var id string
		rows.Scan(&id)
		tagsID = append(tagsID, id)
	}

	tags, err := r.GetTagsByIds(tagsID...)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (r Repo) AddTag(articleID, tagID string) error {
	_, err := r.DB.Exec(`insert into articles_tags(article_id, tag_id) values(?, ?)`,
		articleID, tagID)

	return err
}

func (r Repo) SaveArticle(id, title, md, desc string, tags ...string) error {
	var err error
	if id == "" {
		_, err = r.DB.Exec(`insert into articles(title, markdown, description) values(?, ?, ?)`,
			title, md, desc)
	} else {
		_, err = r.DB.Exec(`insert into articles(id, title, markdown, description) values(?, ?, ?, ?) on conflict(id) do update set title=excluded.title, markdown=excluded.markdown, description=excluded.description`,
			id, title, md, desc)
	}

	return err
}

func (r Repo) GetArticles() ([]Article, error) {
	res := make([]Article, 0, 100)
	rows, err := r.DB.Query(`select id, title, markdown, description from articles`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		a := Article{}
		err := rows.Scan(&a.ID, &a.Title, &a.Markdown, &a.Description)
		if err != nil {
			return nil, err
		}
		res = append(res, a)
	}
	return res, err
}

func (r Repo) GetTagsByIds(ids ...string) ([]Tag, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	tags := make([]Tag, 0, 100)

	q := fmt.Sprintf("select name from tags where id in (%s)", strings.Join(ids, ","))
	rows, err := r.DB.Query(q)
	if err != nil {
		return nil, pkg.MaskErr(err)
	}

	for rows.Next() {
		t := Tag{}
		err = rows.Scan(&t.Name)
		fmt.Println(t.Name)
		if err != nil {
			return nil, pkg.MaskErr(err)
		}
		tags = append(tags, t)
	}
	return tags, nil
}

func (r Repo) GetArticleByID(id string) (Article, error) {
	row := r.DB.QueryRow(`select id, title, markdown, description, status from articles where id = ?`, id)

	article := Article{}
	err := row.Scan(&article.ID, &article.Title, &article.Markdown, &article.Description, &article.Status)
	if err != nil {
		return article, err
	}

	rows, err := r.DB.Query(`select tag_id from articles_tags where article_id = ?`, article.ID)
	if err != nil {
		return article, err
	}
	tagsID := make([]string, 0, 100)
	for rows.Next() {
		var id string
		rows.Scan(&id)
		fmt.Println(id)
		tagsID = append(tagsID, id)
	}

	tags, err := r.GetTagsByIds(tagsID...)
	if err != nil {
		return article, err
	}

	article.Tags = tags
	return article, nil
}

func (r Repo) GetAllTags() ([]Tag, error) {
	tags := make([]Tag, 0, 100)
	rows, err := r.DB.Query(`select id, name from tags`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		t := Tag{}
		err = rows.Scan(&t.ID, &t.Name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, err
}
