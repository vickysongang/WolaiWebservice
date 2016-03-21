package websocket

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"time"

	seelog "github.com/cihub/seelog"
	"github.com/gorilla/websocket"

	"WolaiWebservice/config/settings"
	"WolaiWebservice/models"
	"WolaiWebservice/redis"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func V1WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// Time allowed to read the next pong message from the peer.
	pongWaitInt := settings.WebsocketPongWait()
	pongWait := time.Duration(pongWaitInt) * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriodInt := settings.WebsocketPingPeriod()

	// 将HTTP请求升级为Websocket连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		seelog.Error("V1WebSocketHandler build connection error:", err.Error())
		return
	}

	// 读取Websocket初始化消息
	_, p, err := conn.ReadMessage()
	if err != nil {
		seelog.Error("V1WebSocketHandler init message:", err.Error())
		return
	}

	// 消息反序列化
	var msg WSMessage
	err = json.Unmarshal([]byte(p), &msg)
	if err != nil {
		// Force quit the user if msg is unstructed
		resp := NewWSMessage(msg.MessageId, msg.UserId, WS_FORCE_QUIT)
		resp.Attribute["errCode"] = "2"
		resp.Attribute["errMsg"] = "unstructed message"
		err = conn.WriteJSON(resp)
		seelog.Error("V1WSHandler: unstructed message")
		return
	}

	seelog.Trace("V1WSHandler: recieved content: ", string(p))

	defer func() {
		conn.Close()
		seelog.Debug("V1WebSocketHandler close websocket connection:", msg.UserId)
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()

	// 比对客户端时间和系统时间
	timestamp := time.Now().Unix()
	if math.Abs(msg.Timestamp-float64(timestamp)) > 12*3600 {
		// Force quit the user if timestamp difference is too significant
		resp := NewWSMessage(msg.MessageId, msg.UserId, WS_FORCE_QUIT)
		resp.Attribute["errCode"] = "3"
		resp.Attribute["errMsg"] = "local time not accepted"
		err = conn.WriteJSON(resp)
		seelog.Debug("V1WSHandler: User local time not accepted; UserId: ", msg.UserId)
		return
	}

	// 利用初始化信息登录用户
	userChan, ok := WSUserLogin(msg)
	if !ok {
		// Force quit illegal login
		resp := NewWSMessage(msg.MessageId, msg.UserId, WS_FORCE_QUIT)
		resp.Attribute["errCode"] = "4"
		resp.Attribute["errMsg"] = "illegal websocket login"
		err = conn.WriteJSON(resp)
		return
	} else {
		loginResp := NewWSMessage(msg.MessageId, msg.UserId, msg.OperationCode+1)
		if TeacherManager.IsTeacherOnline(msg.UserId) {
			loginResp.Attribute["online"] = "on"
			if TeacherManager.IsTeacherAssignOpen(msg.UserId) {
				loginResp.Attribute["assign"] = "on"
			} else {
				loginResp.Attribute["assign"] = "off"
			}
		} else {
			loginResp.Attribute["online"] = "off"
			loginResp.Attribute["assign"] = "off"
		}
		loginResp.Attribute["pingPeriod"] = strconv.FormatInt(pingPeriodInt, 10)
		err = conn.WriteJSON(loginResp)
		if err == nil {
			seelog.Trace("send login response to user:", msg.UserId, " ", loginResp)
		} else {
			seelog.Error("send login response to user ", msg.UserId, " fail")
		}
	}

	// 建立处理用户连接的独立goroutine
	userId := msg.UserId
	user, _ := models.ReadUser(userId)
	go WebSocketWriteHandler(conn, userId, userChan)

	// 恢复可能存在的用户被中断的发单请求
	recoverTeacherOrder(userId)
	recoverStudentOrder(userId)
	go RecoverUserSession(userId, msg)

	//处理心跳的pong消息
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(appData string) error {
		err := conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			seelog.Error(err.Error())
			return err
		}
		return nil
	})

	loginTS := UserManager.GetUserOnlineStatus(userId)
	for {
		// 读取Websocket信息
		_, p, err = conn.ReadMessage()
		if err != nil {
			errMsg := err.Error()
			seelog.Debug("WebSocketWriteHandler: socket disconnect; UserId: ", userId, "; ErrorInfo:", errMsg)

			if UserManager.GetUserOnlineStatus(userId) == loginTS {
				WSUserLogout(userId)
				close(userChan)
			}
			return
		}

		// 信息反序列化
		err = json.Unmarshal([]byte(p), &msg)
		if err != nil {
			seelog.Error("V1WSHandler:", err.Error())
			seelog.Debug("V1WSHandler: unstructed message")
			continue
		}

		// 如果信息与本连接的用户Id不符合，忽略信息
		if msg.UserId != userId {
			continue
		}

		if msg.OperationCode != WS_PONG {
			seelog.Trace("V1Handler websocket receive message:", string(p))
		}

		// 比对客户端时间和系统时间
		timestamp = time.Now().Unix()
		if math.Abs(msg.Timestamp-float64(timestamp)) > 12*3600 {
			// Force quit the user if timestamp difference is too significant
			seelog.Debug("V1WSHandler: User local time not accepted; UserId: ", msg.UserId)
			resp := NewWSMessage(msg.MessageId, msg.UserId, WS_FORCE_QUIT)
			resp.Attribute["errCode"] = "3"
			resp.Attribute["errMsg"] = "local time not accepted"
			userChan <- resp
			return
		}

		// 根据信息中的操作码进行对应处理
		switch msg.OperationCode {

		// 用户登出信息
		case WS_LOGOUT:
			seelog.Debug("User:", userId, " logout correctly!")
			resp := NewWSMessage("", userId, WS_LOGOUT_RESP)
			userChan <- resp
			WSUserLogout(userId)
			redis.RemoveUserObjectId(userId)
			if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER {
				TeacherManager.SetOffline(userId)
				TeacherManager.SetAssignOff(userId)
			}
			close(userChan)

		// 上课相关信息，直接转发处理
		case WS_SESSION_PAUSE,
			WS_SESSION_RESUME,
			WS_SESSION_FINISH,
			WS_SESSION_RESUME_ACCEPT,
			WS_SESSION_RESUME_CANCEL:
			resp := NewWSMessage(msg.MessageId, userId, msg.OperationCode+1)

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

			if !SessionManager.IsSessionOnline(sessionId) {
				break
			}

			sessionChan, _ := SessionManager.GetSessionChan(sessionId)
			sessionChan <- msg

		case WS_ORDER2_TEACHER_ONLINE:
			resp := NewWSMessage(msg.MessageId, userId, WS_ORDER2_TEACHER_ONLINE_RESP)
			if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER {
				resp.Attribute["errCode"] = "0"
				resp.Attribute["assign"] = "off"
			} else {
				resp.Attribute["errCode"] = "2"
				resp.Attribute["errMsg"] = "You are not a teacher"
			}
			userChan <- resp
			TeacherManager.SetOnline(userId)

		case WS_ORDER2_TEACHER_OFFLINE:
			resp := NewWSMessage(msg.MessageId, userId, WS_ORDER2_TEACHER_OFFLINE_RESP)
			if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER {
				resp.Attribute["errCode"] = "0"
			} else {
				resp.Attribute["errCode"] = "2"
				resp.Attribute["errMsg"] = "You are not a teacher"
			}
			if err := TeacherManager.SetOffline(userId); err != nil {
				resp.Attribute["errCode"] = "2"
				resp.Attribute["errMsg"] = err.Error()
			}
			userChan <- resp

		case WS_ORDER2_TEACHER_ASSIGNON:
			resp := NewWSMessage(msg.MessageId, userId, WS_ORDER2_TEACHER_ASSIGNON_RESP)
			if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER {
				resp.Attribute["errCode"] = "0"
			} else {
				resp.Attribute["errCode"] = "2"
				resp.Attribute["errMsg"] = "You are not a teacher"
			}
			if err := TeacherManager.SetAssignOn(userId); err != nil {
				resp.Attribute["errCode"] = "2"
				resp.Attribute["errMsg"] = err.Error()
			}
			userChan <- resp

		case WS_ORDER2_TEACHER_ASSIGNOFF:
			resp := NewWSMessage(msg.MessageId, userId, WS_ORDER2_TEACHER_ASSIGNOFF_RESP)
			if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER {
				resp.Attribute["errCode"] = "0"
			} else {
				resp.Attribute["errCode"] = "2"
				resp.Attribute["errMsg"] = "You are not a teacher"
			}
			if err := TeacherManager.SetAssignOff(userId); err != nil {
				resp.Attribute["errCode"] = "2"
				resp.Attribute["errMsg"] = err.Error()
			}
			userChan <- resp

		case WS_ORDER2_CREATE:
			resp := NewWSMessage(msg.MessageId, userId, WS_ORDER2_CREATE_RESP)
			if err := InitOrderDispatch(msg, timestamp); err == nil {
				orderDispatchCountdown := settings.OrderDispatchCountdown()
				resp.Attribute["errCode"] = "0"
				resp.Attribute["countdown"] = strconv.FormatInt(orderDispatchCountdown, 10)
				resp.Attribute["countfrom"] = "0"
			} else {
				resp.Attribute["errCode"] = "2"
				resp.Attribute["errMsg"] = err.Error()
			}
			userChan <- resp

		case WS_ORDER2_PERSONAL_CHECK:
			resp := NewWSMessage(msg.MessageId, userId, WS_ORDER2_PERSONAL_CHECK_RESP)
			resp.Attribute["errCode"] = "0"

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

			status, err := CheckOrderValidation(orderId)
			resp.Attribute["status"] = strconv.FormatInt(status, 10)
			if err != nil {
				resp.Attribute["errMsg"] = err.Error()
			}
			userChan <- resp

		case WS_ORDER2_CANCEL,
			WS_ORDER2_ACCEPT,
			WS_ORDER2_ASSIGN_ACCEPT,
			WS_ORDER2_PERSONAL_REPLY:
			resp := NewWSMessage(msg.MessageId, userId, msg.OperationCode+1)
			resp.Attribute["errCode"] = "0"

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

			if orderChan, err := OrderManager.GetOrderChan(orderId); err != nil {
				resp.Attribute["errCode"] = "2"
				userChan <- resp
			} else {
				orderChan <- msg
			}
		}
	}
}

