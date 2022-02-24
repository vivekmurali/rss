package main

import (
	"context"
	"fmt"
	"time"
)

var domain string = "vivekmurali.in"

func (s *server) sendEmail(recipient string, body string) {
	// mg := mailgun.NewMailgun(domain, s.privateAPIKey)
	sender := "vivek@vivekmurali.in"
	//Add date to the subject
	subject := "Your RSS feed for the day"

	message := s.mg.NewMessage(sender, subject, body, recipient)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, id, err := s.mg.Send(ctx, message)

	if err != nil {
		fmt.Println("NOT WORKING BECAUSE: ", err.Error())
	}
	fmt.Printf("ID: %s Resp: %s\n", id, resp)
}
