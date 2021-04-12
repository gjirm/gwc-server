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
			msg := "Wrong cookie set"
			log.Error(msg)
			c.String(400, msg)
			return
		}

		// Check cookie format
		parts := strings.Split(cookie, "|")
		if len(parts) != 3 {
			msg := "Wrong cookie format"
			log.Error(msg)
			c.String(400, msg)
			return
		}

		log.Info("New request by " + parts[2] + " from IP " + c.ClientIP())

		// Validate cookie
		valid := validate.ValidateCookie(log, config, parts[0], parts[1], parts[2]) // for key, value := range c.Request.Header {
		if valid {

			// Cookie is valid -> run SSH cmd
			log.Info("Running SSH command: " + config.SSH.Command)
			c.String(200, ssh.RunSshCommand(log, config))
		} else {
			msg := "Cookie is not valid!"
			c.String(400, msg)
			return
		}

	})

	// Test
	r.GET("/cookie", func(c *gin.Context) {
		cookie, err := c.Cookie(config.Cookie.Name)
		if err != nil {
			msg := "Wrong cookie set"
			log.Error(msg)
			c.String(400, msg)
			return
		}

		parts := strings.Split(cookie, "|")

		if len(parts) != 3 {
			msg := "Wrong cookie format"
			log.Error(msg)
			c.String(400, msg)
			return
		}

		valid := validate.ValidateCookie(log, config, parts[0], parts[1], parts[2])

		c.JSON(200, gin.H{
			"isValid": valid,
			"hash":    parts[0],
			"expires": parts[1],
			"mail":    parts[2],
		})
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
