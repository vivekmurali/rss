package routes

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Users struct {
	Email string `json:"email"`
}

type LinkData struct {
	Id   int64
	Link string
}

func GetLinks(dbpool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u Users
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var user_id int64

		err = dbpool.QueryRow(context.Background(), "select id from users where email like $1", u.Email).Scan(&user_id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		var linksData []LinkData

		rows, err := dbpool.Query(context.Background(), "select id, link from links where user_id = $1", user_id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		defer rows.Close()

		for rows.Next() {
			var id int64
			var link string

			err = rows.Scan(&id, &link)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			linksData = append(linksData, LinkData{id, link})
		}

		if rows.Err() != nil {
			http.Error(w, rows.Err().Error(), http.StatusBadRequest)
		}

		jsonData, err := json.Marshal(linksData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		w.Write(jsonData)

	}
}
