package utils

import (
	"golang.org/x/crypto/ssh"
	"ospback/logger"
)

func SSHRunCMD(sshSession *ssh.Session, command string) (string, error) {
	cmdrst, err := sshSession.CombinedOutput(command)
	rst := string(cmdrst)
	if err != nil {
		logger.Error("Command Failed:", command, err)
		return rst, err
	} else {
		logger.Info("Command Success:", command)
		return rst, nil
	}
}

func SSHRunCMDS(sshClient *ssh.Client, commands []string) {
	defer CloseSSHClient(sshClient)
	for _, cmd := range commands {
		ss, _ := SSHSession(sshClient)
		SSHRunCMD(ss, cmd)
		defer CloseSSHSession(ss)
	}
}
