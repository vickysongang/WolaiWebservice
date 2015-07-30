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

func print_binary(s []byte) {
	fmt.Printf("Received b:")
	for n := 0; n < len(s); n++ {
		fmt.Printf("%d,", s[n])
	}
	fmt.Printf("\n")
}

func V1WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}

		var msg POIWSMessage
		print_binary(p)
		fmt.Println("WSSocket recieved: ", string(p))
		err = json.Unmarshal([]byte(p), &msg)
		if err != nil {
			fmt.Println(err.Error())
		}
		WsManager.OrderInput <- msg

		err = conn.WriteMessage(messageType, p)
		if err != nil {
			return
		}
	}
}
