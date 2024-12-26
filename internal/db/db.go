package db

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB

func Init(url string) (err error) {
	db, err = sql.Open("pgx", url)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users(
		id        varchar(64)   primary key,
		sid       varchar(64)   not null,
		name      varchar(32)   not null,
		pw_hash   text          not null,
		score     integer       not null,
		bio       varchar(256)  not null,
		joined    timestamp     not null,
		is_admin  boolean       not null,
		is_banned boolean       not null,
	);`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS items(
		id       varchar(64)   primary key,
		title    varchar(128),
		parent   varchar(64),
		content  varchar(2048) not null,
		creator  varchar(64)   not null, 
		points   integer       not null,
		created  timestamp     not null,
		children text[]        not null
	)`)
	return err
}

func Close() error {
	return db.Close()
}
