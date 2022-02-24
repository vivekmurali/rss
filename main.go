package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/mailgun/mailgun-go/v4"
)

type server struct {
	db    *pgxpool.Pool
	mg    *mailgun.MailgunImpl
	store *sessions.CookieStore
}

func main() {
	godotenv.Load()

	// Initialization
	s := &server{mg: mailgun.NewMailgun("vivekmurali.in", os.Getenv("MG_API_KEY"))}
	s.store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Database implementation
	s.db = initDB()
	defer s.db.Close()

	// Job Scheduler
	go s.runTask()

	// Routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("HELLO WORLD"))
	})
	r.Post("/register", s.register)
	r.Post("/login", s.login)
	r.Post("/logout", s.logout)
	r.Post("/add", s.addLink)
	r.Delete("/delete", s.deleteLink)
	r.Get("/get", s.getLinks)

	http.ListenAndServe(":3000", r)
}
