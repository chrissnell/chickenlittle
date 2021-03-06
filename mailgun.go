package main

import (
	"fmt"
	"github.com/mailgun/mailgun-go"
)

// Sends a multipart text and HTML e-mail with a link to the click endpoint for stopping the notification
func SendEmailMailgun(address, subject, plain, html string) {
	from := fmt.Sprint("Chicken Little <chickenlittle@", c.Config.Integrations.Mailgun.Hostname, ">")

	mg := mailgun.NewMailgun(c.Config.Integrations.Mailgun.Hostname, c.Config.Integrations.Mailgun.APIKey, "")

	m := mg.NewMessage(from, subject, plain)
	m.SetHtml(html)
	m.AddRecipient(address)

	mg.Send(m)
}
