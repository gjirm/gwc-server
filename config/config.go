package config

// Configs exported
type Configs struct {
	Log				LogConfig
	Webserver		WebserverConfig
	Database		DatabaseConfig
	Wireguard		WGConfig	
	// Database     DatabaseConfigurations
	// EXAMPLE_PATH string
	// EXAMPLE_VAR  string
}

// LogConfig exported
type LogConfig struct {
	File 	string
	Level	string
}

// WebserverConfig exported
type WebserverConfig struct {
	Port 		int
	Logfile 	string
	Debug		bool
}

// DatabaseConfig exported
type DatabaseConfig struct {
	Name 	string
}

// WGConfig exported
type WGConfig struct {
	SSH 		WGSSH
	WgServer	string	
	WgIPSpace	string
	WgIPBegin	string
}

// WGSSH exported
type WGSSH struct {
	ServerAddress	string
	SSHUser			string
	SSHPrivateKey	string
	SSHKnownHosts	string
}

// // DatabaseConfig exported
// type DatabaseConfig struct {
// 	DBName     string
// 	DBUser     string
// 	DBPassword string
// }

