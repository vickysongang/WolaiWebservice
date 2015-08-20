package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func V1WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// 将HTTP请求升级为Websocket连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// 读取Websocket初始化消息
	_, p, err := conn.ReadMessage()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("V1WSHandler: recieved: ", string(p))

	// 消息反序列化
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

	// 比对客户端时间和系统时间
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

	// 利用初始化信息登录用户
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

	// 建立处理用户连接的独立goroutine
	userId := msg.UserId
	user := QueryUserById(userId)
	go WebSocketWriteHandler(conn, userId, userChan)

	// 恢复可能存在的用户被中断的发单请求
	go RecoverStudentOrder(userId)

	for {

		// 读取Websocket信息
		_, p, err = conn.ReadMessage()
		if err != nil {
			return
		}

		// 信息反序列化
		err = json.Unmarshal([]byte(p), &msg)
		if err != nil {
			fmt.Println("V1WSHandler: recieved: ", string(p))
			fmt.Println("V1WSHandler: unstructed message")
			continue
		}

		// 如果信息与本连接的用户Id不符合，忽略信息
		if msg.UserId != userId {
			continue
		}

		// 比对客户端时间和系统时间
		timestamp = time.Now().Unix()
		if math.Abs(msg.Timestamp-float64(timestamp)) > 12*3600 {
			// Force quit the user if timestamp difference is too significant
			fmt.Println("V1WSHandler: User local time not accepted; UserId: ", msg.UserId)

			resp := NewPOIWSMessage(msg.MessageId, msg.UserId, WS_FORCE_QUIT)
			resp.Attribute["errMsg"] = "local time not accepted"
			userChan <- resp
			return
		}

		// 输出
		if msg.OperationCode != WS_PONG {
			fmt.Println("V1WSHandler: recieved: ", string(p))
		}

		// 根据信息中的操作码进行对应处理
		switch msg.OperationCode {

		// 心跳信息，直接转发处理
		case WS_PONG:
			userChan <- msg

		// 用户登出信息
		case WS_LOGOUT:
			resp := NewPOIWSMessage("", userId, WS_LOGOUT_RESP)
			userChan <- resp

		// 订单中心老师上线信息
		case WS_ORDER_TEACHER_ONLINE:
			resp := NewPOIWSMessage(msg.MessageId, userId, WS_ORDER_TEACHER_RESP)
			if user.AccessRight == 2 {
				WsManager.SetTeacherOnline(userId, timestamp)
				go RecoverTeacherOrder(userId)
				resp.Attribute["errCode"] = "0"
			} else {
				resp.Attribute["errCode"] = "2"
				resp.Attribute["errMsg"] = "You are not a teacher"
			}
			userChan <- resp

		// 订单中心老师下线信息
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

		// 创建发单请求信息
		case WS_ORDER_CREATE:
			resp := NewPOIWSMessage(msg.MessageId, userId, WS_ORDER_CREATE_RESP)
			if InitOrderDispatch(msg, userId, timestamp) {
				resp.Attribute["errCode"] = "0"
			} else {
				resp.Attribute["errCode"] = "2"
				resp.Attribute["errMsg"] = "Error on order creation"
			}
			userChan <- resp

		// 订单相关信息，直接转发处理
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

		// 上课相关信息，直接转发处理
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
	}
}

func WebSocketWriteHandler(conn *websocket.Conn, userId int64, userChan chan POIWSMessage) {

	// 初始化心跳计时器
	pingTicker := time.NewTicker(time.Second * 15)
	pongTicker := time.NewTicker(time.Second * 15)
	pingpong := true

	for {
		select {

		// 发送心跳
		case <-pingTicker.C:
			pingMsg := NewPOIWSMessage("", userId, WS_PING)
			err := conn.WriteJSON(pingMsg)
			if err != nil {
				fmt.Println("WebSocket Write Error: UserId", userId, "ErrMsg: ", err.Error())

				WSUserLogout(userId)
				close(userChan)
				conn.Close()
				return
			}

		// 检验用户是否连接超时
		case <-pongTicker.C:
			fmt.Println("HEARTBEAT: UserId: ", userId, "pingpong bool: ", pingpong)
			if pingpong {
				pingpong = false
			} else {
				fmt.Println("WebSocketWriteHandler: user timed out; UserId: ", userId)
				WSUserLogout(userId)
				close(userChan)
				conn.Close()
				return
			}

		// 处理向用户发送消息
		case msg := <-userChan:
			// 特殊处理，收到用户心跳信息
			if msg.OperationCode == WS_PONG {
				pingpong = true
			} else {
				err := conn.WriteJSON(msg)
				if err != nil {
					fmt.Println("WebSocket Write Error: UserId", userId, "ErrMsg: ", err.Error())

					WSUserLogout(userId)
					close(userChan)
					conn.Close()
					return
				}

				if msg.OperationCode == WS_FORCE_QUIT ||
					msg.OperationCode == WS_FORCE_LOGOUT ||
					msg.OperationCode == WS_LOGOUT_RESP {

					WSUserLogout(userId)
					close(userChan)
					conn.Close()
					return
				}
			}
		}
	}
}
