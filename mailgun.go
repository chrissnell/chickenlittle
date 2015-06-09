package main

import (
	"fmt"
	"github.com/mailgun/mailgun-go"
	"log"
)

// Sends a multipart text and HTML e-mail with a link to the click endpoint for stopping the notification
func SendEmail(address, message, uuid string) {

	log.Println("[", uuid, "] Sending email to:", address)

	from := fmt.Sprint("Chicken Little <chickenlittle@", c.Config.Integrations.Mailgun.Hostname, ">")

	plain := fmt.Sprint("You've received a message from the Chicken Little alert system:\n\n", message,
		"\n\n", "Stop notifications for this alert: ", c.Config.Service.ClickURLBase, "/", uuid, "/stop")

	html := fmt.Sprint("<HTML><BODY>You've received a message from the Chicken Little alert system:<BR><BR>",
		message, "<BR><BR><A HREF='", c.Config.Service.ClickURLBase, "/", uuid, "/stop'>Stop notifications for this alert</A></BODY></HTML>")

	mg := mailgun.NewMailgun(c.Config.Integrations.Mailgun.Hostname, c.Config.Integrations.Mailgun.APIKey, "")

	m := mg.NewMessage(from, "Chicken Little message received", plain)
	m.SetHtml(html)
	m.AddRecipient(address)

	mg.Send(m)
}
