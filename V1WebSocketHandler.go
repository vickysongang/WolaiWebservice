package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func V1WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	_, p, err := conn.ReadMessage()
	if err != nil {
		return
	}

	var msg POIWSMessage

	fmt.Println("WebSocketHandler: recieved: ", string(p))
	err = json.Unmarshal([]byte(p), &msg)
	if err != nil {
		fmt.Println("WebSocketHandler: unstructed message")
		return
	}

	userChan := make(chan POIWSMessage)
	WsManager.SetUserChan(msg.UserId, userChan)
	go WebSocketWriteHandler(conn, userChan)

	for {
		_, p, err = conn.ReadMessage()
		if err != nil {
			return
		}

		fmt.Println("WSSocket recieved: ", string(p))
		err = json.Unmarshal([]byte(p), &msg)

		WsManager.OrderInput <- msg
	}
}

func WebSocketWriteHandler(conn *websocket.Conn, userChan chan POIWSMessage) {
	for {
		select {
		case msg := <-userChan:
			if msg.OperationCode == -1 {
				return
			}
			err := conn.WriteJSON(msg)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}
