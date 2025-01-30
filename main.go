package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

type Category struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ParentID *int   `json:"parent_id"`
}

func getCategories(db *sql.DB, id int) ([]Category, error) {
	//я загуглил
	query := `
	WITH RECURSIVE category_tree AS (
		SELECT id, name, parent_id FROM categories WHERE parent_id = $1
		UNION ALL
		SELECT c.id, c.name, c.parent_id FROM categories c
		INNER JOIN category_tree ct ON c.parent_id = ct.id
	)
	SELECT * FROM category_tree;
	`
	rows, err := db.Query(query, id)
	if err != nil {
		fmt.Printf("error query %v", err)
		return nil, err

	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.ID, &category.Name, &category.ParentID)
		if err != nil {
			fmt.Printf("error scan %v", err)
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func main() {
	db, err := sql.Open("postgres", "user=*** password=*** dbname=*** sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		queryValues := r.URL.Query()
		categoryIDStr := queryValues.Get("id")
		categoryID, err := strconv.Atoi(categoryIDStr)
		if err != nil {
			fmt.Printf("err conv categoryID %v", err)
			return
		}
		categories, err := getCategories(db, categoryID)
		if err != nil {
			fmt.Printf("err get categories %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(categories)
	})

	http.ListenAndServe(":8080", nil)
	fmt.Println("server start 8080")

}
