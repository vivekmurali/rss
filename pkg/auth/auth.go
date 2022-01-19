package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RegisterHandler(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u Users
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		pswd, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tag, err := dbpool.Exec(context.Background(), "insert into users (email, password) values ($1, $2)", u.Email, pswd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !tag.Insert() {
			http.Error(w, "Couldn't insert", http.StatusBadRequest)
		}

		w.Write([]byte("Successfully written"))

	}
}

func LoginHandler(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
