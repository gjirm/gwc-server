package internal

import (
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"jirm.cz/gwc-server/internal/config"
	"jirm.cz/gwc-server/internal/ssh"
	"jirm.cz/gwc-server/internal/validate"
	"jirm.cz/gwc-server/internal/wg"
)

// MyServer server instance
func MyServer(log *logrus.Logger, config config.Configs) {

	log.Info("Starting Webserver")
	if !(config.Webserver.Debug) {
		gin.SetMode(gin.ReleaseMode)
	}
	// Disable Console Color
	gin.DisableConsoleColor()

	// Disable gin logging
	gin.DefaultWriter = ioutil.Discard

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {

		// Check if right cookie is present
		cookie, err := c.Cookie(config.Cookie.Name)
		if err != nil {
			msg := "Auth cookie " + config.Cookie.Name + " not found"
			log.Error(msg)
			c.String(400, msg)
			return
		}

		// Validate cookie
		valid, msg := validate.ValidateCookie(config, cookie)
		if valid {
			// Cookie is valid -> run SSH cmd - avtivate users WireGuard peers
			log.Info("Valid request by " + msg + " from IP " + c.ClientIP())
			user := strings.Split(msg, "@")
			command := user[0] + " " + c.ClientIP()
			log.Info("Running SSH command: " + command)
			c.String(200, ssh.RunSshCommand(log, config, command))
		} else {
			log.Error(msg)
			c.String(400, msg)
			return
		}
	})

	// Generate WireGuard keys
	r.GET("/generate", func(c *gin.Context) {
		log.Info("Generating WireGuard keys")
		privKey, pubKey, err := wg.GenerateWGKey()
		if err != nil {
			msg := "Unable to generate WireGuard Keys"
			log.Error(msg)
			c.String(400, msg)
		} else {
			c.JSON(200, gin.H{
				"PrivateKey": pubKey,
				"PublicKey":  privKey,
			})
		}

	})

	portNumber := ":" + strconv.Itoa(config.Webserver.Port)
	r.Run(portNumber)
}
