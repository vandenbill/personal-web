package internal

import (
	"database/sql"
)

func Migrate(db *sql.DB) error {
	_, err := db.Exec(`create table if not exists articles (
		id integer primary key,
		title text,
		article text,
		created_at datetime default CURRENT_TIMESTAMP,
		updated_at datetime default CURRENT_TIMESTAMP,
		deleted_at datetime
	);`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`create index if not exists title_index on articles(title);`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`create table if not exists article_tags (
		id integer primary key,
		article_id integer references articles(id) on delete cascade on update cascade,
		tags_id integer references tags(id) on delete cascade on update cascade
	);`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`create table if not exists tags (
		id integer primary key,
		name text,
		created_at datetime default CURRENT_TIMESTAMP,
		updated_at datetime default CURRENT_TIMESTAMP,
		deleted_at datetime null
	);`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`create index if not exists name_index on tags(name)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`create table if not exists article_drafts (
		id integer primary key,
		title text,
		draft text
	);`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`insert into article_drafts(id, title, draft) values(1, "title", "draft") on conflict(id) do nothing`)
	if err != nil {
		return err
	}

	return nil
}
