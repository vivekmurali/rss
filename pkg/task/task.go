package task

import (
	"context"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mmcdole/gofeed"
)

// Overlay to run from main function through passing the database
func RunTask(dbpool *pgxpool.Pool) {
	fp := gofeed.NewParser()
	var task = func() {
		rows, err := dbpool.Query(context.Background(), "select link, user_id from links")
		if err != nil {
			fmt.Println("NOT WORKING IDK WHY ", err.Error())
		}

		for rows.Next() {
			var link string
			var user_id int64
			var email string

			err = rows.Scan(&link, &user_id)
			if err != nil {
			}

			err = dbpool.QueryRow(context.Background(), "select email from users where id = $1", user_id).Scan(&email)
			if err != nil {
				fmt.Println("NOT WORKING BECAUSE ", err.Error())
			}
			go run(fp, link, email)
		}

		if rows.Err() != nil {
			fmt.Println("ERROR WITH ", rows.Err())
		}
	}

	s := gocron.NewScheduler(time.UTC)

	// Change time
	s.Every(1).Day().At("16:07:20").Do(task)

	s.StartAsync()
}

// Function to find link's thing

func run(fp *gofeed.Parser, link, email string) {
	fmt.Println(link, email)
	feed, err := fp.ParseURL(link)
	if err != nil {
		fmt.Println("NOT WORKING BECAUSE ", err.Error())
	}
	// fmt.Println(feed.Items)
	for i, v := range feed.Items {
		t := *v.PublishedParsed
		if !time.Now().Add(-time.Hour * 24).After(t) {
			fmt.Println(i, v.PublishedParsed)
		}
	}
}
