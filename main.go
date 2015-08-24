package main

import (
	"net/http"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/astaxie/beego/orm"
	seelog "github.com/cihub/seelog"
)

var (
	RedisManager    POIRedisManager
	WsManager       POIWSManager
	SessionTicker   *time.Ticker
	LCMessageTicker *time.Ticker
	Config          POIConfig
)

func init() {
	//加载seelog的配置文件，使用配置文件里的方式输出日志信息
	logger, err := seelog.LoggerFromConfigAsFile("/var/lib/poi/logs/config/seelog.xml")
	if err != nil {
		panic(err)
	}
	seelog.ReplaceLogger(logger)

	//加载系统使用到的配置信息
	if _, err := toml.DecodeFile("/var/lib/poi/POIWolaiWebService.toml", &Config); err != nil {
		seelog.Critical(err.Error())
	}

	RedisManager = NewPOIRedisManager()
	WsManager = NewPOIWSManager()
	SessionTicker = time.NewTicker(time.Millisecond * 5000)
	LCMessageTicker = time.NewTicker(time.Minute * 1)
	err = orm.RegisterDataBase("default", "mysql", Config.Database.Username+":"+
		Config.Database.Password+"@"+
		Config.Database.Method+"("+
		Config.Database.Address+":"+
		Config.Database.Port+")/"+
		Config.Database.Database+"?charset=utf8&loc=Asia%2FShanghai", 30)
	if err != nil {
		seelog.Critical(err.Error())
	}
}

func main() {
	orm.Debug = false

	go POISessionTickerHandler()
	go POILeanCloudTickerHandler()

	router := NewRouter()
	seelog.Critical(http.ListenAndServe(Config.Server.Port, router))
}
