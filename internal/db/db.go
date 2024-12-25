package db

import (
	"database/sql"
	"log"

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
		name      varchar(32)   not null,
		pw_hash   text          not null,
		score     integer       not null,
		bio       varchar(256),
		joined    timestamp     not null,
		is_admin  boolean       not null,
		is_banned boolean       not null,
		items     text[]
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
		children text[]
	)`)
	return err
}

func Close() error {
	return db.Close()
}

func GetUser(id string) (user User, err error) {
	err = db.QueryRow(`SELECT * FROM users WHERE id=$1`, id).Scan(&user)
	return user, err
}

func AddUser(user User) error {
	_, err := db.Exec(
		`INSERT INTO users(id, name, pw_hash, score, bio, joined, is_admin, is_banned) VALUES($1, $2, $3, $4, $5, $6, $7, $8);`,
		user.ID,
		user.Name,
		user.PasswordHash,
		0,
		"",
		user.Joined,
		false,
		false,
	)
	if err != nil {
		log.Println(err)
	}
	return err
}