func WebSocketWriteHandler(conn *websocket.Conn, userId int64, userChan chan WSMessage) {
	// Time allowed to write a message to the peer.
	writeWaitInt := settings.WebsocketWriteWait()
	writeWait := time.Duration(writeWaitInt) * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriodInt := settings.WebsocketPingPeriod()
	pingPeriod := time.Duration(pingPeriodInt) * time.Second

	// 初始化心跳计时器
	pingTicker := time.NewTicker(pingPeriod)
	defer func() {
		pingTicker.Stop()
		conn.Close()
		seelog.Debug("WebSocketWriteHandler close websocket connection:", userId)
		if r := recover(); r != nil {
			seelog.Error(r)
		}
	}()

	loginTS := UserManager.GetUserOnlineStatus(userId)

	for {
		select {
		// 发送心跳
		case <-pingTicker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				seelog.Error("WebSocket Write Error: UserId", userId, "ErrMsg: ", err.Error())
				if UserManager.GetUserOnlineStatus(userId) == loginTS {
					WSUserLogout(userId)
					close(userChan)
				}
				return
			}
		// 处理向用户发送消息
		case msg, ok := <-userChan:
			if ok {
				conn.SetWriteDeadline(time.Now().Add(writeWait))
				err := conn.WriteJSON(msg)

				if err != nil {
					seelog.Error("WebSocket Write Error: UserId", userId, "ErrMsg: ", err.Error())
					if UserManager.GetUserOnlineStatus(userId) == loginTS {
						WSUserLogout(userId)
						close(userChan)
					}
					return
				}

				msgByte, err := json.Marshal(msg)
				if err == nil {
					seelog.Trace("WebSocketWriter: UserId: ", userId, "Msg: ", string(msgByte))
				}

				if msg.OperationCode == WS_FORCE_QUIT ||
					msg.OperationCode == WS_FORCE_LOGOUT ||
					msg.OperationCode == WS_LOGOUT_RESP {
					if UserManager.GetUserOnlineStatus(userId) == loginTS {
						WSUserLogout(userId)
						close(userChan)
						seelog.Trace("WebSocketWriter:User ", userId, " quit or logout!")
					}
					return
				}
			} else {
				return
			}
		}
	}
}
