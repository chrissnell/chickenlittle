package main

type Config struct {
	Service      ServiceConfig `yaml:"service"`
	Integrations Integrations  `yaml:"integrations"`
}

type ServiceConfig struct {
	ListenAddr string `yaml:"listen_address"`
	DBFile     string `yaml:"db_file"`
}

type Integrations struct {
	HipChat   HipChat   `yaml:"hipchat"`
	VictorOps VictorOps `yaml:"victorops"`
}

type VictorOps struct {
	APIKey string `yaml:"api_key"`
}

type HipChat struct {
	HipChatAuthToken    string `yaml:"hipchat_auth_token"`
	HipChatAnnounceRoom string `yaml:"hipchat_announce_room"`
}
