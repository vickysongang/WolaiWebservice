package main

import (
	"log"
	"net/http"
)

var (
	DbManager    POIDBManager
	RedisManager POIRedisManager
)

func init() {
	DbManager = NewPOIDBManager()
	RedisManager = NewPOIRedisManager()
}

func main() {
	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
