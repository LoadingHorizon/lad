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
			return seelog.Errorf("init woklog failed. err: [%v]", err)
		}
		if satlog, err = seelog.LoggerFromConfigAsFile(fmt.Sprintf("%s/syncer_satlog.xml", confdir)); err != nil {
			return seelog.Errorf("init satlog failed. err: [%v]", err)
		}
		if acclog, err = seelog.LoggerFromConfigAsFile(fmt.Sprintf("%s/syncer_acclog.xml", confdir)); err != nil {
			return seelog.Errorf("init acclog failed. err: [%v]", err)
		}
		seelog.Infof("init logger success.")
		seelog.ReplaceLogger(woklog)
	}

	{
		config = viper.New()
		config.AddConfigPath(confdir)
		config.SetConfigName(confname)
		config.SetConfigType(conftype)
		if err := config.ReadInConfig(); err != nil {
			return seelog.Errorf("init viper failed. err: [%v]", err)
		}
		seelog.Infof("init config success.")
	}

	{
		MysqlSyncerInstance = &MysqlSyncer{}
		err := MysqlSyncerInstance.Init(
			uint32(viper.GetInt("mysql.server_id")),
			config.GetString("mysql.flavor"),
			config.GetString("mysql.host"),
			uint16(config.GetInt("mysql.port")),
			config.GetString("mysql.user"),
			config.GetString("mysql.password"),
			config.GetString("mysql.database"),
			config.GetStringSlice("mysql.tables"),
			config.GetString("syncpoint"),
		)
		if err != nil {
			return seelog.Errorf("init syncer failed.", err)
		}
		seelog.Infof("init syncer success.")
	}

	return nil
}
