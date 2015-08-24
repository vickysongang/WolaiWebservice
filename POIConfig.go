package main

import (
	"time"

	seelog "github.com/cihub/seelog"
)

type POIConfig struct {
	Title     string
	Server    serverConf
	Database  databaseConf
	Redis     redisConf
	LeanCloud leancloudConf
	Reminder  reminderConf
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
	}
	return err
}
