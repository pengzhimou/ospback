package utils

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"ospback/logger"
	"time"
)

func publicKeyAuthFunc(kPath string) ssh.AuthMethod {
	keyPath, err := homedir.Expand(kPath)
	if err != nil {
		logger.Error("find key's home dir failed", err)
	}
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		logger.Error("ssh key file read failed", err)
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		logger.Error("ssh key signer failed", err)
	}
	return ssh.PublicKeys(signer)
}

func SSHClient(sshHost, sshUser, sshPassword, sshKeyPath string, sshPort int) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		Timeout:         10 * time.Second,
		User:            sshUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	if sshKeyPath == "" {
		config.Auth = []ssh.AuthMethod{ssh.Password(sshPassword)}
	} else {
		config.Auth = []ssh.AuthMethod{publicKeyAuthFunc(sshKeyPath)}
	}

	addr := fmt.Sprintf("%s:%d", sshHost, sshPort)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		logger.Error("create ssh client fail:", err)
		defer func() {
			fmt.Println("defer close client!!!!!")
			CloseSSHClient(sshClient)
		}()
		return nil, err
	} else {
		logger.Info("create ssh client success:", sshHost, sshUser, sshPassword, sshPort, sshKeyPath)
		return sshClient, nil
	}
}

func SSHSession(sshClient *ssh.Client) (*ssh.Session, error) {
	sshSession, err := sshClient.NewSession()
	if err != nil {
		logger.Error("create ssh session fail", err)
		return nil, err
	} else {
		logger.Info("create ssh session success")
		return sshSession, nil
	}
}

func CloseSSHClient(client *ssh.Client) {
	logger.Warn("ssh client is closing")
	client.Close()
}

func CloseSSHSession(session *ssh.Session) {
	logger.Warn("ssh session is closing")
	session.Close()
}
