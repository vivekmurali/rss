package routes

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Link struct {
	Link  string `json:"link"`
	Email string `json:"email"`
}

func AddLink(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var l Link
		err := json.NewDecoder(r.Body).Decode(&l)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var user_id int64

		err = dbpool.QueryRow(context.Background(), "select id from users where email like $1", l.Email).Scan(&user_id)

		tag, err := dbpool.Exec(context.Background(), "insert into links (user_id, link)values($1, $2)", user_id, l.Link)
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
