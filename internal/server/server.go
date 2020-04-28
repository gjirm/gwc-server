package internal

import (
	"strings"
	"strconv"
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	db "jirm.cz/gwc-server/db"
	config "jirm.cz/gwc-server/config"
	ssh "jirm.cz/gwc-server/internal/ssh"
	wg "jirm.cz/gwc-server/internal/wg"
)

// TODO: .. Configuration https://stackoverflow.com/questions/16465705/how-to-handle-configuration-in-go

// getNewIP 
func getNewIP() (string) {
	ipCounter := db.GetValue("options", "ipcounter")
	ipPrefix := db.GetValue("options", "ipspace")
	newIP := ipPrefix + "." + ipCounter
	tmpIPCounter, _ := strconv.Atoi(ipCounter)
	tmpIPCounter++
	ipCounter =strconv.Itoa(tmpIPCounter)
	db.PutValue("options", "ipcounter", ipCounter)
	return newIP
}

// MyServer server instance
func MyServer(log *logrus.Logger, config config.Configs) {
	
	log.Info("Starting Webserver")
	if !(config.Webserver.Debug) {
		gin.SetMode(gin.ReleaseMode)
	}
	db.InitDB(log, config)
    // Disable Console Color, you don't need console color when writing the logs to file.
    gin.DisableConsoleColor()

	// Logging to a file.
	log.Info("Webserver access log file: " + config.Webserver.Logfile)
    f, _ := os.OpenFile(config.Webserver.Logfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
    gin.DefaultWriter = io.MultiWriter(f, os.Stdout) 

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")

		if kv := db.GetValue("keys", name); kv == "" {
			log.Info("Generating new key for user: ", name )
			
			pKey, pubKey, err := wg.GenerateWGKey()
			
			if err != nil {
				log.Error("Failed to generate Wireguard key: ", err)
			
				c.JSON(500, gin.H{
					"Error":  err,
				})
			} else {
				db.PutValue("keys", name, pubKey)

				ip := getNewIP()
				db.PutValue("ip", name, ip)
				
				c.JSON(200, gin.H{
					"Name":  name,
					"newPrivateKey": pKey,
					"newPublicKey": pubKey,
					"ip": ip,
				})
			}
		} else {
			log.Info("Key for user " + name + " already exists" )
			pubKey := db.GetValue("keys", name)
			ip := db.GetValue("ip", name)
			c.JSON(200, gin.H{
				"Name":  name,
				"publicKey": pubKey,
				"ip": ip,

			})
		}





	})

	r.GET("/jirm", func(c *gin.Context) {
		//fmt.Println(c.Request.Header)
		cookie, _ := c.Cookie("_authreponse")
		parts := strings.Split(cookie, "|")

		if len(parts) != 3 {
			return
		}

		// for key, value := range c.Request.Header {
		// 	// for each pair in the map, print key and value
		// 	fmt.Printf("key: %s, value: %s\n", key, value)
		// }
		c.JSON(200, gin.H{
			"hash":  parts[0],
			"mail":  parts[1],
			"valid": parts[2],
		})
	})

	r.GET("/test", func(c *gin.Context) {
		//fmt.Println(c.Request.Header)
		cookie, _ := c.Cookie("_authreponse")
		parts := strings.Split(cookie, "|")

		if len(parts) != 3 {
			return
		}

		// for key, value := range c.Request.Header {
		// 	// for each pair in the map, print key and value
		// 	fmt.Printf("key: %s, value: %s\n", key, value)
		// }
		myname := strings.Split(parts[1], "@")
		mypubkey := db.GetValue("pubkeys", myname[0])
		c.JSON(200, gin.H{
			"pubkey":  mypubkey,
		})
	})

	r.GET("/ssh", func(c *gin.Context) {
		
		c.JSON(200, gin.H{
			"sshOut":  ssh.InitSSH(log, config),
		})
	})
	
	r.GET("/generate", func(c *gin.Context) {
		wg.GenerateWGKey()
		c.JSON(200, gin.H{
			"ahoj": "svete",
		})
	})

	portNumber := ":" + strconv.Itoa(config.Webserver.Port)
	r.Run(portNumber) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
