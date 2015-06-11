package main

import (
	"fmt"
	"log"
)

func SendEmail(address, message, uuid string) {
	log.Println("[", uuid, "] Sending email to:", address)
	subject := "Chicken Little message received"
	plain := fmt.Sprint("You've received a message from the Chicken Little alert system:\n\n", message,
		"\n\n", "Stop notifications for this alert: ", c.Config.Service.ClickURLBase, "/", uuid, "/stop")
	html := fmt.Sprint("<HTML><BODY>You've received a message from the Chicken Little alert system:<BR><BR>",
		message, "<BR><BR><A HREF='", c.Config.Service.ClickURLBase, "/", uuid, "/stop'>Stop notifications for this alert</A></BODY></HTML>")

	if c.Config.Integrations.Mailgun.Enabled {
		SendEmailMailgun(address, subject, plain, html)
	} else {
		SendEmailSMTP(address, subject, plain, html)
	}
}
