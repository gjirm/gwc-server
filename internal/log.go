package gwc

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

// NewDefaultLogger creates a new logger based on the current configuration
func NewDefaultLogger() *logrus.Logger {
	// Setup logger
	//log = logrus.StandardLogger()

	// Init logging
	log = logrus.New()
	//Formatter := new(logrus.TextFormatter)
	Formatter := new(logrus.JSONFormatter)
	log.SetFormatter(Formatter)
	log.SetOutput(os.Stdout)

	// Set logger format
	// switch config.LogFormat {
	// case "pretty":
	// 	break
	// case "json":
	// 	logrus.SetFormatter(&logrus.JSONFormatter{})
	// // "text" is the default
	// default:
	// 	logrus.SetFormatter(&logrus.TextFormatter{
	// 		DisableColors: true,
	// 		FullTimestamp: true,
	// 	})
	// }

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

	return log
}
