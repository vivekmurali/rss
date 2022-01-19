package main

import (
	"net/http"
	"rss/pkg/auth"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	r.Post("/register", auth.RegisterHandler)
	r.Post("/login", auth.LoginHandler)
	http.ListenAndServe(":3000", r)
}
