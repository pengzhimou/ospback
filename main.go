package main

import (
	"ospback/utils"
)

func main() {
	clt, _ := utils.SSHClient(
		"192.168.122.21",
		"root",
		"redhat",
		"",
		22)

	var cmds = []string{"uname -a", "ls /"}
	utils.SSHRunCMDS(clt, cmds)

	//utils.SSHRunCMDS2(clt, cmds)

}
