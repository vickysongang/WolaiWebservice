package main

import (
	"encoding/json"
	"fmt"
	"net/http"
<<<<<<< HEAD
	"strconv"
	"time"
=======
>>>>>>> beegoorm

	"github.com/gorilla/websocket"
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

<<<<<<< HEAD
	timestamp := time.Now().Unix()
	if math.Abs(msg.Timestamp-float64(timestamp)) > 12*3600 {
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
	user := DbManager.QueryUserById(userId)
	go WebSocketWriteHandler(conn, userId, userChan)
=======
	userChan := make(chan POIWSMessage)
	WsManager.SetUserChan(msg.UserId, userChan)
	go WebSocketWriteHandler(conn, userChan)

	WsManager.OrderInput <- msg
>>>>>>> beegoorm

	for {
		_, p, err = conn.ReadMessage()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println("V1WSOrderHandler recieved: ", string(p))
		err = json.Unmarshal([]byte(p), &msg)

<<<<<<< HEAD
		if msg.UserId != userId {
			continue
		}

		timestamp = time.Now().Unix()
		if math.Abs(msg.Timestamp-float64(timestamp)) > 12*3600 {
			// Force quit the user if timestamp difference is too significant
			resp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_FORCE_QUIT)
			resp.Attribute["errMsg"] = "local time not accepted"
			err = conn.WriteJSON(resp)
			conn.Close()

			fmt.Println("V1WSHandler: User local time not accepted; UserId: ", msg.UserId)
			return
		}

		if msg.OperationCode != WS_PONG {
			fmt.Println("V1WSHandler: recieved: ", string(p))
		}

		switch msg.OperationCode {
		case WS_PONG:
			userChan <- msg

		case WS_LOGOUT:
			_, _ = WSUserLogout(msg.UserId)
			resp := NewPOIWSMessage("", userId, WS_LOGOUT_RESP)
			userChan <- resp

		case WS_ORDER_TEACHER_ONLINE:
			resp := NewPOIWSMessage(msg.MessageId, userId, WS_ORDER_TEACHER_RESP)
			if user.AccessRight == 2 {
				WsManager.SetTeacherOnline(userId, timestamp)
				resp.Attribute["errCode"] = "0"
			} else {
				resp.Attribute["errCode"] = "2"
				resp.Attribute["errMsg"] = "You are not a teacher"
			}
			userChan <- resp

		case WS_ORDER_TEACHER_OFFLINE:
			resp := NewPOIWSMessage(msg.MessageId, userId, WS_ORDER_TEACHER_OFFLINE_RESP)
			if user.AccessRight == 2 {
				WsManager.SetTeacherOffline(userId)
				resp.Attribute["errCode"] = "0"
			} else {
				resp.Attribute["errCode"] = "2"
				resp.Attribute["errMsg"] = "You are not a teacher"
			}
			userChan <- resp

		case WS_ORDER_CREATE:
			resp := NewPOIWSMessage(msg.MessageId, userId, WS_ORDER_CREATE_RESP)
			if InitOrderDispatch(msg, userId, timestamp) {
				resp.Attribute["errCode"] = "0"
			} else {
				resp.Attribute["errCode"] = "2"
				resp.Attribute["errMsg"] = "Error on order creation"
			}
			userChan <- resp

		case WS_ORDER_REPLY,
			WS_ORDER_CONFIRM,
			WS_ORDER_CANCEL:
			resp := NewPOIWSMessage(msg.MessageId, userId, msg.OperationCode+1)

			orderIdStr, ok := msg.Attribute["orderId"]
			if !ok {
				resp.Attribute["errCode"] = "2"
				userChan <- resp
				break
			}

			orderId, err := strconv.ParseInt(orderIdStr, 10, 64)
			if err != nil {
				resp.Attribute["errCode"] = "2"
				userChan <- resp
				break
			}

			if !WsManager.HasOrderChan(orderId) {
				break
			}
			orderChan := WsManager.GetOrderChan(orderId)
			orderChan <- msg

		case WS_SESSION_START,
			WS_SESSION_ACCEPT,
			WS_SESSION_PAUSE,
			WS_SESSION_RESUME,
			WS_SESSION_FINISH,
			WS_SESSION_CANCEL:
			resp := NewPOIWSMessage(msg.MessageId, userId, msg.OperationCode+1)

			sessionIdStr, ok := msg.Attribute["sessionId"]
			if !ok {
				resp.Attribute["errCode"] = "2"
				userChan <- resp
				break
			}

			sessionId, err := strconv.ParseInt(sessionIdStr, 10, 64)
			if err != nil {
				resp.Attribute["errCode"] = "2"
				userChan <- resp
				break
			}

			if !WsManager.HasSessionChan(sessionId) {
				break
			}
			sessionChan := WsManager.GetSessionChan(sessionId)
			sessionChan <- msg

		}
=======
		WsManager.OrderInput <- msg
>>>>>>> beegoorm
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

<<<<<<< HEAD
				WsManager.SetUserOffline(userId)
				WsManager.SetTeacherOffline(userId)
				close(userChan)
				conn.Close()
				return
			}
=======
	fmt.Println("V1WSSessionHandler: recieved: ", string(p))
	err = json.Unmarshal([]byte(p), &msg)
	if err != nil {
		fmt.Println("V1WSSessionHandler: unstructed message")
		return
	}
>>>>>>> beegoorm

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
