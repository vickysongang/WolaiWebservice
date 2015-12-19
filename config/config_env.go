package config

import (
	"time"

	"github.com/BurntSushi/toml"
	"github.com/cihub/seelog"
)

const (
	CONFIG_FILE_LOC = "/var/lib/poi/WolaiWebservice.toml"
)

//加载系统使用到的配置信息
func init() {
	if _, err := toml.DecodeFile(CONFIG_FILE_LOC, &Env); err != nil {
		seelog.Critical(err.Error())
	}
}

type EnvironmentConf struct {
	Title     string
	Server    serverConf
	Database  databaseConf
	Redis     redisConf
	LeanCloud leancloudConf
	Pingpp    pingppConf
	SendCloud sendcloudConf
}

var Env EnvironmentConf

type serverConf struct {
	Port    string
	RpcPort string
}

type databaseConf struct {
	Type     string
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

type pingppConf struct {
	Key   string
	AppId string
}

type sendcloudConf struct {
	SmsUser    string
	TemplateId string
	SmsKey     string
	AppKey     string
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
