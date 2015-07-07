package notification

import (
	"fmt"

	"github.com/mailgun/mailgun-go"
)

// SendEmailMailgun sends a multipart text and HTML e-mail with a link to the click endpoint for stopping the notification
func (e *Engine) SendEmailMailgun(address, subject, plain, html string) error {
	from := fmt.Sprint("Chicken Little <chickenlittle@", e.Config.Integrations.Mailgun.Hostname, ">")

	mg := mailgun.NewMailgun(e.Config.Integrations.Mailgun.Hostname, e.Config.Integrations.Mailgun.APIKey, "")

	m := mg.NewMessage(from, subject, plain)
	m.SetHtml(html)
	m.AddRecipient(address)

	_, _, err := mg.Send(m)
	return err
}
