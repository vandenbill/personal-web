package internal

import (
	"database/sql"
)

func Migrate(db *sql.DB) error {
	_, err := db.Exec(`create table if not exists articles (
		id integer primary key,
		title text,
		markdown text,
		description text,
		status text default 'draft',
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

	_, err = db.Exec(`create table if not exists articles_tags (
		id integer primary key,
		article_id integer references articles(id) on delete cascade on update cascade,
		tag_id integer references tags(id) on delete cascade on update cascade
	);`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`create table if not exists tags (
		id integer primary key,
		name text
	);`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`create index if not exists name_index on tags(name)`)
	if err != nil {
		return err
	}

	return nil
}
