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
		items     text[]        not null
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

func GetUser(id string) (user User, err error) {
	var temp any
	err = db.QueryRow(`SELECT * FROM users WHERE id=$1`, id).Scan(&user.ID, &user.SID, &user.Name, &user.PasswordHash, &user.Score, &user.Bio, &user.Joined, &user.IsAdmin, &user.IsBanned, &temp)
	return user, err
}

func GetUsers(limit int, offset int) (users []User, err error) {
	rows, err := db.Query(`SELECT * FROM users LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return users, err
	}
	for rows.Next() {
		var user User
		var temp any
		if err := rows.Scan(&user.ID, &user.SID, &user.Name, &user.PasswordHash, &user.Score, &user.Bio, &user.Joined, &user.IsAdmin, &user.IsBanned, &temp); err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

func AddUser(user User) error {
	_, err := db.Exec(
		`INSERT INTO users(id, sid, name, pw_hash, score, bio, joined, is_admin, is_banned, items) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`,
		user.ID,
		user.SID,
		user.Name,
		user.PasswordHash,
		0,
		"",
		user.Joined,
		false,
		false,
		user.Items,
	)
	return err
}

func SetUserBio(id, bio string) error {
	_, err := db.Exec(`UPDATE users SET bio = $1 WHERE id = $2;`, bio, id)
	return err
}

func SetUserPasswordHash(id, passwordHash string) error {
	_, err := db.Exec(`UPDATE users SET pw_hash = $1 WHERE id = $2;`, passwordHash, id)
	return err
}

func SetUserName(id, name string) error {
	_, err := db.Exec(`UPDATE users SET name = $1 WHERE id = $2`, name, id)
	return err
}

func GetUserBySID(sid string) (user User, err error) {
	row := db.QueryRow(`SELECT * FROM users WHERE sid = $1`, sid)

	err = row.Scan(&user.ID, &user.SID, &user.Name, &user.PasswordHash, &user.Score, &user.Bio, &user.Joined, &user.IsAdmin, &user.IsBanned, &user.Items)
	return user, err
}

func SetUserSID(id, sid string) error {
	_, err := db.Exec(`UPDATE users SET sid = $1 WHERE id = $2`, sid, id)
	return err
}
