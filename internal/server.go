package gwc

import (
	"html/template"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// MyServer server instance
func MyServer() {

	if !(config.Webserver.Debug) {
		gin.SetMode(gin.ReleaseMode)
	}
	// Disable Console Color
	gin.DisableConsoleColor()

	// Disable gin logging
	gin.DefaultWriter = ioutil.Discard

	r := gin.Default()
	r.Use(location.Default())
	r.LoadHTMLGlob("templates/*")

	// Activate all peers of validated user
	r.GET("/", func(c *gin.Context) {

		if !(config.Api.ActivateAll) {
			c.String(400, "Not enabled")
			return
		}

		cLog := log.WithFields(logrus.Fields{
			"request":   "direct",
			"operation": "activate",
			"object":    "all",
		})
		// rURL := location.Get(c)
		// cLog.Info(rURL.Host)
		// rURL.

		// Check if right cookie is present
		cookie, err := c.Cookie(config.Cookie.Name)
		if err != nil {
			msg := "Auth cookie " + config.Cookie.Name + " not found"
			cLog.Error(msg)
			c.String(400, msg)
			return
		}

		// Validate cookie
		valid, msg := ValidateCookie(cookie)
		if valid {
			// Cookie is valid -> run SSH cmd - avtivate users WireGuard peers
			cLog.Info("Valid request by " + msg + " from IP " + c.ClientIP())

			user := strings.Split(msg, "@")
			command := user[0] + " " + c.ClientIP()
			cLog.Info("Running SSH command: " + command)
			sshOut, err := RunSshCommand(command)
			if err != nil {
				cLog.Debugf(sshOut, err)
				cLog.Error("Error running command")
				c.String(400, "Error running command")
				return
			}

			if strings.Contains(sshOut, "failed") || strings.Contains(sshOut, "Failed") {
				c.HTML(200, "status.html", gin.H{
					"message": sshOut,
					"alert":   "alert-danger",
				})
			} else {
				c.HTML(200, "status.html", gin.H{
					"message": sshOut,
					"alert":   "alert-success",
				})
			}

		} else {
			cLog.Error(msg)
			c.String(400, msg)
			return
		}
	})

	// List all peers of validated user
	r.GET("/list", func(c *gin.Context) {

		rURL := location.Get(c)

		cLog := log.WithFields(logrus.Fields{
			"request":   "direct",
			"operation": "list",
			"object":    "peers",
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
		valid, msg := ValidateCookie(cookie)
		if valid {
			// Cookie is valid -> run SSH cmd - list user WireGuard peers
			cLog.Info("Valid request by " + msg + " from IP " + c.ClientIP())

			user := strings.Split(msg, "@")

			command := user[0] + " " + c.ClientIP() + " list"

			cLog.Info("Running SSH command: " + command)
			sshOut, err := RunSshCommand(command)
			if err != nil {
				cLog.Debugf(sshOut, err)
				cLog.Error("Error running command")
				c.String(400, "Error running command")
				return
			}

			peersList := strings.Fields(sshOut)

			htmlList := ""
			deviceName := ""
			for i := 0; i < len(peersList); i++ {

				if strings.HasPrefix(peersList[i], "pc") {
					deviceName = "Computer (" + peersList[i] + ")"
				} else if strings.HasPrefix(peersList[i], "ntb") {
					deviceName = "Notebook (" + peersList[i] + ")"
				} else if strings.HasPrefix(peersList[i], "mac") {
					deviceName = "MacBook (" + peersList[i] + ")"
				} else {
					deviceName = "Mobile phone (" + peersList[i] + ")"
				}
				htmlList += "<a href='" + rURL.Scheme + "://" + c.Request.Host + "/activate/" + peersList[i] + "' class='list-group-item list-group-item-action list-group-item-primary'>" + deviceName + "</a></li>"
			}
			if htmlList != "" {
				c.HTML(200, "listPeers.html", gin.H{
					"user": user[0],
					"list": template.HTML(htmlList),
				})
			} else {
				c.HTML(200, "listPeers.html", gin.H{
					"user": user[0],
					"list": template.HTML("<a href='/list' class='list-group-item list-group-item-action list-group-item-primary disabled' >You have no devices</a></li>"),
				})
			}

		} else {
			cLog.Error(msg)
			c.String(400, msg)
			return
		}
	})

	// Show URL for generating and download config using token
	r.GET("/d/token/:token", func(c *gin.Context) {

		rURL := location.Get(c)

		cLog := log.WithFields(logrus.Fields{
			"request":   "direct",
			"operation": "token",
			"object":    "peer",
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
		valid, msg := ValidateCookie(cookie)
		if valid {
			// Cookie is valid -> run SSH cmd - list user WireGuard peers
			cLog.Info("Valid request by " + msg + " from IP " + c.ClientIP())

			token := c.Param("token")

			c.HTML(200, "downloadConfig.html", gin.H{
				"download": template.HTML(rURL.Scheme + "://" + c.Request.Host + "/token/" + token),
			})
		} else {
			cLog.Error(msg)
			c.String(400, msg)
			return
		}
	})

	// Activate users peer
	r.GET("/activate/:peer", func(c *gin.Context) {

		cLog := log.WithFields(logrus.Fields{
			"request":   "direct",
			"operation": "activate",
			"object":    "peer",
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
		valid, msg := ValidateCookie(cookie)
		if valid {
			// Cookie is valid -> run SSH cmd - avtivate users WireGuard peers
			cLog.Info("Valid request by " + msg + " from IP " + c.ClientIP())

			user := strings.Split(msg, "@")

			// Peer name
			peer := c.Param("peer")

			command := user[0] + " " + c.ClientIP() + " " + peer

			cLog.Info("Running SSH command: " + command)

			sshOut, err := RunSshCommand(command)
			if err != nil {
				cLog.Debugf(sshOut, err)
				cLog.Error("Error running command")
				c.String(400, "Error running command")
				return
			}

			if strings.Contains(sshOut, "failed") || strings.Contains(sshOut, "Failed") || strings.Contains(sshOut, "not found") {
				c.HTML(200, "status.html", gin.H{
					"message": sshOut,
					"alert":   "alert-danger",
				})
			} else {
				c.HTML(200, "status.html", gin.H{
					"message": sshOut,
					"alert":   "alert-success",
				})
			}

		} else {
			cLog.Error(msg)
			c.String(400, msg)
			return
		}
	})

	// Activate users peer
	r.GET("/token/:token", func(c *gin.Context) {

		cLog := log.WithFields(logrus.Fields{
			"request":   "direct",
			"operation": "token",
			"object":    "peer",
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
		valid, msg := ValidateCookie(cookie)
		if valid {
			// Cookie is valid -> run SSH cmd - avtivate users WireGuard peers
			cLog.Info("Valid request by " + msg + " from IP " + c.ClientIP())

			user := strings.Split(msg, "@")

			token := c.Param("token")

			command := user[0] + " " + c.ClientIP() + " token " + token

			cLog.Info("Running SSH command: " + command)

			sshOut, err := RunSshCommand(command)
			if err != nil {
				cLog.Debugf(sshOut, err)
				cLog.Error("Error running command")
				c.String(400, "Error running command")
				return
			}

			if strings.HasPrefix(sshOut, "[Interface]") {
				// Download configuration
				peerType := strings.Split(token, "-")
				c.Header("Content-Disposition", "attachment; filename=wg_"+user[0]+"_"+peerType[0]+".conf")
				c.Data(200, "application/octet-stream", []byte(sshOut))
			} else {
				if strings.Contains(sshOut, "failed") || strings.Contains(sshOut, "Failed") || strings.Contains(sshOut, "not valid") {
					c.HTML(200, "status.html", gin.H{
						"message": sshOut,
						"alert":   "alert-danger",
					})
				} else {
					c.String(200, sshOut)
				}
			}

		} else {
			cLog.Error(msg)
			c.String(400, msg)
			return
		}
	})

	// Create new token
	r.GET("/api/token/:suffix/:totp", func(c *gin.Context) {

		rURL := location.Get(c)

		if !(config.Api.Admin) {
			c.String(400, "Api Not enabled")
			return
		}

		cLog := log.WithFields(logrus.Fields{
			"request":   "api",
			"operation": "token",
			"object":    "peer",
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
		valid, msg := ValidateCookie(cookie)
		if valid {
			// Cookie is valid -> generate new WG configuration
			cLog.Info("Valid request by " + msg + " from IP " + c.ClientIP())

			user := strings.Split(msg, "@")
			suffix := c.Param("suffix")
			totp := c.Param("totp")

			command := user[0] + " " + c.ClientIP() + " " + totp + " " + "token " + suffix

			cLog.Info("Running SSH command: " + command)

			sshOut, err := RunSshCommand(command)
			if err != nil {
				cLog.Debugf(sshOut, err)
				cLog.Error("Error running command")
				c.String(400, "Error running command")
				return
			}

			if strings.Contains(sshOut, "failed") || strings.Contains(sshOut, "Failed") {
				c.String(200, sshOut)
			} else {
				outText := rURL.Scheme + "://" + c.Request.Host + "/d/token/" + sshOut
				c.String(200, outText)
			}

		} else {
			cLog.Error(msg)
			c.String(400, msg)
			return
		}

	})

	// Add new user peer to VPN
	r.GET("/api/add/:name/:suffix/:totp", func(c *gin.Context) {

		if !(config.Api.Admin) {
			c.String(400, "Api Not enabled")
			return
		}

		cLog := log.WithFields(logrus.Fields{
			"request":   "api",
			"operation": "add",
			"object":    "user",
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
		valid, msg := ValidateCookie(cookie)
		if valid {
			// Cookie is valid -> generate new WG configuration
			cLog.Info("Valid request by " + msg + " from IP " + c.ClientIP())

			user := strings.Split(msg, "@")
			name := c.Param("name")
			suffix := c.Param("suffix")
			totp := c.Param("totp")

			// :admin :ip :totp add :user :suffix
			command := user[0] + " " + c.ClientIP() + " " + totp + " " + "add " + name + " " + suffix

			cLog.Info("Running SSH command: " + command)

			sshOut, err := RunSshCommand(command)
			if err != nil {
				cLog.Debugf(sshOut, err)
				cLog.Error("Error running command")
				c.String(400, "Error running command")
				return
			}

			c.String(200, sshOut)

		} else {
			cLog.Error(msg)
			c.String(400, msg)
			return
		}

	})

	// List all users and peers
	r.GET("/api/list/users/:totp", func(c *gin.Context) {

		if !(config.Api.Admin) {
			c.String(400, "Api Not enabled")
			return
		}

		cLog := log.WithFields(logrus.Fields{
			"request":   "api",
			"operation": "list",
			"object":    "users",
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
		valid, msg := ValidateCookie(cookie)
		if valid {
			// Cookie is valid
			cLog.Info("Valid request by " + msg + " from IP " + c.ClientIP())

			user := strings.Split(msg, "@")
			totp := c.Param("totp")

			command := user[0] + " " + c.ClientIP() + " " + totp + " " + "list users"

			cLog.Info("Running SSH command: " + command)

			sshOut, err := RunSshCommand(command)
			if err != nil {
				cLog.Debugf(sshOut, err)
				cLog.Error("Error running command")
				c.String(400, "Error running command")
				return
			}

			c.String(200, sshOut)

		} else {
			cLog.Error(msg)
			c.String(400, msg)
			return
		}

	})

	// List activated peers
	r.GET("/api/list/activated/:totp", func(c *gin.Context) {

		if !(config.Api.Admin) {
			c.String(400, "Api Not enabled")
			return
		}

		cLog := log.WithFields(logrus.Fields{
			"request":   "api",
			"operation": "list",
			"object":    "activated",
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
		valid, msg := ValidateCookie(cookie)
		if valid {
			// Cookie is valid
			cLog.Info("Valid request by " + msg + " from IP " + c.ClientIP())

			user := strings.Split(msg, "@")
			totp := c.Param("totp")

			command := user[0] + " " + c.ClientIP() + " " + totp + " " + "list activated"

			cLog.Info("Running SSH command: " + command)

			sshOut, err := RunSshCommand(command)
			if err != nil {
				cLog.Debugf(sshOut, err)
				cLog.Error("Error running command")
				c.String(400, "Error running command")
				return
			}

			c.String(200, sshOut)

		} else {
			cLog.Error(msg)
			c.String(400, msg)
			return
		}

	})

	// Expire activated peer
	r.GET("/api/expire/:owner/:totp", func(c *gin.Context) {

		if !(config.Api.Admin) {
			c.String(400, "Api Not enabled")
			return
		}

		cLog := log.WithFields(logrus.Fields{
			"request":   "api",
			"operation": "expire",
			"object":    "user",
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
		valid, msg := ValidateCookie(cookie)
		if valid {
			// Cookie is valid
			cLog.Info("Valid request by " + msg + " from IP " + c.ClientIP())

			user := strings.Split(msg, "@")
			owner := c.Param("owner")
			totp := c.Param("totp")

			command := user[0] + " " + c.ClientIP() + " " + totp + " " + "expire " + owner

			cLog.Info("Running SSH command: " + command)

			sshOut, err := RunSshCommand(command)
			if err != nil {
				cLog.Debugf(sshOut, err)
				cLog.Error("Error running command")
				c.String(400, "Error running command")
				return
			}

			c.String(200, sshOut)

		} else {
			cLog.Error(msg)
			c.String(400, msg)
			return
		}

	})

	// Delete user
	r.GET("/api/del/:user/:totp", func(c *gin.Context) {

		if !(config.Api.Admin) {
			c.String(400, "Api Not enabled")
			return
		}

		cLog := log.WithFields(logrus.Fields{
			"request":   "api",
			"operation": "del",
			"user":      "user",
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
		valid, msg := ValidateCookie(cookie)
		if valid {
			// Cookie is valid
			cLog.Info("Valid request by " + msg + " from IP " + c.ClientIP())

			user := strings.Split(msg, "@")
			userDelete := c.Param("user")
			totp := c.Param("totp")

			command := user[0] + " " + c.ClientIP() + " " + totp + " " + "del " + userDelete

			cLog.Info("Running SSH command: " + command)

			sshOut, err := RunSshCommand(command)
			if err != nil {
				cLog.Debugf(sshOut, err)
				cLog.Error("Error running command")
				c.String(400, "Error running command")
				return
			}

			c.String(200, sshOut)

		} else {
			cLog.Error(msg)
			c.String(400, msg)
			return
		}

	})

	// Delete user peer
	r.GET("/api/del/peer/:user/:peer/:totp", func(c *gin.Context) {

		if !(config.Api.Admin) {
			c.String(400, "Api Not enabled")
			return
		}

		cLog := log.WithFields(logrus.Fields{
			"request":   "api",
			"operation": "del",
			"object":    "peer",
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
		valid, msg := ValidateCookie(cookie)
		if valid {
			// Cookie is valid
			cLog.Info("Valid request by " + msg + " from IP " + c.ClientIP())

			user := strings.Split(msg, "@")
			userDelete := c.Param("user")
			peerDelete := c.Param("peer")
			totp := c.Param("totp")

			command := user[0] + " " + c.ClientIP() + " " + totp + " " + "del " + userDelete + " " + peerDelete

			cLog.Info("Running SSH command: " + command)

			sshOut, err := RunSshCommand(command)
			if err != nil {
				cLog.Debugf(sshOut, err)
				cLog.Error("Error running command")
				c.String(400, "Error running command")
				return
			}

			c.String(200, sshOut)

		} else {
			cLog.Error(msg)
			c.String(400, msg)
			return
		}

	})

	// Generate WireGuard keys
	r.GET("/api/generate", func(c *gin.Context) {

		if !(config.Api.Admin) {
			c.String(400, "Api Not enabled")
			return
		}

		cLog := log.WithFields(logrus.Fields{
			"request":   "api",
			"operation": "generate",
			"object":    "keys",
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
		valid, msg := ValidateCookie(cookie)
		if valid {
			// Cookie is valid -> generate WireGuard keys
			cLog.Info("Valid request by " + msg + " from IP " + c.ClientIP())

			privKey, pubKey, err := GenerateWGKey()
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
