package gwc

import (
	"bytes"
	"time"

	"io/ioutil"

	"golang.org/x/crypto/ssh"
	kh "golang.org/x/crypto/ssh/knownhosts"
)

// RunSshCommand exported
func RunSshCommand(command string) (string, error) {

	// Read SSH private key from file
	key, err := ioutil.ReadFile(config.SSH.SSHPrivateKey)
	if err != nil {
		return "SSH unable to read private key", err
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return "SSH unable to parse private key", err
	}

	// Read known_hosts from file
	hostKeyCallback, err := kh.New(config.SSH.SSHKnownHosts)
	if err != nil {
		return "SSH could not create hostkeycallback function", err
	}

	sshConfig := &ssh.ClientConfig{
		User: config.SSH.SSHUser,
		Auth: []ssh.AuthMethod{
			// Add in password check here for moar security.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: hostKeyCallback,

		// optional host key algo list
		HostKeyAlgorithms: []string{
			ssh.KeyAlgoECDSA256,
			// ssh.KeyAlgoECDSA384,
			// ssh.KeyAlgoECDSA521,
			// ssh.KeyAlgoED25519,
		},
		// optional tcp connect timeout
		Timeout: 5 * time.Second,
	}
	// Connect to the remote server and perform the SSH handshake.
	client, err := ssh.Dial("tcp", config.SSH.ServerAddress+":"+config.SSH.Port, sshConfig)
	if err != nil {
		return "SSH unable to connect", err
	}
	defer client.Close()
	ss, err := client.NewSession()
	if err != nil {
		return "SSH unable to create SSH session", err
	}
	defer ss.Close()
	// Creating the buffer which will hold the remotly executed command's output.
	var stdoutBuf bytes.Buffer
	ss.Stdout = &stdoutBuf
	ss.Run(command)

	return stdoutBuf.String(), nil
}
