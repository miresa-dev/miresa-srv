package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/miresa-dev/miresa-srv/internal/db"

	"github.com/go-chi/chi/v5"

	"github.com/Kaamkiya/nanoid-go"

	"golang.org/x/crypto/bcrypt"
)

// GetUser gets a single user by ID from the database and returns them as JSON.
func GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user, err := db.GetUser(id)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to get user: %v\n", err)
		return
	}
	if err = json.NewEncoder(w).Encode(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to marshal JSON: %v\n", err)
	}
}

// GetAllUsers gets all the users from the database, with a limit and an offset
// specified in the query parameters.
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	limit, err := strconv.Atoi(queryParams.Get("limit"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	offset, err := strconv.Atoi(queryParams.Get("offset"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users, err := db.GetUsers(limit, offset)
	if err != nil {
		// There must be at least one user on the instance, so if we
		// can't get any, something must be wrong.
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to get users: %v\n", err)
		return
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		// If we can't marshal a list of structs, something is probably
		// wrong on our end.
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to marshal JSON: %v\n", err)
		return
	}
}

// CreateUser creates a new user and adds them to the database.
func CreateUser(w http.ResponseWriter, r *http.Request) {
	user := map[string]string{}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: actually check the SID and CAPTCHA.
	if user["password"] == "" || user["username"] == "" || user["sid"] == "" || user["captcha"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var u db.User

	u.ID = nanoid.Nanoid(64, nanoid.DefaultCharset)
	u.Name = user["username"]
	u.Joined = time.Now()
	u.Bio = ""
	u.Items = []string{}

	data, err := bcrypt.GenerateFromPassword([]byte(user["password"]), 10)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to hash password: %v\n", err)
		return
	}

	u.PasswordHash = string(data)
	if err := db.AddUser(u); err != nil {
		log.Printf("Failed to add user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	if err := json.NewEncoder(w).Encode(u); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to marshal JSON: %v\n", err)
	}
}
