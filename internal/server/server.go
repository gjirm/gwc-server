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
		cLog := log.WithField("action", "activate")

		// Check if right cookie is present
		cookie, err := c.Cookie(config.Cookie.Name)
		if err != nil {
			msg := "Auth cookie " + config.Cookie.Name + " not found"
			cLog.Error(msg)
			c.String(400, msg)
			return
		}
		// Validate cookie
		valid, msg := validate.ValidateCookie(config, cookie)
		if valid {
			// Cookie is valid -> run SSH cmd - avtivate users WireGuard peers
			cLog.Info("Valid request by " + msg + " from IP " + c.ClientIP())

			user := strings.Split(msg, "@")
			command := user[0] + " " + c.ClientIP()
			cLog.Info("Running SSH command: " + command)
			c.String(200, ssh.RunSshCommand(cLog, config, command))
		} else {
			cLog.Error(msg)
			c.String(400, msg)
			return
		}
	})

	// Add new user peer to VPN
	// Peer name: <user>-<suffixe>
	r.GET("/api/add/:name/:suffix/:totp", func(c *gin.Context) {

		cLog := log.WithFields(logrus.Fields{
			"action":    "api",
			"operation": "add",
		})

		// Check if right cookie is present
		cookie, err := c.Cookie(config.Cookie.Name)
		if err != nil {
			msg := "Auth cookie " + config.Cookie.Name + " not found"
			cLog.Error(msg)
			c.String(400, msg)
			return
		}

		// Validate cookie
		valid, msg := validate.ValidateCookie(config, cookie)
		if valid {
			// Cookie is valid -> generate new WG configuration
			cLog.Info("Valid request by " + msg + " from IP " + c.ClientIP())

			user := strings.Split(msg, "@")
			name := c.Param("name")
			suffix := c.Param("suffix")
			totp := c.Param("totp")

			command := user[0] + " " + c.ClientIP() + " " + "add " + name + " " + suffix + " " + totp

			cLog.Info("Running SSH command: " + command)
			c.String(200, ssh.RunSshCommand(cLog, config, command))
			//c.String(200, command)

		} else {
			cLog.Error(msg)
			c.String(400, msg)
			return
		}

	})

	// Generate WireGuard keys
	r.GET("/api/generate", func(c *gin.Context) {

		cLog := log.WithFields(logrus.Fields{
			"action":    "api",
			"operation": "generate",
		})

		// Check if right cookie is present
		cookie, err := c.Cookie(config.Cookie.Name)
		if err != nil {
			msg := "Auth cookie " + config.Cookie.Name + " not found"
			cLog.Error(msg)
			c.String(400, msg)
			return
		}

		// Validate cookie
		valid, msg := validate.ValidateCookie(config, cookie)
		if valid {
			// Cookie is valid -> generate WireGuard keys
			cLog.Info("Valid request by " + msg + " from IP " + c.ClientIP())

			privKey, pubKey, err := wg.GenerateWGKey()
			if err != nil {
				msg := "Unable to generate WireGuard Keys"
				cLog.Error(msg)
				c.String(400, msg)
			} else {
				c.JSON(200, gin.H{
					"PrivateKey": pubKey,
					"PublicKey":  privKey,
				})
			}

		} else {
			cLog.Error(msg)
			c.String(400, msg)
			return
		}
	})

	portNumber := ":" + strconv.Itoa(config.Webserver.Port)
	r.Run(portNumber)
}
