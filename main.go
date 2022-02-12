package main

import (
	"net/http"
	"rss/pkg/auth"
	"rss/pkg/db"
	"rss/pkg/routes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	pool := db.InitDB()
	defer pool.Close()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	r.Post("/register", auth.RegisterHandler(pool))
	r.Post("/login", auth.LoginHandler(pool))
	r.Post("/add", routes.AddLink(pool))
	r.Delete("/delete", routes.DeleteLink(pool))
	r.Get("/get", routes.GetLinks(pool))
	http.ListenAndServe(":3000", r)
}
