package main

import (
	"github.com/cihub/seelog"
	"gopkg.in/alecthomas/kingpin.v2"
	"syncer"
)

func main() {
	kingpin.Parse()

	if err := syncer.Init("conf/syncer.json"); err != nil {
		seelog.Errorf("init failed.")
		return
	}

	syncer.MysqlSyncerInstance.Dump()
}
