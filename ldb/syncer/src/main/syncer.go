package main

import (
	"fmt"
	"github.com/cihub/seelog"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"os/signal"
	"syncer"
	"syscall"
)

func WaitSignalToStop() {
	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	s := <-stop
	seelog.Infof("get signal[%v]. syncer will exit", s)
}

func main() {
	kingpin.Parse()

	if err := syncer.Init("conf/syncer.json"); err != nil {
		seelog.Errorf("syncer init failed.")
		fmt.Println(err)
		return
	}

	if err := syncer.MysqlSyncerInstance.Dump(); err != nil {
		seelog.Errorf("syncer dump failed. err[%v]", err)
		fmt.Println(err)
	} else {
		WaitSignalToStop()
	}
}
