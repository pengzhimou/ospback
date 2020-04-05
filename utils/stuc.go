package utils

type Server struct {
	ip       string
	username string
	password string
	conntype string
}

type LinuxServer struct {
	Server
	shell  string
	python string
}

type WindowsServer struct {
	Server
}

type BackupServer struct {
}

type Command struct {
	cmd string
	opt string
}

type BackLoc struct {
	sourceloc []string
	destloc   string
}
