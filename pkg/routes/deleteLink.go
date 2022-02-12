package routes

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
)

type LinkID struct {
	LinkID int64  `json:"link"`
	Email  string `json:"email"`
}

func DeleteLink(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var l LinkID
		err := json.NewDecoder(r.Body).Decode(&l)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var user_id int64

		err = dbpool.QueryRow(context.Background(), "select id from users where email like $1", l.Email).Scan(&user_id)

		tag, err := dbpool.Exec(context.Background(), "delete from links where id = $1 and user_id = $2", l.LinkID, user_id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !tag.Delete() {
			http.Error(w, "Couldn't delete", http.StatusBadRequest)
		}

		w.Write([]byte("Successfully deleted"))

	}
}
