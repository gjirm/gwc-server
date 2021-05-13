package main

import (
	internal "jirm.cz/gwc-server/internal"
)

func main() {

	// Parse options
	config := internal.NewGlobalConfig()

	// Setup logger
	log := internal.NewDefaultLogger()

	// ##################################
	// Init webserver
	log.Info("Starting GWC Server")
	log.Infof("Listening on :%d", config.Webserver.Port)

	internal.MyServer()

}
