package api

import (
	"net/http"
	"strconv"
	"errors"
	"time"
	"log"
	"encoding/json"
	"database/sql"

	"github.com/miresa-dev/miresa-srv/internal/db"

	"github.com/go-chi/chi/v5"

	"github.com/Kaamkiya/nanoid-go"
)

func GetItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	item, err := db.GetItem(id)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(item); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func GetItems(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	limit, err := strconv.Atoi(queryParams.Get("limit"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	offset, err := strconv.Atoi(queryParams.Get("offset"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	creator := queryParams.Get("creator")
	if creator != "" {
		items, err := db.GetItemsByCreator(creator, limit, offset)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err = json.NewEncoder(w).Encode(items); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Failed to marshal JSON: %v\n", err)
		}
		return
	}

	items, err := db.GetItems(limit, offset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to get items from DB: %v\n", err)
		return
	}

	if err = json.NewEncoder(w).Encode(items); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to marshal JSON: %v\n", err)
	}
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
	var item db.Item

	if len(r.Header["Authorization"]) < 1 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, err := db.GetUserBySID(r.Header["Authorization"][0])
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to get user from DB: %v\n", err)
		return
	}

	if item.Title != "" && item.Parent != "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	item.ID = nanoid.Nanoid(64, nanoid.DefaultCharset)
	item.Creator = u.ID
	item.Points = 0
	item.Children = []string{}
	item.Published = time.Now()

	if err := db.AddItem(item); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to add item to DB: %v\n", err)
		return
	}

	if err := json.NewEncoder(w).Encode(item); err != nil {
		// Don't return the error to the user. If everything succeeded,
		// and they get a 500, they will try again, which is bad
		// because everything worked. So log the error and move on.
		log.Printf("Failed to marshal JSON: %v\n", err)
	}
}

func UpdateItem(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todo"))
}
