package db

func GetItem(id string) (item Item, err error) {
	row := db.QueryRow(`SELECT * FROM items WHERE id = $1`, id)
	err = row.Scan(
		&item.ID,
		&item.Creator,
		&item.Points,
		&item.Title,
		&item.Content,
		&item.Parent,
		&item.Children,
		&item.Published,
	)
	return item, err
}

func AddItem(item Item) error {
	_, err := db.Exec(
		`INSERT INTO items(id, title, parent, content, creator, points, created, children) VALUES($1, $2, $3, $4, $5, $6, $7, $8);`,
		item.ID,
		item.Title,
		item.Parent,
		item.Content,
		item.Creator,
		item.Points,
		item.Published,
		item.Children,
	)
	return err
}
