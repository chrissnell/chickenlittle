package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Service      ServiceConfig `yaml:"service"`
	Integrations Integrations  `yaml:"integrations"`
}

type ServiceConfig struct {
	APIListenAddr      string `yaml:"api_listen_address"`
	ClickListenAddr    string `yaml:"click_listen_address"`
	ClickURLBase       string `yaml:"click_url_base"`
	CallbackListenAddr string `yaml:"callback_listen_address"`
	CallbackURLBase    string `yaml:"callback_url_base"`
	DBFile             string `yaml:"db_file"`
}

type Integrations struct {
	HipChat   HipChat   `yaml:"hipchat"`
	VictorOps VictorOps `yaml:"victorops"`
	Twilio    Twilio    `yaml:"twilio"`
	Mailgun   Mailgun   `yaml:"mailgun"`
	SMTP      SMTP      `yaml:"smtp"`
}

type Twilio struct {
	AccountSID     string `yaml:"account_sid"`
	AuthToken      string `yaml:"auth_token"`
	CallFromNumber string `yaml:"call_from_number"`
	APIBaseURL     string `yaml:"api_base_url"`
}

type Mailgun struct {
	Enabled  bool   `yaml:"enabled"`
	APIKey   string `yaml:"api_key"`
	Hostname string `yaml:"hostname"`
}

type SMTP struct {
	Hostname string `yaml:"hostname"`
	Port     int    `yaml:"port"`
	Login    string `yaml:"login"`
	Password string `yaml:"password"`
	Sender   string `yaml:"sender"`
}

type VictorOps struct {
	APIKey string `yaml:"api_key"`
}

type HipChat struct {
	HipChatAuthToken    string `yaml:"hipchat_auth_token"`
	HipChatAnnounceRoom string `yaml:"hipchat_announce_room"`
}

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
