package main

import (
	"log"
	"net/http"
	"time"
)

var (
	DbManager    POIDBManager
	RedisManager POIRedisManager
	WsManager    POIWSManager
	Ticker       *time.Ticker
)

const (
	APP_ID     = "fyug6fiiadinzpha6nnlaajo22kam8rhba28oc9n86girasu"
	APP_KEY    = "r8pjshqr1edfvsgi0m17pq64j86pru7buae5bcw5f8yjxxbq"
	MASTER_KEY = "7e5nby4ljia5sqei97v5efvelf1a5cgplkasubm1q3gugs9u"

	DB_URL = "poi:public11223@tcp(poianalytics.mysql.rds.aliyuncs.com:3306)/wolai"

	REDIS_HOST     = "121.41.108.66:"
	REDIS_PORT     = "6381"
	REDIS_DB       = 0
	REDIS_PASSWORD = ""

/*
	APP_ID     = "h8tqr42vy41bsameqmxf5c23x2g5hbzh9t7qmhzjehk8ully"
	APP_KEY    = "6samy9ruqg9u1pfvlcgyophmcyvo8l45ytlwx68z6bs5x2hx"
	MASTER_KEY = "l8fukae8disscv46z1810dorna1z6ugh2pjv9edkribbfbhq"

	DB_URL = "webservice:Wolai11223@tcp(rdse58d3a61484hx8c69.mysql.rds.aliyuncs.com:3306)/wolai"

	REDIS_HOST     = "10.171.232.244:"
	REDIS_PORT     = "6379"
	REDIS_DB       = 0
	REDIS_PASSWORD = "Poi11223"
*/
)

func init() {
	DbManager = NewPOIDBManager()
	RedisManager = NewPOIRedisManager()
	WsManager = NewPOIWSManager()
	Ticker = time.NewTicker(time.Millisecond * 5000)
}

func main() {
	go POIOrderHandler()
	go POISessionHandler()
	go POISessionTickerHandler()

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
