package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/mmcdole/gofeed"
)

type ItemData struct {
	title string
	desc  string
	link  string
}

// Overlay to run from main function through passing the database
func (s *server) runTask() {
	var task = func() {
		s.findUsers()
	}

	sc := gocron.NewScheduler(time.UTC)

	// Only for testing purposes
	// currentTime := time.Now().UTC()
	// currentTime = currentTime.Add(time.Second * 1)
	// timeString := currentTime.Format("15:04:05")
	// sc.Every(1).Day().At(timeString).Do(task)

	// Real statement
	sc.Every(1).Day().At("15:00:00").Do(task)

	sc.StartAsync()
}

// Function to find users

func (s *server) findUsers() {
	fp := gofeed.NewParser()
	rows, err := s.db.Query(context.Background(), "select id, email from users")
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

		s.singleUser(fp, id, email)
	}

}

// Function to run a single user's posts

func (s *server) singleUser(fp *gofeed.Parser, id int64, email string) {

	rows, err := s.db.Query(context.Background(), "select link from links where user_id = $1", id)
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
	// Change to HTML template
	emailString := htmlTemplateGenerator(data)
	// fmt.Println(emailString)

	// Send this data as an email
	s.sendEmail(email, emailString)
}

// Find the items that are published in the last day
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
			items = append(items, ItemData{v.Title, v.Description, v.Link})
		}
	}
	return items
}

// Generates HTML template from the given data
func htmlTemplateGenerator(data []ItemData) string {
	var emailString string

	for i, v := range data {
		if i >= 11 {
			break
		}
		emailString += strconv.Itoa(i)
		emailString += ". "
		// emailString += "TITLE: "
		emailString += v.title
		emailString += "\n\n"
		// emailString += "\n\n DESCRIPTION: "
		emailString += v.desc
		// emailString += "\n\n LINK: "
		emailString += "\n\n"
		emailString += v.link
		emailString += "\n\n\n\n"
	}

	return emailString
}
