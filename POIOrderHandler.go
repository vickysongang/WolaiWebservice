package main

import (
	"fmt"
)

func POIOrderHandler() {
	var msg POIWSMessage
	for {
		select {
		case msg = <-WsManager.OrderInput:
			fmt.Println("WSHandler recieve: ", msg.MessageId)
			userChan := WsManager.GetUserChan(msg.UserId)
			userChan <- NewType2Message()
		}
	}
}
