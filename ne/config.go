package ne

type Config struct {
	Service ServiceConfig
	Twilio  TwilioConfig
	Mailgun MailgunConfig
	SMTP    SMTPConfig
}

type ServiceConfig struct {
	ClickURLBase    string
	CallbackURLBase string
}

type TwilioConfig struct {
	AccountSID     string
	AuthToken      string
	CallFromNumber string
	APIBaseURL     string
}

type MailgunConfig struct {
	Enabled  bool
	APIKey   string
	Hostname string
}

type SMTPConfig struct {
	Hostname string
	Port     int
	Login    string
	Password string
	Sender   string
}
