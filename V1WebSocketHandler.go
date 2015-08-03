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

func V1WSOrderHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	_, p, err := conn.ReadMessage()
	if err != nil {
		return
	}

	var msg POIWSMessage

	fmt.Println("V1WSOrderHandler: recieved: ", string(p))
	err = json.Unmarshal([]byte(p), &msg)
	if err != nil {
		fmt.Println("V1WSOrderHandler: unstructed message")
		return
	}

	userChan := make(chan POIWSMessage)
	WsManager.SetUserChan(msg.UserId, userChan)
	fmt.Println("aaa")
	go WebSocketWriteHandler(conn, userChan)
	fmt.Println("bbb")
	WsManager.OrderInput <- msg
	fmt.Println("ccc")

	for {
		_, p, err = conn.ReadMessage()
		if err != nil {
			return
		}

		fmt.Println("V1WSOrderHandler recieved: ", string(p))
		err = json.Unmarshal([]byte(p), &msg)

		WsManager.OrderInput <- msg
	}
}

func V1WSSessionHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	_, p, err := conn.ReadMessage()
	if err != nil {
		return
	}

	var msg POIWSMessage

	fmt.Println("V1WSSessionHandler: recieved: ", string(p))
	err = json.Unmarshal([]byte(p), &msg)
	if err != nil {
		fmt.Println("V1WSSessionHandler: unstructed message")
		return
	}

	userChan := make(chan POIWSMessage)
	WsManager.SetUserChan(msg.UserId, userChan)
	go WebSocketWriteHandler(conn, userChan)
	WsManager.SessionInput <- msg

	for {
		_, p, err = conn.ReadMessage()
		if err != nil {
			return
		}

		fmt.Println("V1WSSessionHandler recieved: ", string(p))
		err = json.Unmarshal([]byte(p), &msg)

		WsManager.SessionInput <- msg
	}
}
