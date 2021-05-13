package gwc

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var config *Configs

// Configs exported
type Configs struct {
	Log       LogConfig
	Webserver WebserverConfig
	Cookie    CookieConfig
	SSH       WGSSH
	Api       ApiConfig
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

// ApiConfig exported
type ApiConfig struct {
	Admin       bool
	ActivateAll bool
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

func readConfig(configName, configType string) (*Configs, error) {
	//log.Info("Reading configuration")
	// Set the file name of the configurations file
	viper.SetConfigName(configName)

	// Set the path to look for the configurations file
	viper.AddConfigPath(".")
	viper.AddConfigPath("/gwc")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType(configType)
	//var configuration c.Configs

	if err := viper.ReadInConfig(); err != nil {
		//log.Error("Error reading config file, %s", err)
		return config, err
	}

	// Set undefined variables
	//viper.SetDefault("database.dbname", "test_db")

	err := viper.Unmarshal(&config)
	if err != nil {
		//log.Error("Unable to decode into struct, %v", err)
		return config, err
	}

	return config, nil
}

// NewGlobalConfig creates a new global config from file or arguments
func NewGlobalConfig() *Configs {

	var err error
	// Init config
	config, err := readConfig("config", "yml")
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
	return config
}
