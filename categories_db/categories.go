package dbgeter

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Category struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ParentID int    `json:"parent_id"`
}

type DB struct {
	*sql.DB
}

func NewDB() (*DB, error) {
	db, err := sql.Open("postgres", "user=*** password=*** dbname=*** sslmode=disable")
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db *DB) GetCategories(id int) ([]Category, error) {
	query := `
	SELECT id, name, parent_id
	FROM categories
	WHERE path LIKE '%'|| $1 || '%'
	ORDER BY path;
	`

	rows, err := db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("error query: %w", err)
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.Name, &category.ParentID); err != nil {
			return nil, fmt.Errorf("error scan row: %w", err)
		}
		categories = append(categories, category)
	}

	return categories[1:], nil
}
