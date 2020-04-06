package main

import (
	. "ospback/cls"
)

func main() {

	var linuxserver1 Server = LinuxOS{
		OS{
			"192.168.122.21",
			"root",
			"redhat",
			"ssh",
		},
		"",
		22,
	}

	var linuxserver1_backup Backup = BackupCMD{
		ComdDataStore: Command{
			Cmd: "mount",
			Opt: "-t nfs",
		},
		ComdPrpDir: Command{
			Cmd: "mkdir",
			Opt: "-p",
		},
		ComdBK: Command{
			Cmd: "cp",
			Opt: "-rL",
		},
		Backupserver: BackupServer{
			Os: OS{
				Ip: "192.168.122.21",
			},
			Srvtype:   "nfs",        // nfs/samba
			Datastore: "/datastore", // /datastore/ip/date/dirs

		},
		Projbase: "/nfsclient",
		Project:  map[string]string{"vxflexic": "/etc", "xtremioci": "/etc"},
	}

	linuxserver1_backup.PrepDataStore(linuxserver1)

	linuxserver1_backup.PrjBackup(linuxserver1)

}
