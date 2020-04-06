package cls

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"ospback/utils"
)

// Server part
type Server interface {
	Conn() *ssh.Client
}

type OS struct {
	Ip       string
	Username string
	Password string
	Conntype string // ssh/rdp
}

type LinuxOS struct {
	Os         OS
	Sshkeypath string
	Sshport    int
}

func (self LinuxOS) Conn() *ssh.Client {
	sshclient, _ := utils.SSHClient(
		self.Os.Ip,
		self.Os.Username,
		self.Os.Password,
		self.Sshkeypath,
		self.Sshport,
	)
	return sshclient
}

// Backup data store
type BackupServer struct {
	Os        OS
	Srvtype   string // nfs/samba
	Datastore string // /datastore/ip/date/dirs
}

// Commands part

type Backup interface {
	PrepDataStore(Server)
	PrepDataTar(Server)
	PrjBackup(Server)
}

type Command struct {
	Cmd string
	Opt string
}

type BackupCMD struct {
	ComdDataStore Command
	ComdPrpDir    Command
	ComdBK        Command
	Projbase      string
	Project       map[string]string // cijenkins: "/var/lib/jenkins/"
	Backupserver  BackupServer
}

func (self BackupCMD) PrepDataStore(srv Server) {
	cmds := []string{}
	if self.Backupserver.Srvtype == "nfs" {
		cmdnfs := fmt.Sprintf(
			"%s %s %s:%s %s",
			self.ComdDataStore.Cmd,
			self.ComdDataStore.Opt,
			self.Backupserver.Os.Ip,
			self.Backupserver.Datastore,
			self.Projbase,
		)
		cmdmkdir := fmt.Sprintf(
			"%s %s %s",
			self.ComdPrpDir.Cmd,
			self.ComdPrpDir.Opt,
			self.Projbase,
		)
		cmds = append(cmds, cmdnfs)
		cmds = append(cmds, cmdmkdir)

		for k, _ := range self.Project {
			cmdmkdirprj := fmt.Sprintf(
				"%s %s %s/%s",
				self.ComdPrpDir.Cmd,
				self.ComdPrpDir.Opt,
				self.Projbase,
				k,
			)
			cmds = append(cmds, cmdmkdirprj)
		}

		utils.SSHRunCMDS(srv.Conn(), cmds)
	}
}

func (self BackupCMD) PrepDataTar(srv Server) {
}

func (self BackupCMD) PrjBackup(srv Server) {
	cmds := []string{}
	for k, v := range self.Project {
		cmdprjbk := fmt.Sprintf(
			"%s %s %s %s/%s",
			self.ComdBK.Cmd,
			self.ComdBK.Opt,
			v,
			self.Projbase,
			k,
		)
		cmds = append(cmds, cmdprjbk)
		fmt.Println(cmds)
	}

	// handle the project dir, backup the tar to the datastore/proj
	utils.SSHRunCMDS(srv.Conn(), cmds)
}
