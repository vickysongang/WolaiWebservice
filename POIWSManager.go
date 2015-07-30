package main

import (
	_ "fmt"
)

type POIWSManager struct {
	OrderInput chan POIWSMessage
	UserMap    map[int64](chan POIWSMessage)
}

func NewPOIWSManager() POIWSManager {
	return POIWSManager{
		OrderInput: make(chan POIWSMessage),
		UserMap:    make(map[int64](chan POIWSMessage)),
	}
}
