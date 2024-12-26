package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/miresa-dev/miresa-srv/internal/db"
	"github.com/miresa-dev/miresa-srv/internal/verifier"

	"github.com/go-chi/chi/v5"

	"github.com/Kaamkiya/nanoid-go"

	"golang.org/x/crypto/bcrypt"
)

// CaptchaAndSID is hosted on /init. It contains a text captcha for a user to
// complete and a session ID associated with it.
func CaptchaAndSID(w http.ResponseWriter, r *http.Request) {
	sid, captcha, err := verifier.GenCaptchaSIDPair()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	pair := fmt.Sprintf(`{"sid":"%s","captcha":"%s"}`, sid, captcha)

	if _, err := w.Write([]byte(pair)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Login serves an endpoint where users can authenticate themselves.
func Login(w http.ResponseWriter, r *http.Request) {
	var input map[string]string

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if input["id"] == "" || input["password"] == "" || input["sid"] == "" || input["captcha"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := db.GetUser(input["id"])
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input["password"])); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !verifier.VerifyPair(input["sid"], input["captcha"]) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := db.SetUserSID(input["id"], input["sid"]); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Logout serves an endpoint that allows users to remove their session ID,
// effectively ending their sessions.
func Logout(w http.ResponseWriter, r *http.Request) {
	if len(r.Header["Authorization"]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, err := db.GetUserBySID(r.Header["Authorization"][0])
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to get user: %v\n", err)
		return
	}

	if u.SID != r.Header["Authorization"][0] {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := db.SetUserSID(u.ID, ""); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to clear user SID: %v\n", err)
	}
}

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

	if user["password"] == "" || user["username"] == "" || user["sid"] == "" || user["captcha"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !verifier.VerifyPair(user["sid"], user["captcha"]) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var u db.User

	u.ID = nanoid.Nanoid(64, nanoid.DefaultCharset)
	u.Name = user["username"]
	u.Joined = time.Now()
	u.Bio = ""
	u.SID = user["sid"]

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
		return
	}

	if err := json.NewEncoder(w).Encode(u); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to marshal JSON: %v\n", err)
	}
}

// UpdateUser rewrites fields of a user from an HTTP patch request.
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	u, err := db.GetUser(id)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(r.Header["Authorization"]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if u.SID != "" && u.SID != r.Header["Authorization"][0] {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var user struct {
		db.User
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		// If we can't unmarshal request data to a user, it must be an
		// invalid request.
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if user.Name != "" {
		// Display names can be updated.
		if err := db.SetUserName(id, user.Name); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Failed to update user: %v\n", err)
			return
		}
	}
	if user.Bio != "" {
		if err := db.SetUserBio(id, user.Bio); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Failed to update user: %v\n", err)
			return
		}
	}

	if user.Password != "" {
		data, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Failed to hash password")
			return
		}
		if err = db.SetUserPasswordHash(id, string(data)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Failed to update user: %v\n", err)
			return
		}
	}
}
