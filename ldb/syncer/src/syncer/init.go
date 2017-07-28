package syncer

import (
	"fmt"
	"github.com/cihub/seelog"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

var MysqlSyncerInstance *MysqlSyncer
var woklog seelog.LoggerInterface
var satlog seelog.LoggerInterface
var acclog seelog.LoggerInterface
var config *viper.Viper

func Init(confpath string) error {
	confdir := filepath.Dir(confpath)
	confbase := filepath.Base(confpath)
	nametype := strings.Split(confbase, ".")
	if len(nametype) != 2 {
		return seelog.Warnf("confbase has no suffix type. confbase: [%v]", confbase)
	}
	confname := nametype[0]
	conftype := nametype[1]

	{
		var err error
		if woklog, err = seelog.LoggerFromConfigAsFile(fmt.Sprintf("%s/syncer_woklog.xml", confdir)); err != nil {
			return seelog.Errorf("init hlog failed. err: [%v]", err)
		}
		if satlog, err = seelog.LoggerFromConfigAsFile(fmt.Sprintf("%s/syncer_satlog.xml", confdir)); err != nil {
			return seelog.Errorf("init hlog failed. err: [%v]", err)
		}
		if acclog, err = seelog.LoggerFromConfigAsFile(fmt.Sprintf("%s/syncer_acclog.xml", confdir)); err != nil {
			return seelog.Errorf("init hlog failed. err: [%v]", err)
		}
		seelog.ReplaceLogger(woklog)
	}

	{
		viper.AddConfigPath(confdir)
		viper.SetConfigName(confname)
		viper.SetConfigType(conftype)
		if err := config.ReadInConfig(); err != nil {
			return seelog.Errorf("init config failed. err: [%v]", err)
		}
	}

	{
		MysqlSyncerInstance = &MysqlSyncer{}
		err := MysqlSyncerInstance.Init(
			uint32(viper.GetInt("mysql.server_id")),
			viper.GetString("mysql.flavor"),
			viper.GetString("mysql.host"),
			uint16(viper.GetInt("mysql.port")),
			viper.GetString("mysql.user"),
			viper.GetString("mysql.password"),
		)
		if err != nil {
			return seelog.Errorf("init syncer failed.", err)
		}
	}

	return nil
}
