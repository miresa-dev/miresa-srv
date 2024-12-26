package db

import "github.com/lib/pq"

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

func GetItems(limit, offset int) ([]Item, error) {
	items := make([]Item, limit)

	rows, err := db.Query(`SELECT * FROM items LIMIT $1 OFFSET $2;`, limit, offset)
	if err != nil {
		return items, err
	}
	i := 0
	for rows.Next() {
		var it Item
		err = rows.Scan(&it.ID, &it.Title, &it.Parent, &it.Content, &it.Creator, &it.Points, &it.Published, (*pq.StringArray)(&it.Children))
		if err != nil {
			return items, err
		}
		items[i] = it
		i++
	}
	return items, nil
}

func GetItemsByCreator(id string, limit, offset int) ([]Item, error) {
	items := make([]Item, limit)

	rows, err := db.Query(`SELECT * FROM items LIMIT $1 OFFSET $2 WHERE creator = $3`, limit, offset, id)
	if err != nil {
		return items, err
	}

	i := 0
	for rows.Next() {
		var it Item
		err = rows.Scan(&it.ID, &it.Title, &it.Parent, &it.Content, &it.Creator, &it.Points, &it.Published, (*pq.StringArray)(&it.Children))
		if err != nil {
			return items, err
		}
		items[i] = it
		i++
	}
	return items, nil
}
