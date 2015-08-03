package main

import (
	"log"
	"net/http"
)

var (
	DbManager    POIDBManager
	RedisManager POIRedisManager
	WsManager    POIWSManager
)

const (
	APP_ID     = "fyug6fiiadinzpha6nnlaajo22kam8rhba28oc9n86girasu"
	APP_KEY    = "r8pjshqr1edfvsgi0m17pq64j86pru7buae5bcw5f8yjxxbq"
	MASTER_KEY = "7e5nby4ljia5sqei97v5efvelf1a5cgplkasubm1q3gugs9u"
)

func init() {
	DbManager = NewPOIDBManager()
	RedisManager = NewPOIRedisManager()
	WsManager = NewPOIWSManager()
}

func main() {
	go POIOrderHandler()
	go POISessionHandler()

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
