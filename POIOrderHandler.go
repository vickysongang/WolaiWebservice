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
			//WsManager.UserMap[msg.UserId] <- NewType2Message
		}
	}
}
