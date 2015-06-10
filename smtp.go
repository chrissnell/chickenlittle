package main

import ()

func SendEmailSMTP(address, plain, html string) {
	// Set up authentication information
	auth := smtp.PlainAuth(
		"",
		c.Config.Integrations.SMTP.Login,
		c.Config.Integrations.SMTP.Password,
		c.Config.Integrations.SMTP.Hostname,
	)
	// TODO set plain and HTML multi-parts
	// Connect to the server, auth and send
	err := smtp.SendMail(
		c.Config.Integrations.SMTP.Hostname+":"+c.Config.Integrations.SMTP.Port,
		auth,
		c.Config.Integrations.SMTP.Sender,
		[]string{address},
		[]byte(plain),
	)
	if err != nil {
		log.Println(err)
	}
}
