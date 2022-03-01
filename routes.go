package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

type UsersData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LinkData struct {
	Id   int64
	Link string
}

// Add a link
func (s *server) addLink(w http.ResponseWriter, r *http.Request) {
	session, _ := s.store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	link := r.FormValue("link")
	// fmt.Println(link)

	// _, err := url.Parse(link)
	// if err != nil {
	// 	http.Error(w, "not a valid link", http.StatusBadRequest)
	// 	return
	// }
	// host := u.Host

	// tries := []string{"feed", "index.xml", "rss"}

	_, err := s.parser.ParseURL(link)
	if err != nil {
		resp, err := http.Get(link)
		if err != nil {
			http.Error(w, "not a valid link", http.StatusBadRequest)
			return
		}
		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			fmt.Println(err.Error())
		}

		doc.Find("link").Each(func(i int, s *goquery.Selection) {
			t, ok := s.Attr("type")
			if !ok {
				return
			}
			if t == "application/rss+xml" || t == "application/atom+xml" {
				// fmt.Println(s.Attr("href"))
				link, _ = s.Attr("href")
			}
		})
		// }

		// resp, err := http.Get(link)
		// if err != nil {
		// 	http.Error(w, "not a valid link", http.StatusBadRequest)
		// 	return
		// }
		// defer resp.Body.Close()

		// doc, err := goquery.NewDocumentFromReader(resp.Body)
		// if err != nil {
		// 	fmt.Println(err.Error())
		// }

		// doc.Find("link").Each(func(i int, s *goquery.Selection) {
		// 	t, ok := s.Attr("type")
		// 	if !ok {
		// 		return
		// 	}
		// 	if t == "application/rss+xml" || t == "application/atom+xml" {
		// 		// fmt.Println(s.Attr("href"))
		// 		link, _ = s.Attr("href")
		// 	}
		// })

		var user_id int64

		err = s.db.QueryRow(context.Background(), "select id from users where email like $1", session.Values["email"]).Scan(&user_id)

		tag, err := s.db.Exec(context.Background(), "insert into links (user_id, link)values($1, $2)", user_id, link)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !tag.Insert() {
			http.Error(w, "Couldn't insert", http.StatusBadRequest)
		}

		// w.Write([]byte("Successfully written"))
		w.WriteHeader(http.StatusOK)
	}
}

// Delete link
func (s *server) deleteLink(w http.ResponseWriter, r *http.Request) {
	session, _ := s.store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	linkID := chi.URLParam(r, "id")

	var user_id int64

	err := s.db.QueryRow(context.Background(), "select id from users where email like $1", session.Values["email"]).Scan(&user_id)

	tag, err := s.db.Exec(context.Background(), "delete from links where user_id = $1 and id = $2", user_id, linkID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !tag.Delete() {
		http.Error(w, "Couldn't delete", http.StatusBadRequest)
	}

	w.Write([]byte("Successfully deleted"))
	// http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
}

// Get all links of a particular user DASHBOARD
func (s *server) getLinks(w http.ResponseWriter, r *http.Request) {

	session, _ := s.store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		// http.Error(w, "Forbidden", http.StatusForbidden)
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	var user_id int64

	err := s.db.QueryRow(context.Background(), "select id from users where email like $1", session.Values["email"]).Scan(&user_id)
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

	tmpl, _ := template.ParseFiles("template/dashboard.html")
	tmpl.Execute(w, linksData)
}

// Register account
func (s *server) register(w http.ResponseWriter, r *http.Request) {

	var u UsersData
	fmt.Println("USER: ", u.Email)
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

	session, _ := s.store.Get(r, "cookie-name")

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
	session.Values["authenticated"] = true
	session.Values["email"] = u.Email
	session.Save(r, w)
	w.Write([]byte("Login successful"))
}

func (s *server) logout(w http.ResponseWriter, r *http.Request) {
	session, _ := s.store.Get(r, "cookie-name")
	session.Values["authenticated"] = false
	session.Save(r, w)
	w.Write([]byte("Logged out"))
}
