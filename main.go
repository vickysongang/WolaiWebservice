package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/astaxie/beego/orm"
)

var (
	RedisManager    POIRedisManager
	WsManager       POIWSManager
	Ticker          *time.Ticker
	LCMessageTicker *time.Ticker
	Config          POIConfig
)

func init() {
	if _, err := toml.DecodeFile("/var/lib/poi/POIWolaiWebService.toml", &Config); err != nil {
		fmt.Println(err.Error())
	}
	RedisManager = NewPOIRedisManager()
	WsManager = NewPOIWSManager()
	Ticker = time.NewTicker(time.Millisecond * 5000)
	LCMessageTicker = time.NewTicker(time.Minute * 1)
	orm.RegisterDataBase("default", "mysql", Config.Database.Username+":"+
		Config.Database.Password+"@"+
		Config.Database.Method+"("+
		Config.Database.Address+":"+
		Config.Database.Port+")/"+
		Config.Database.Database+"?charset=utf8&loc=Asia%2FShanghai", 30)
}

func main() {
	orm.Debug = false
	go POISessionTickerHandler()
	go POILeanCloudTickerHandler()
	router := NewRouter()
	log.Fatal(http.ListenAndServe(Config.Server.Port, router))
}
