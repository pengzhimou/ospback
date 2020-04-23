package backup

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
	"os"
	"ospback/utils"
	"time"
)

//var Jobtime string = time.Now().Format("20060102_150405")

type Server interface {
	Conn() *ssh.Client
}

type Backup interface {
	NFSMount(Server)
	PrjBackup(Server)
}

// server part
type ServerConfigs struct {
	LinuxOsConfig map[string]LinuxOS `yaml:"serverconfigs"`
}

type LinuxOS struct {
	Ip      string      `yaml:"ip"`
	Conntp  Conntype    `yaml:"conntype"`
	Bcktsks BackupTasks `yaml:"backuptasks"`
}

type Conntype struct {
	SshCn    SshConn    `yaml:"sshconn"`
	TelnetCn TelnetConn `yaml:"telnetconn"`
}

type SshConn struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Keypath  string `yaml:"keypath"`
	Port     int    `yaml:"port"`
}

type TelnetConn struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Port     int    `yaml:"port"`
}

type BackupTasks struct {
	Projbase string            `yaml:"projbase"`
	Tasks    map[string]string `yaml:"tasks"`
}

// backup server part
type BKServerConfigs struct {
	BKSCfgs BKServerType `yaml:"bkserverconfigs"`
}

type BKServerType struct {
	NFSSrv NFSServer   `yaml:"nfsserver"`
	SmbSrv SambaServer `yaml:"sambaserver"`
}

type NFSServer struct {
	Ip        string `yaml:"ip"`
	Datastore string `yaml:"datastore"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
}

type SambaServer struct {
	Ip        string `yaml:"ip"`
	Datastore string `yaml:"datastore"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
}

type Command struct {
	Cmd string `yaml:"cmd"`
	Opt string `yaml:"opt"`
}

func (os *LinuxOS) Conn() *ssh.Client {
	sshclient, _ := utils.SSHClient(
		os.Ip,
		os.Conntp.SshCn.Username,
		os.Conntp.SshCn.Password,
		os.Conntp.SshCn.Keypath,
		os.Conntp.SshCn.Port,
	)
	return sshclient
}

func (los *LinuxOS) NFSMount(srv Server, nfssrv *NFSServer) {
	cmds := []string{}

	cmdmkdir := fmt.Sprintf(
		"mkdir -p %s",
		los.Bcktsks.Projbase,
	)
	cmdmount := fmt.Sprintf(
		"mount -t nfs %s:%s %s",
		nfssrv.Ip,
		nfssrv.Datastore,
		los.Bcktsks.Projbase,
	)
	cmds = append(cmds, cmdmkdir, cmdmount)
	utils.SSHRunCMDS(srv.Conn(), cmds)
}

func (los *LinuxOS) NFSBackup(srv Server) {
	cmds := []string{}

	for task, loc := range los.Bcktsks.Tasks {
		cmdmkdir := fmt.Sprintf(
			"mkdir -p %s/%s/%s/%s",
			los.Bcktsks.Projbase,
			los.Ip,
			task,
			fmt.Sprintf(time.Now().Format("20060102_150405")),
		)
		cmdbackup := fmt.Sprintf(
			"cp -rL %s %s/%s/%s/%s",
			loc,
			los.Bcktsks.Projbase,
			los.Ip,
			task,
			fmt.Sprintf(time.Now().Format("20060102_150405")),
		)
		cmds = append(cmds, cmdmkdir, cmdbackup)
	}
	utils.SSHRunCMDS(srv.Conn(), cmds)
}

func ReadBKServerYaml(path string) (*BKServerConfigs, error) {
	conf := &BKServerConfigs{}
	if f, err := os.Open(path); err != nil {
		return nil, err
	} else {
		yaml.NewDecoder(f).Decode(conf)
	}
	return conf, nil
}

func ReadServerYaml(path string) (*ServerConfigs, error) {
	conf := &ServerConfigs{}
	if f, err := os.Open(path); err != nil {
		return nil, err
	} else {
		yaml.NewDecoder(f).Decode(conf)
	}
	return conf, nil
}
