package ssh

import (
	"bytes"
	"time"

	//"fmt"
	"io/ioutil"
	//"log"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	kh "golang.org/x/crypto/ssh/knownhosts"

	config "jirm.cz/gwc-server/internal/config"
)

// var (
// 	sshConfig
// )

// InitSSH exported
func InitSSH(log *logrus.Logger, config config.Configs) string {
	// user := "user"
	// address := "192.168.0.17"
	user := config.Wireguard.SSH.SSHUser
	address := config.Wireguard.SSH.ServerAddress
	command := "sudo wg"
	port := "22"

	// key, err := ioutil.ReadFile("/Users/user/.ssh/id_rsa")
	key, err := ioutil.ReadFile(config.Wireguard.SSH.SSHPrivateKey)
	if err != nil {
		log.Error("unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Error("unable to parse private key: %v", err)
	}

	// hostKeyCallback, err := kh.New("/Users/user/.ssh/known_hosts")
	hostKeyCallback, err := kh.New(config.Wireguard.SSH.SSHKnownHosts)
	if err != nil {
		log.Error("could not create hostkeycallback function: ", err)
	}

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			// Add in password check here for moar security.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: hostKeyCallback,

		// optional host key algo list
		HostKeyAlgorithms: []string{
			// ssh.KeyAlgoRSA,
			// ssh.KeyAlgoDSA,
			// ssh.KeyAlgoECDSA256,
			// ssh.KeyAlgoECDSA384,
			// ssh.KeyAlgoECDSA521,
			ssh.KeyAlgoED25519,
		},
		// optional tcp connect timeout
		Timeout: 5 * time.Second,
	}
	// Connect to the remote server and perform the SSH handshake.
	client, err := ssh.Dial("tcp", address+":"+port, sshConfig)
	if err != nil {
		log.Error("unable to connect: %v", err)
	}
	defer client.Close()
	ss, err := client.NewSession()
	if err != nil {
		log.Error("unable to create SSH session: ", err)
	}
	defer ss.Close()
	// Creating the buffer which will hold the remotly executed command's output.
	var stdoutBuf bytes.Buffer
	ss.Stdout = &stdoutBuf
	ss.Run(command)
	// Let's print out the result of command.
	//fmt.Println(stdoutBuf.String())
	return stdoutBuf.String()
}
