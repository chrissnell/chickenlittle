package main

import (
	"fmt"
	"github.com/jpoehls/gophermail"
	"log"
	"net/mail"
	"net/smtp"
	"time"
)

func SendEmailSMTP(address, subject, plain, html string) {
	// Set up authentication information
	auth := smtp.PlainAuth(
		"",
		c.Config.Integrations.SMTP.Login,
		c.Config.Integrations.SMTP.Password,
		c.Config.Integrations.SMTP.Hostname,
	)

	from := mail.Address{Address: c.Config.Integrations.SMTP.Sender}
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
	host := fmt.Sprintf("%s:%d", c.Config.Integrations.SMTP.Hostname, c.Config.Integrations.SMTP.Port)
	err := gophermail.SendMail(
		host,
		auth,
		message,
	)
	if err != nil {
		log.Println("SMTP-Error:", err)
	}
}
