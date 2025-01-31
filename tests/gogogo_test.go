package gogogo_test

import (
	"encoding/json"
	categories_db "gogogo/categories_db"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestInitCategoriesDB(t *testing.T) {
	_, err := categories_db.NewDB()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGetCategories(t *testing.T) {
	db, err := categories_db.NewDB()
	if err != nil {
		t.Fatalf("err connect to database: %v", err)
	}
	defer db.Close()
	id := 2

	categories, err := db.GetCategories(id)
	if err != nil {
		t.Fatalf("got err %v", err)
	}

	expected := []categories_db.Category{
		{ID: 4, Name: "Смартфоны", ParentID: 2},
		{ID: 5, Name: "Аксессуары", ParentID: 2},
		{ID: 6, Name: "Чехлы", ParentID: 5},
		{ID: 7, Name: "Зарядки", ParentID: 5},
	}

	assert.Equal(t, expected, categories, "should match")
}

func TestCategoriesHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/categories/2", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r := mux.NewRouter()

	r.HandleFunc("/categories/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		categoryID, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "err categoryID", http.StatusBadRequest)
			return
		}

		db, err := categories_db.NewDB()
		if err != nil {
			http.Error(w, "err conn db", http.StatusInternalServerError)
			return
		}

		categories, err := db.GetCategories(categoryID)
		if err != nil {
			http.Error(w, "err GetCategories", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(categories)
	}).Methods("GET")

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "status code != 200")

	expectedResponse := []categories_db.Category{
		{ID: 4, Name: "Смартфоны", ParentID: 2},
		{ID: 5, Name: "Аксессуары", ParentID: 2},
		{ID: 6, Name: "Чехлы", ParentID: 5},
		{ID: 7, Name: "Зарядки", ParentID: 5},
	}

	var actualResponse []categories_db.Category
	err = json.NewDecoder(rr.Body).Decode(&actualResponse)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedResponse, actualResponse, "should be equal")
}
