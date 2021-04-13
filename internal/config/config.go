package config

// Configs exported
type Configs struct {
	Log       LogConfig
	Webserver WebserverConfig
	Cookie    CookieConfig
	SSH       WGSSH
}

// LogConfig exported
type LogConfig struct {
	Level string
}

// WebserverConfig exported
type WebserverConfig struct {
	Port  int
	Debug bool
}

// CookieConfig exported
type CookieConfig struct {
	Name   string
	Secret string
	Domain string
}

// WGSSH exported
type WGSSH struct {
	ServerAddress string
	Port          string
	SSHUser       string
	SSHPrivateKey string
	SSHKnownHosts string
}
