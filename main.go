package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/BurntSushi/toml"
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
}

func main() {
	go POISessionTickerHandler()

	router := NewRouter()

	log.Fatal(http.ListenAndServe(Config.Server.Port, router))
}
