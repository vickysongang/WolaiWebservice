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
	DbManager    POIDBManager
	RedisManager POIRedisManager
	WsManager    POIWSManager
	Ticker       *time.Ticker
	Config       POIConfig
)

func init() {
	if _, err := toml.DecodeFile("/var/lib/poi/POIWolaiWebService.toml", &Config); err != nil {
		fmt.Println(err.Error())
	}
	DbManager = NewPOIDBManager()
	RedisManager = NewPOIRedisManager()
	WsManager = NewPOIWSManager()
	Ticker = time.NewTicker(time.Millisecond * 5000)
	orm.RegisterDataBase("default", "mysql", Config.Database.Username+":"+
		Config.Database.Password+"@"+
		Config.Database.Method+"("+
		Config.Database.Address+":"+
		Config.Database.Port+")/"+
		Config.Database.Database,30)
}

func main() {
	 orm.Debug = true
	go POIOrderHandler()
	go POISessionHandler()
	go POISessionTickerHandler()

	router := NewRouter()
	log.Fatal(http.ListenAndServe(Config.Server.Port, router))
}
