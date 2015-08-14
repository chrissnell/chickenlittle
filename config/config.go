package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config contains the ChickenLittle configuration
type Config struct {
	Service      ServiceConfig `yaml:"service"`
	Integrations Integrations  `yaml:"integrations"`
}

// ServiceConfig contains the webserver and DB configuration (URLs, Ports, ...)
type ServiceConfig struct {
	APIListenAddr      string `yaml:"api_listen_address"`
	ClickListenAddr    string `yaml:"click_listen_address"`
	ClickURLBase       string `yaml:"click_url_base"`
	CallbackListenAddr string `yaml:"callback_listen_address"`
	CallbackURLBase    string `yaml:"callback_url_base"`
	DBFile             string `yaml:"db_file"`
}

// Integrations contains the configuration for our plugin integrations
type Integrations struct {
	HipChat   HipChat   `yaml:"hipchat"`
	VictorOps VictorOps `yaml:"victorops"`
	Twilio    Twilio    `yaml:"twilio"`
	Mailgun   Mailgun   `yaml:"mailgun"`
	SMTP      SMTP      `yaml:"smtp"`
}

// Twilio contains the twilio API credentials. Twilio is an telephony SaaS provider.
type Twilio struct {
	AccountSID     string `yaml:"account_sid"`
	AuthToken      string `yaml:"auth_token"`
	CallFromNumber string `yaml:"call_from_number"`
	APIBaseURL     string `yaml:"api_base_url"`
}

// Mailgun contains the Mailgun API credentials. Mailgun is an email SaaS provider.
type Mailgun struct {
	Enabled  bool   `yaml:"enabled"`
	APIKey   string `yaml:"api_key"`
	Hostname string `yaml:"hostname"`
}

// SMTP contains the SMTP client credentials.
type SMTP struct {
	Hostname string `yaml:"hostname"`
	Port     int    `yaml:"port"`
	Login    string `yaml:"login"`
	Password string `yaml:"password"`
	Sender   string `yaml:"sender"`
}

// VictorOps contains the VictorOps API credentials.
type VictorOps struct {
	APIKey string `yaml:"api_key"`
}

// HipChat contains the HipChat credentials.
type HipChat struct {
	HipChatAuthToken    string `yaml:"hipchat_auth_token"`
	HipChatAnnounceRoom string `yaml:"hipchat_announce_room"`
}

// New creates an new config object from the given filename.
func New(filename string) (Config, error) {
	cfgFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}
	c := Config{}
	err = yaml.Unmarshal(cfgFile, &c)
	if err != nil {
		return Config{}, err
	}
	return c, nil
}

// NewDefault creates valid default config.
func NewDefault() Config {
	c := Config{
		Service: ServiceConfig{
			APIListenAddr:      ":21001",
			ClickListenAddr:    ":21002",
			ClickURLBase:       "http://localhost:21002/",
			CallbackListenAddr: ":21003",
			CallbackURLBase:    "http://localhost:21003/",
			DBFile:             "chickenlittle.db",
		},
		Integrations: Integrations{
			HipChat:   HipChat{},
			VictorOps: VictorOps{},
			Twilio:    Twilio{},
			Mailgun:   Mailgun{},
			SMTP: SMTP{
				Hostname: "localhost",
				Port:     25,
				Sender:   "chickenlittle@localhost",
			},
		},
	}
	return c
}
