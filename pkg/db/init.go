package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

var PrivateAPIKey string

func InitDB() *pgxpool.Pool {
	PrivateAPIKey = os.Getenv("MG_API_KEY")
	dbpool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Println(err)
	}
	return dbpool
}
