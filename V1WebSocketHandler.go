package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func V1WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, p, err := conn.ReadMessage()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("V1WSHandler: recieved: ", string(p))

	var msg POIWSMessage
	err = json.Unmarshal([]byte(p), &msg)
	if err != nil {
		// Force quit the user if msg is unstructed
		resp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_FORCE_QUIT)
		resp.Attribute["errMsg"] = "unstructed message"
		err = conn.WriteJSON(resp)
		conn.Close()

		fmt.Println("V1WSHandler: unstructed message")
		return
	}

	timestampNano := time.Now().UnixNano()
	timestampMillis := timestampNano / 1000
	timestamp := float64(timestampMillis) / 1000000.0
	if math.Abs(msg.Timestamp-timestamp) > 12*3600 {
		// Force quit the user if timestamp difference is too significant
		resp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_FORCE_QUIT)
		resp.Attribute["errMsg"] = "local time not accepted"
		err = conn.WriteJSON(resp)
		conn.Close()

		fmt.Println("V1WSHandler: User local time not accepted; UserId: ", msg.UserId)
		return
	}

	userChan, ok := WSUserLogin(msg)
	if !ok {
		// Force quit illegal login
		resp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_FORCE_QUIT)
		resp.Attribute["errMsg"] = "illegal websocket login"
		err = conn.WriteJSON(resp)
		conn.Close()

		fmt.Println("V1WSHandler: illegal websocket login; UserId: ", msg.UserId)
		return
	} else {
		loginResp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_LOGIN_RESP)
		err = conn.WriteJSON(loginResp)
	}

	userId := msg.UserId
	go WebSocketWriteHandler(conn, userId, userChan)

	for {
		_, p, err = conn.ReadMessage()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		err = json.Unmarshal([]byte(p), &msg)
		if err != nil {
			fmt.Println("V1WSHandler: recieved: ", string(p))
			fmt.Println("V1WSHandler: unstructed message")
			continue
		}

		if msg.UserId != userId {
			continue
		}

		if msg.OperationCode != WS_PONG {
			fmt.Println("V1WSHandler: recieved: ", string(p))
		}

		switch msg.OperationCode {
		case WS_PONG:
			userChan <- msg
		case WS_LOGOUT:
			_, _ = WSUserLogout(msg.UserId)
			logoutResp := NewPOIWSMessage("", userId, WS_LOGOUT_RESP)
			userChan <- logoutResp
		}
	}
}

func WebSocketWriteHandler(conn *websocket.Conn, userId int64, userChan chan POIWSMessage) {
	pingTicker := time.NewTicker(time.Second * 15)
	pongTicker := time.NewTicker(time.Second * 20)
	pingpong := true

	for {
		select {
		case <-pingTicker.C:
			pingMsg := NewPOIWSMessage("", userId, WS_PING)
			err := conn.WriteJSON(pingMsg)
			if err != nil {
				fmt.Println(err.Error())
			}

		case <-pongTicker.C:
			if pingpong {
				pingpong = false
			} else {
				_, _ = WSUserLogout(userId)
				fmt.Println("WebSocketWriteHandler: user timed out; UserId: ", userId)

				close(userChan)
				conn.Close()
				return
			}

		case msg := <-userChan:
			if msg.OperationCode == WS_PONG {
				pingpong = true
			} else {
				err := conn.WriteJSON(msg)
				if err != nil {
					fmt.Println(err.Error())
				}

				if msg.OperationCode == WS_FORCE_QUIT || msg.OperationCode == WS_FORCE_LOGOUT {
					close(userChan)
					conn.Close()
					return
				}
			}
		}
	}
}
