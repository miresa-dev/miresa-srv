package db

func GetUser(id string) (user User, err error) {
	err = db.QueryRow(`SELECT * FROM users WHERE id=$1`, id).Scan(&user.ID, &user.SID, &user.Name, &user.PasswordHash, &user.Score, &user.Bio, &user.Joined, &user.IsAdmin, &user.IsBanned)
	return user, err
}

func GetUsers(limit int, offset int) (users []User, err error) {
	rows, err := db.Query(`SELECT * FROM users LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return users, err
	}
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.SID, &user.Name, &user.PasswordHash, &user.Score, &user.Bio, &user.Joined, &user.IsAdmin, &user.IsBanned); err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}

func AddUser(user User) error {
	_, err := db.Exec(
		`INSERT INTO users(id, sid, name, pw_hash, score, bio, joined, is_admin, is_banned) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9);`,
		user.ID,
		user.SID,
		user.Name,
		user.PasswordHash,
		0,
		"",
		user.Joined,
		false,
		false,
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

	err = row.Scan(&user.ID, &user.SID, &user.Name, &user.PasswordHash, &user.Score, &user.Bio, &user.Joined, &user.IsAdmin, &user.IsBanned)
	return user, err
}

func SetUserSID(id, sid string) error {
	_, err := db.Exec(`UPDATE users SET sid = $1 WHERE id = $2`, sid, id)
	return err
}
