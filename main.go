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

func init() {
	DbManager = NewPOIDBManager()
	RedisManager = NewPOIRedisManager()
	WsManager = NewPOIWSManager()
}

func main() {
	go POIOrderHandler()

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
