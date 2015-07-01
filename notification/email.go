package notification

import (
	"fmt"
	"log"
)

// SendEmail will send an email using an appropriate integration.
func (e *Engine) SendEmail(address, message, uuid string) error {
	log.Println("[", uuid, "] Sending email to:", address)
	subject := "Chicken Little message received"
	plain := fmt.Sprint("You've received a message from the Chicken Little alert system:\n\n", message,
		"\n\n", "Stop notifications for this alert: ", e.Config.Service.ClickURLBase, "/", uuid, "/stop")
	html := fmt.Sprint("<HTML><BODY>You've received a message from the Chicken Little alert system:<BR><BR>",
		message, "<BR><BR><A HREF='", e.Config.Service.ClickURLBase, "/", uuid, "/stop'>Stop notifications for this alert</A></BODY></HTML>")

	if e.Config.Integrations.Mailgun.Enabled {
		return e.SendEmailMailgun(address, subject, plain, html)
	}
	return e.SendEmailSMTP(address, subject, plain, html)
}
