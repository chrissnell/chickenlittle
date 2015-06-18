package ne

import (
	"fmt"

	"github.com/mailgun/mailgun-go"
)

// Sends a multipart text and HTML e-mail with a link to the click endpoint for stopping the notification
func (e *Engine) SendEmailMailgun(address, subject, plain, html string) {
	from := fmt.Sprint("Chicken Little <chickenlittle@", e.Config.Mailgun.Hostname, ">")

	mg := mailgun.NewMailgun(e.Config.Mailgun.Hostname, e.Config.Mailgun.APIKey, "")

	m := mg.NewMessage(from, subject, plain)
	m.SetHtml(html)
	m.AddRecipient(address)

	mg.Send(m)
}
