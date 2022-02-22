package task

import (
	"context"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mmcdole/gofeed"
)

type ItemData struct {
	title string
	desc  string
}

// Overlay to run from main function through passing the database
func RunTask(dbpool *pgxpool.Pool) {
	var task = func() {
		findUsers(dbpool)
	}

	s := gocron.NewScheduler(time.UTC)

	// Change time
	s.Every(1).Day().At("20:07:20").Do(task)

	s.StartAsync()
}

// Function to find users

func findUsers(dbpool *pgxpool.Pool) {
	fp := gofeed.NewParser()
	// rows, err := dbpool.Query(context.Background(), "select link, user_id from links")
	// if err != nil {
	// 	fmt.Println("NOT WORKING IDK WHY ", err.Error())
	// }

	// for rows.Next() {
	// 	var link string
	// 	var user_id int64
	// 	var email string

	// 	err = rows.Scan(&link, &user_id)
	// 	if err != nil {
	// 	}

	// 	err = dbpool.QueryRow(context.Background(), "select email from users where id = $1", user_id).Scan(&email)
	// 	if err != nil {
	// 		fmt.Println("NOT WORKING BECAUSE ", err.Error())
	// 	}
	// 	go run(fp, link, email)
	// }

	// if rows.Err() != nil {
	// 	fmt.Println("ERROR WITH ", rows.Err().Error())
	// }

	rows, err := dbpool.Query(context.Background(), "select id, email from users")
	if err != nil {
		fmt.Println("NOT WORKING IDK WHY", err.Error())
	}

	for rows.Next() {
		var id int64
		var email string

		err = rows.Scan(&id, &email)
		if err != nil {
			fmt.Println("Something wrong with rows", err.Error())
		}
		if rows.Err() != nil {
			fmt.Println("ERROR WITH ", rows.Err().Error())
		}

		singleUser(fp, id, email, dbpool)
	}

}

// Function to run a single user's posts

func singleUser(fp *gofeed.Parser, id int64, email string, dbpool *pgxpool.Pool) {

	rows, err := dbpool.Query(context.Background(), "select link from links where user_id = $1", id)
	if err != nil {
		fmt.Println("Something wrong with rows", err.Error())
	}
	var links []string

	for rows.Next() {
		var link string
		err = rows.Scan(&link)
		if err != nil {
			fmt.Println("ERROR WITH ", err.Error())
		}
		links = append(links, link)
	}

	if rows.Err() != nil {
		fmt.Println("ERROR WITH ", rows.Err().Error())
	}

	data := make([]ItemData, 0, 10)

	for _, link := range links {
		data = append(data, run(fp, link)...)
	}

	// Send this data as an email

}

func run(fp *gofeed.Parser, link string) []ItemData {
	feed, err := fp.ParseURL(link)
	if err != nil {
		fmt.Println("NOT WORKING BECAUSE ", err.Error())
	}
	items := make([]ItemData, 0, 3)
	// fmt.Println(feed.Items)
	for _, v := range feed.Items {
		t := *v.PublishedParsed
		// Within the last 24 hours
		if !time.Now().Add(-time.Hour * 24).After(t) {
			items = append(items, ItemData{v.Title, v.Description})
		}
	}
	return items
}
