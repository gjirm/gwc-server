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

// TODO: .. Configuration https://stackoverflow.com/questions/16465705/how-to-handle-configuration-in-go

// MyServer server instance
func MyServer(log *logrus.Logger, config config.Configs) {

	log.Info("Starting Webserver")
	if !(config.Webserver.Debug) {
		gin.SetMode(gin.ReleaseMode)
	}
	// Disable Console Color, you don't need console color when writing the logs to file.
	gin.DisableConsoleColor()

	// Logging to a file.
	//log.Info("Webserver access log file: " + config.Webserver.Logfile)
	//f, _ := os.OpenFile(config.Webserver.Logfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	//gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	gin.DefaultWriter = ioutil.Discard

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
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

		valid := validate.ValidateCookie(log, config, parts[0], parts[1], parts[2]) // for key, value := range c.Request.Header {
		if valid {
			c.String(200, ssh.RunSshCommand(log, config))
		} else {
			msg := "Cookie is not valid!"
			log.Error(msg)
			c.String(400, msg)
			return
		}

	})

	r.GET("/ping", func(c *gin.Context) {
		//c.JSON(200, c.Request.Header)
		var cki, err = c.Request.Cookie("foo")
		if err != nil {
			log.Error("Unable to obtain cookie: foo")
			c.String(400, "Unable to obtain cookie: foo")
		} else {
			c.JSON(200, cki.Value)
		}
	})

	r.GET("/cookie", func(c *gin.Context) {
		//fmt.Println(c.Request.Header)
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

		valid := validate.ValidateCookie(log, config, parts[0], parts[1], parts[2]) // for key, value := range c.Request.Header {

		c.JSON(200, gin.H{
			"isValid": valid,
			"hash":    parts[0],
			"expires": parts[1],
			"mail":    parts[2],
		})
	})

	r.GET("/ssh", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"sshOut": ssh.RunSshCommand(log, config),
		})
	})

	r.GET("/generate", func(c *gin.Context) {
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
	r.Run(portNumber) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
