package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type UsersData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Link struct {
	Link  string `json:"link"`
	Email string `json:"email"`
}

type Users struct {
	Email string `json:"email"`
}

type LinkID struct {
	LinkID int64  `json:"link"`
	Email  string `json:"email"`
}

type LinkData struct {
	Id   int64
	Link string
}

// Add a link
func (s *server) addLink(w http.ResponseWriter, r *http.Request) {
	var l Link
	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user_id int64

	err = s.db.QueryRow(context.Background(), "select id from users where email like $1", l.Email).Scan(&user_id)

	tag, err := s.db.Exec(context.Background(), "insert into links (user_id, link)values($1, $2)", user_id, l.Link)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !tag.Insert() {
		http.Error(w, "Couldn't insert", http.StatusBadRequest)
	}

	w.Write([]byte("Successfully written"))

}

// Delete link
func (s *server) deleteLink(w http.ResponseWriter, r *http.Request) {

	var l LinkID
	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user_id int64

	err = s.db.QueryRow(context.Background(), "select id from users where email like $1", l.Email).Scan(&user_id)

	tag, err := s.db.Exec(context.Background(), "delete from links where user_id = $1 and id = $2", user_id, l.LinkID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !tag.Delete() {
		http.Error(w, "Couldn't delete", http.StatusBadRequest)
	}

	w.Write([]byte("Successfully deleted"))

}

// Get all links of a particular user
func (s *server) getLinks(w http.ResponseWriter, r *http.Request) {

	var u Users
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user_id int64

	err = s.db.QueryRow(context.Background(), "select id from users where email like $1", u.Email).Scan(&user_id)
	if err != nil {
		// fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var linksData []LinkData

	rows, err := s.db.Query(context.Background(), "select id, link from links where user_id = $1", user_id)
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

// Register account
func (s *server) register(w http.ResponseWriter, r *http.Request) {

	var u UsersData
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

	tag, err := s.db.Exec(context.Background(), "insert into users (email, password) values ($1, $2)", u.Email, pswd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !tag.Insert() {
		http.Error(w, "Couldn't insert", http.StatusBadRequest)
	}

	w.Write([]byte("Successfully written"))
}

// Login
func (s *server) login(w http.ResponseWriter, r *http.Request) {

	var u UsersData

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var hashedPassword string

	err = s.db.QueryRow(context.Background(), "select password from users where email like $1", u.Email).Scan(&hashedPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// w.Write([]byte(hashedPassword))
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(u.Password))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// Return encrypted username
	w.Write([]byte(encrypt(u.Email)))
}

// Convert to base64
func encrypt(s string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(s))
	return encoded
}
