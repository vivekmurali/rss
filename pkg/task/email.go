package task

import (
	"context"
	"fmt"
	"rss/pkg/db"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

var domain string = "vivekmurali.in"

func sendEmail(recipient string, body string) {
	mg := mailgun.NewMailgun(domain, db.PrivateAPIKey)
	sender := "vivek@vivekmurali.in"
	//Add date to the subject
	subject := "Your RSS feed for the day"

	message := mg.NewMessage(sender, subject, body, recipient)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		fmt.Println("NOT WORKING BECAUSE: ", err.Error())
	}
	fmt.Printf("ID: %s Resp: %s\n", id, resp)
}
