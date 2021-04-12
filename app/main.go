package main

import (
	"os"

	c "jirm.cz/gwc-server/internal/config"
	server "jirm.cz/gwc-server/internal/server"

	//db "jirm.cz/gwc-server/db"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func readConfig(log *logrus.Logger, configName, configType string) c.Configs {
	log.Info("Reading configuration")
	// Set the file name of the configurations file
	viper.SetConfigName(configName)

	// Set the path to look for the configurations file
	viper.AddConfigPath(".")
	viper.AddConfigPath("/gwc")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType(configType)
	var configuration c.Configs

	if err := viper.ReadInConfig(); err != nil {
		log.Error("Error reading config file, %s", err)
	}

	// Set undefined variables
	//viper.SetDefault("database.dbname", "test_db")

	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Error("Unable to decode into struct, %v", err)
	}

	return configuration
}

func main() {

	// Init logging
	var log = logrus.New()
	Formatter := new(logrus.TextFormatter)
	//Formatter.TimestampFormat = "2006-01-02 15:04:05"
	Formatter.FullTimestamp = true
	Formatter.DisableColors = true
	log.SetFormatter(Formatter)
	log.SetOutput(os.Stdout)

	// Init config
	config := readConfig(log, "config", "yml")

	switch config.Log.Level {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
		log.Warnf("Home: invalid log level supplied: '%s'", config.Log.Level)
	}

	log.SetOutput(os.Stdout)

	// Init webserver
	server.MyServer(log, config)

}
