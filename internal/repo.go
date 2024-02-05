package internal

import "database/sql"

type Repo struct {
	DB *sql.DB
}

func (r Repo) SaveDraft(title string, draft string) error {
	_, err := r.DB.Exec(`insert into article_drafts(id, title, draft) values(1, ?, ?) on conflict(id) do update set title=excluded.title, draft=excluded.draft`, title, draft)
	return err
}

func (r Repo) GetDraft() (Draft, error) {
	d := Draft{}
	row := r.DB.QueryRow(`select title, draft from article_drafts where id = 1`)
	err := row.Scan(&d.Title, &d.Draft)
	return d, err
}
