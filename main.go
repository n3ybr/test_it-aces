package main

import (
	"encoding/json"
	"fmt"
	categories_db "gogogo/categories_db"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func categoriesHandler(cdb *categories_db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		categoryID, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "err strconv ID", http.StatusBadRequest)
			return
		}
		categories, err := cdb.GetCategories(categoryID)
		if err != nil {
			http.Error(w, fmt.Sprintf("err get categories: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if len(categories) == 0 {
			if err := json.NewEncoder(w).Encode("No Categories Found"); err != nil {
				http.Error(w, "err messjson response", http.StatusInternalServerError)
			}
		} else {
			if err := json.NewEncoder(w).Encode(categories); err != nil {
				http.Error(w, "err cat json response", http.StatusInternalServerError)
			}
		}
	}
}

func main() {
	db, err := categories_db.NewDB()
	if err != nil {
		log.Fatalf("err conn to database: %v", err)
	}
	defer db.Close()
	r := mux.NewRouter()

	r.HandleFunc("/categories/{id:[0-9]+}", categoriesHandler(db))

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("err serve: %v", err)
	}

	fmt.Println("Server start port 8080")

}
