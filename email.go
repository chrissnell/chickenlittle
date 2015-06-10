package main

func SendEmail(address, message, uuid string) {
	log.Println("[", uuid, "] Sending email to:", address)
	plain := fmt.Sprint("You've received a message from the Chicken Little alert system:\n\n", message,
		"\n\n", "Stop notifications for this alert: ", c.Config.Service.ClickURLBase, "/", uuid, "/stop")
	html := fmt.Sprint("<HTML><BODY>You've received a message from the Chicken Little alert system:<BR><BR>",
		message, "<BR><BR><A HREF='", c.Config.Service.ClickURLBase, "/", uuid, "/stop'>Stop notifications for this alert</A></BODY></HTML>")

	if c.Config.Integrations.Mailgun.Enabled {
		SendEmailMailgun(address, plain, html)
	} else {
		SendEmailSMTP(address, plain, html)
	}
}
