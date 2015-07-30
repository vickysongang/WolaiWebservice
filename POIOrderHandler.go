package main

import (
	"fmt"
)

func POIOrderHandler() {
	var msg POIWSMessage
	for {
		select {
		case msg = <-WsManager.OrderInput:
			fmt.Printf("WSHandler recieve: ", msg.MessageId)
			//WsManager.UserMap[msg.UserId] <- NewType2Message
		}
	}
}
