package websocket

import (
	"encoding/json"
	"io"
	"math"
	"net/http"
	"strconv"
	"time"

	"POIWolaiWebService/managers"
	"POIWolaiWebService/models"

	seelog "github.com/cihub/seelog"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func V1WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// 将HTTP请求升级为Websocket连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		seelog.Error("V1WebSocketHandler:", err.Error())
		return
	}
	defer func() {
		conn.Close()
		seelog.Debug("V1WebSocketHandler close websocket connection ......")
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()
	// 读取Websocket初始化消息
	_, p, err := conn.ReadMessage()
	if err != nil {
		seelog.Error("V1WebSocketHandler:", err.Error())
		//暂时注释掉，用于debug websocket
		//		return
	}

	// 消息反序列化
	var msg models.POIWSMessage
	err = json.Unmarshal([]byte(p), &msg)
	if err != nil {
		// Force quit the user if msg is unstructed
		resp := models.NewPOIWSMessage(msg.MessageId, msg.UserId, models.WS_FORCE_QUIT)
		resp.Attribute["errCode"] = "2"
		resp.Attribute["errMsg"] = "unstructed message"
		err = conn.WriteJSON(resp)
		//		conn.Close()
		seelog.Debug("V1WSHandler: unstructed message")
		return
	}

	seelog.Debug("V1WSHandler: recieved: ", string(p))

	// 比对客户端时间和系统时间
	timestamp := time.Now().Unix()
	if math.Abs(msg.Timestamp-float64(timestamp)) > 12*3600 {
		// Force quit the user if timestamp difference is too significant
		resp := models.NewPOIWSMessage(msg.MessageId, msg.UserId, models.WS_FORCE_QUIT)
		resp.Attribute["errCode"] = "3"
		resp.Attribute["errMsg"] = "local time not accepted"
		err = conn.WriteJSON(resp)
		//		conn.Close()
		seelog.Debug("V1WSHandler: User local time not accepted; UserId: ", msg.UserId)
		return
	}

	// 利用初始化信息登录用户
	userChan, ok := WSUserLogin(msg)
	if !ok {
		// Force quit illegal login
		resp := models.NewPOIWSMessage(msg.MessageId, msg.UserId, models.WS_FORCE_QUIT)
		resp.Attribute["errCode"] = "4"
		resp.Attribute["errMsg"] = "illegal websocket login"
		err = conn.WriteJSON(resp)
		//		conn.Close()
		seelog.Debug("V1WSHandler: illegal websocket login; UserId: ", msg.UserId)
		return
	} else {
		loginResp := models.NewPOIWSMessage(msg.MessageId, msg.UserId, msg.OperationCode+1)
		err = conn.WriteJSON(loginResp)
		if err != nil {
			seelog.Debug("V1WSHandler:Send Code 12 to ", msg.UserId)
		}
	}

	// 建立处理用户连接的独立goroutine
	userId := msg.UserId
	user := models.QueryUserById(userId)
	go WebSocketWriteHandler(conn, userId, userChan)

	// 恢复可能存在的用户被中断的发单请求
	go RecoverStudentOrder(userId)
	go RecoverUserSession(userId)

	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		seelog.Debug("@@@@@@@@@@@@@@@Recieve pong message:", msg.UserId)
		err := conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			seelog.Error(err.Error())
			return err
		}
		return nil
	})
	loginTS := managers.WsManager.GetUserOnlineStatus(userId)
	for {
		// 读取Websocket信息
		_, p, err = conn.ReadMessage()
		if err != nil {
			errMsg := err.Error()
			seelog.Debug("WebSocketWriteHandler: user timed out; UserId: ", userId, " ErrorInfo:", errMsg)
			errUnexpectedEOF := &websocket.CloseError{Code: websocket.CloseAbnormalClosure, Text: io.ErrUnexpectedEOF.Error()}
			if errMsg != errUnexpectedEOF.Error() {
				if managers.WsManager.GetUserOnlineStatus(userId) == loginTS {
					WSUserLogout(userId)
					close(userChan)
				}
				return
			}
		}

		// 信息反序列化
		err = json.Unmarshal([]byte(p), &msg)
		if err != nil {
			seelog.Error("V1WSHandler:", err.Error())
			seelog.Debug("V1WSHandler recieved: UserId", msg.UserId, "Msg: ", string(p))
			seelog.Debug("V1WSHandler: unstructed message")
			continue
		}

		// 如果信息与本连接的用户Id不符合，忽略信息
		if msg.UserId != userId {
			continue
		}

		if msg.OperationCode != models.WS_PONG {
			seelog.Debug("V1Handler websocket recieve message:", string(p))
		}

		// 比对客户端时间和系统时间
		timestamp = time.Now().Unix()
		if math.Abs(msg.Timestamp-float64(timestamp)) > 12*3600 {
			// Force quit the user if timestamp difference is too significant
			seelog.Debug("V1WSHandler: User local time not accepted; UserId: ", msg.UserId)
			resp := models.NewPOIWSMessage(msg.MessageId, msg.UserId, models.WS_FORCE_QUIT)
			resp.Attribute["errCode"] = "3"
			resp.Attribute["errMsg"] = "local time not accepted"
			userChan <- resp
			return
		}

		// 根据信息中的操作码进行对应处理
		switch msg.OperationCode {

		// 用户登出信息
		case models.WS_LOGOUT:
			seelog.Debug("User:", userId, " COMMON LOGOUT!!!!!")
			resp := models.NewPOIWSMessage("", userId, models.WS_LOGOUT_RESP)
			userChan <- resp
			WSUserLogout(userId)
			managers.RedisManager.RemoveUserObjectId(userId)
			close(userChan)

		// 订单中心老师上线信息
		case models.WS_ORDER_TEACHER_ONLINE:
			resp := models.NewPOIWSMessage(msg.MessageId, userId, models.WS_ORDER_TEACHER_RESP)
			if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER {
				managers.WsManager.SetTeacherOnline(userId, timestamp)
				go RecoverTeacherOrder(userId)
				resp.Attribute["errCode"] = "0"
			} else {
				resp.Attribute["errCode"] = "2"
				resp.Attribute["errMsg"] = "You are not a teacher"
			}
			userChan <- resp

		// 订单中心老师下线信息
		case models.WS_ORDER_TEACHER_OFFLINE:
			resp := models.NewPOIWSMessage(msg.MessageId, userId, models.WS_ORDER_TEACHER_OFFLINE_RESP)
			if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER {
				managers.WsManager.SetTeacherOffline(userId)
				resp.Attribute["errCode"] = "0"
			} else {
				resp.Attribute["errCode"] = "2"
				resp.Attribute["errMsg"] = "You are not a teacher"
			}
			userChan <- resp

		// 创建发单请求信息
		case models.WS_ORDER_CREATE:
			resp := models.NewPOIWSMessage(msg.MessageId, userId, models.WS_ORDER_CREATE_RESP)
			if InitOrderDispatch(msg, userId, timestamp) {
				resp.Attribute["errCode"] = "0"
				resp.Attribute["countdown"] = "120"
			} else {
				resp.Attribute["errCode"] = "2"
				resp.Attribute["errMsg"] = "Error on order creation"
			}
			userChan <- resp

		// 订单相关信息，直接转发处理
		case models.WS_ORDER_REPLY,
			models.WS_ORDER_CONFIRM,
			models.WS_ORDER_CANCEL:
			resp := models.NewPOIWSMessage(msg.MessageId, userId, msg.OperationCode+1)

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

			if !managers.WsManager.HasOrderChan(orderId) {
				break
			}
			orderChan := managers.WsManager.GetOrderChan(orderId)
			orderChan <- msg

		// 上课相关信息，直接转发处理
		case models.WS_SESSION_START,
			models.WS_SESSION_ACCEPT,
			models.WS_SESSION_PAUSE,
			models.WS_SESSION_RESUME,
			models.WS_SESSION_FINISH,
			models.WS_SESSION_CANCEL,
			models.WS_SESSION_RESUME_ACCEPT,
			models.WS_SESSION_RESUME_CANCEL:
			resp := models.NewPOIWSMessage(msg.MessageId, userId, msg.OperationCode+1)

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

			if !managers.WsManager.HasSessionChan(sessionId) {
				break
			}
			sessionChan := managers.WsManager.GetSessionChan(sessionId)
			sessionChan <- msg

		}
	}
}

