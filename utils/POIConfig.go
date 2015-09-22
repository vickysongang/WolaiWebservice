package utils

import (
	"time"

	"github.com/BurntSushi/toml"
	seelog "github.com/cihub/seelog"
)

const DB_TYPE = "mysql"

type POIConfig struct {
	Title     string
	Server    serverConf
	Database  databaseConf
	Redis     redisConf
	LeanCloud leancloudConf
	Reminder  reminderConf
}

var Config POIConfig

func init() {
	//加载系统使用到的配置信息
	if _, err := toml.DecodeFile("/var/lib/poi/POIWolaiWebService.toml", &Config); err != nil {
		seelog.Critical(err.Error())
	}
}

type serverConf struct {
	Port string
}

type databaseConf struct {
	Username string
	Password string
	Method   string
	Address  string
	Port     string
	Database string
}

type redisConf struct {
	Host     string
	Port     string
	Db       int64
	Password string
}

type leancloudConf struct {
	AppId     string
	AppKey    string
	MasterKey string
}

type reminderConf struct {
	Durations []duration
}

// Parse duration
type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	if err != nil {
		seelog.Error(string(text), " ", err.Error())
		return err
	}
	return nil
}