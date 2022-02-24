package main

import (
	"net/http"
	"rss/pkg/auth"
	"rss/pkg/db"
	"rss/pkg/routes"
	"rss/pkg/task"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	// Initialization
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Database implementation
	pool := db.InitDB()
	defer pool.Close()

	// Job Scheduler
	go task.RunTask(pool)

	// Routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("HELLO WORLD"))
	})
	r.Post("/register", auth.RegisterHandler(pool))
	r.Post("/login", auth.LoginHandler(pool))
	r.Post("/add", routes.AddLink(pool))
	r.Delete("/delete", routes.DeleteLink(pool))
	r.Get("/get", routes.GetLinks(pool))

	http.ListenAndServe(":3000", r)
}