func WebSocketWriteHandler(conn *websocket.Conn, userId int64, userChan chan models.POIWSMessage) {
	// 初始化心跳计时器
	pingTicker := time.NewTicker(pingPeriod)
	defer func() {
		pingTicker.Stop()
		conn.Close()
		seelog.Debug("WebSocketWriteHandler close websocket connection")
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()

	loginTS := managers.WsManager.GetUserOnlineStatus(userId)

	for {
		select {
		// 发送心跳
		case <-pingTicker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				seelog.Error("WebSocket Write Error: UserId", userId, "ErrMsg: ", err.Error())
				if managers.WsManager.GetUserOnlineStatus(userId) == loginTS {
					WSUserLogout(userId)
					close(userChan)
				}
				return
			}
			seelog.Debug("&&&&&&&&&&&&Send ping message:", userId)
		// 处理向用户发送消息
		case msg, ok := <-userChan:
			if ok {
				conn.SetWriteDeadline(time.Now().Add(writeWait))
				err := conn.WriteJSON(msg)
				if err != nil {
					seelog.Error("WebSocket Write Error: UserId", userId, "ErrMsg: ", err.Error())
					if managers.WsManager.GetUserOnlineStatus(userId) == loginTS {
						WSUserLogout(userId)
						close(userChan)
					}
					return
				}

				msgByte, err := json.Marshal(msg)
				if err == nil {
					seelog.Debug("WebSocketWriter: UserId: ", userId, "Msg: ", string(msgByte))
				}

				if msg.OperationCode == models.WS_FORCE_QUIT ||
					msg.OperationCode == models.WS_FORCE_LOGOUT ||
					msg.OperationCode == models.WS_LOGOUT_RESP {
					if managers.WsManager.GetUserOnlineStatus(userId) == loginTS {
						WSUserLogout(userId)
						close(userChan)
						seelog.Debug("quit or logout .....")
					}
					return
				}
			} else {
				return
			}
		}
	}
}
