package main

import (
	"errors"
	"fmt"
	. "ospback/backup"
	"ospback/logger"
	"sync"
	"time"
)

func getKeys(m map[string]LinuxOS) []string {
	// 数组默认长度为map长度,后面append时,不需要重新申请内存和拷贝,效率较高
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
func linuxNfsBackup(gorlimit int) {
	bksrvyml, err := ReadBKServerYaml("etc/backup.yml")
	if err != nil {
		errors.New("yaml fmt not correct!")
		logger.Error("yaml bkserverconfigs part fmt not correct")
	}

	srvyml, err := ReadServerYaml("etc/backup.yml")
	if err != nil {
		errors.New("yaml fmt not correct!")
		logger.Error("yaml serverconfigs part fmt not correct")
	}

	data := srvyml.LinuxOsConfig
	keys := getKeys(data)
	gorNum := 0
	gorLimit := gorlimit
	wg := sync.WaitGroup{}

	for i := 0; i < len(keys); i++ {
		wg.Add(1)
		if gorNum < gorLimit {
			gorNum++
			go func(n int) {
				fmt.Println(data[keys[i]])
				server := data[keys[i]]
				/////////////////////
				servertmp := LinuxOS{
					Ip: server.Ip,
					Conntp: Conntype{
						SshCn: SshConn{
							server.Conntp.SshCn.Username,
							server.Conntp.SshCn.Password,
							server.Conntp.SshCn.Keypath,
							server.Conntp.SshCn.Port,
						},
						TelnetCn: TelnetConn{
							server.Conntp.TelnetCn.Username,
							server.Conntp.TelnetCn.Password,
							server.Conntp.TelnetCn.Port,
						},
					},
					Bcktsks: BackupTasks{
						Projbase: server.Bcktsks.Projbase,
						//Tasks:    map[string]string{"vxflexcu": "/etc/", "xetea": "/etc/hosts"},
						Tasks: server.Bcktsks.Tasks,
					},
				}
				backupservertemp := BKServerConfigs{
					BKSCfgs: BKServerType{
						NFSSrv: NFSServer{
							Ip:        bksrvyml.BKSCfgs.NFSSrv.Ip,
							Password:  bksrvyml.BKSCfgs.NFSSrv.Password,
							Datastore: bksrvyml.BKSCfgs.NFSSrv.Datastore,
							Username:  bksrvyml.BKSCfgs.NFSSrv.Username,
						},
						SmbSrv: SambaServer{
							Ip:        bksrvyml.BKSCfgs.SmbSrv.Ip,
							Password:  bksrvyml.BKSCfgs.SmbSrv.Password,
							Datastore: bksrvyml.BKSCfgs.SmbSrv.Datastore,
							Username:  bksrvyml.BKSCfgs.SmbSrv.Username,
						},
					},
				}
				var server_itf Server
				server_itf = &servertmp
				servertmp.Conn()
				servertmp.NFSMount(server_itf, &backupservertemp.BKSCfgs.NFSSrv)
				servertmp.NFSBackup(server_itf, keys[i])
				////////////////////////
				wg.Done()
			}(i)
			wg.Wait()
		} else {
			time.Sleep(1 * time.Second)
			gorNum = 0
			i--
			wg.Done()
		}
	}

}

//func linuxNfsBackupOrig() {
//	bksrvyml, err := ReadBKServerYaml("etc/backup.yml")
//	if err != nil {
//		errors.New("yaml fmt not correct!")
//		logger.Error("yaml bkserverconfigs part fmt not correct")
//	}
//
//	srvyml, err := ReadServerYaml("etc/backup.yml")
//	if err != nil {
//		errors.New("yaml fmt not correct!")
//		logger.Error("yaml serverconfigs part fmt not correct")
//	}
//
//	wg := sync.WaitGroup{}
//
//	for prj, server := range srvyml.LinuxOsConfig {
//		servertmp := LinuxOS{
//			Ip: server.Ip,
//			Conntp: Conntype{
//				SshCn: SshConn{
//					server.Conntp.SshCn.Username,
//					server.Conntp.SshCn.Password,
//					server.Conntp.SshCn.Keypath,
//					server.Conntp.SshCn.Port,
//				},
//				TelnetCn: TelnetConn{
//					server.Conntp.TelnetCn.Username,
//					server.Conntp.TelnetCn.Password,
//					server.Conntp.TelnetCn.Port,
//				},
//			},
//			Bcktsks: BackupTasks{
//				Projbase: server.Bcktsks.Projbase,
//				//Tasks:    map[string]string{"vxflexcu": "/etc/", "xetea": "/etc/hosts"},
//				Tasks: server.Bcktsks.Tasks,
//			},
//		}
//		backupservertemp := BKServerConfigs{
//			BKSCfgs: BKServerType{
//				NFSSrv: NFSServer{
//					Ip:        bksrvyml.BKSCfgs.NFSSrv.Ip,
//					Password:  bksrvyml.BKSCfgs.NFSSrv.Password,
//					Datastore: bksrvyml.BKSCfgs.NFSSrv.Datastore,
//					Username:  bksrvyml.BKSCfgs.NFSSrv.Username,
//				},
//				SmbSrv: SambaServer{
//					Ip:        bksrvyml.BKSCfgs.SmbSrv.Ip,
//					Password:  bksrvyml.BKSCfgs.SmbSrv.Password,
//					Datastore: bksrvyml.BKSCfgs.SmbSrv.Datastore,
//					Username:  bksrvyml.BKSCfgs.SmbSrv.Username,
//				},
//			},
//		}
//		wg.Add(1)
//		go func(prjname string) {
//			var server_itf Server
//			server_itf = &servertmp
//			servertmp.Conn()
//			servertmp.NFSMount(server_itf, &backupservertemp.BKSCfgs.NFSSrv)
//			servertmp.NFSBackup(server_itf, prjname)
//			wg.Done()
//		}(prj) //注意这里的值传递
//	}
//	wg.Wait()
//}

func linuxNfsBackupOrig() {
	bksrvyml, err := ReadBKServerYaml("etc/backup.yml")
	if err != nil {
		errors.New("yaml fmt not correct!")
		logger.Error("yaml bkserverconfigs part fmt not correct")
	}

	srvyml, err := ReadServerYaml("etc/backup.yml")
	if err != nil {
		errors.New("yaml fmt not correct!")
		logger.Error("yaml serverconfigs part fmt not correct")
	}

	wg := sync.WaitGroup{}

	for prj, server := range srvyml.LinuxOsConfig {
		wg.Add(1)
		go func(prjname string) {
			var server_itf Server
			server_itf = &server
			server.Conn()
			server.NFSMount(server_itf, &bksrvyml.BKSCfgs.NFSSrv)
			server.NFSBackup(server_itf, prjname)
			wg.Done()
		}(prj) //注意这里的值传递
	}
	wg.Wait()
}

func linuxNfsBackupOrig2() {
	bksrvyml, err := ReadBKServerYaml2("etc/backup.yml")
	if err != nil {
		errors.New("yaml fmt not correct!")
		logger.Error("yaml bkserverconfigs part fmt not correct")
	}

	srvyml, err := ReadServerYaml2("etc/backup.yml")
	if err != nil {
		errors.New("yaml fmt not correct!")
		logger.Error("yaml serverconfigs part fmt not correct")
	}

	wg := sync.WaitGroup{}

	for prj, server := range srvyml.LinuxOsConfig {
		wg.Add(1)
		go func(prjname string) {
			var server_itf Server
			server_itf = &server
			server.Conn()
			server.NFSMount(server_itf, &bksrvyml.BKSCfgs.NFSSrv)
			server.NFSBackup(server_itf, prjname)
			wg.Done()
		}(prj) //注意这里的值传递
	}
	wg.Wait()
}

func main() {
	//linuxNfsBackup(3)
	//linuxNfsBackupOrig()
	linuxNfsBackupOrig2()
}
