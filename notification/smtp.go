package notification

import (
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"time"

	"github.com/jpoehls/gophermail"
)

// SendEmailSMTP will send an email via SMTP
func (e *Engine) SendEmailSMTP(address, subject, plain, html string) {
	// Set up authentication information
	auth := smtp.PlainAuth(
		"",
		e.Config.Integrations.SMTP.Login,
		e.Config.Integrations.SMTP.Password,
		e.Config.Integrations.SMTP.Hostname,
	)

	from := mail.Address{Address: e.Config.Integrations.SMTP.Sender}
	to := mail.Address{Address: address}
	headers := mail.Header{}
	headers["Date"] = []string{time.Now().Format(time.RFC822Z)}

	message := &gophermail.Message{
		From:     from,
		To:       []mail.Address{to},
		Subject:  subject,
		Body:     plain,
		HTMLBody: html,
		Headers:  headers,
	}

	// Connect to the server, auth and send
	host := fmt.Sprintf("%s:%d", e.Config.Integrations.SMTP.Hostname, e.Config.Integrations.SMTP.Port)
	err := gophermail.SendMail(
		host,
		auth,
		message,
	)
	if err != nil {
		log.Println("SMTP-Error:", err)
	}
}
