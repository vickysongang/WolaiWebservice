package main

import (
	"time"

	"github.com/satori/go.uuid"
)

type POIWSMessage struct {
	MessageId     string            `json:"msgId"`
	UserId        int64             `json:"userId"`
	OperationCode int64             `json:"oprCode"`
	Timestamp     float64           `json:"timestamp"`
	Attribute     map[string]string `json:"attr"`
}

const (
	WS_FORCE_QUIT = -1

	WS_PING = 0
	WS_PONG = 1

	WS_LOGIN          = 11
	WS_LOGIN_RESP     = 12
	WS_LOGOUT         = 13
	WS_LOGOUT_RESP    = 14
	WS_FORCE_LOGOUT   = 16
	WS_RECONNECT      = 17
	WS_RECONNECT_RESP = 18

	WS_ORDER_TEACHER_ONLINE       = 101
	WS_ORDER_TEACHER_RESP         = 102
	WS_ORDER_TEACHER_OFFLINE      = 103
	WS_ORDER_TEACHER_OFFLINE_RESP = 104
	WS_ORDER_CREATE               = 105
	WS_ORDER_CREATE_RESP          = 106
	WS_ORDER_DISPATCH             = 107
	WS_ORDER_DISPATCH_RESP        = 108
	WS_ORDER_REPLY                = 109
	WS_ORDER_REPLY_RESP           = 110
	WS_ORDER_PRESENT              = 111
	WS_ORDER_PRESENT_RESP         = 112
	WS_ORDER_CONFIRM              = 113
	WS_ORDER_CONFIRM_RESP         = 114
	WS_ORDER_RESULT               = 115
	WS_ORDER_RESULT_RESP          = 116
	WS_ORDER_CANCEL               = 117
	WS_ORDER_CANCEL_RESP          = 118
	WS_ORDER_EXPIRE               = 119
	WS_ORDER_EXPIRE_RESP          = 120
	WS_ORDER_RECOVER_STU          = 121
	WS_ORDER_RECOVER_TEACHER      = 122

	WS_SESSION_ALERT           = 201
	WS_SESSION_ALERT_RESP      = 202
	WS_SESSION_START           = 203
	WS_SESSION_START_RESP      = 204
	WS_SESSION_ACCEPT          = 205
	WS_SESSION_ACCEPT_RESP     = 206
	WS_SESSION_PAUSE           = 207
	WS_SESSION_PAUSE_RESP      = 208
	WS_SESSION_RESUME          = 209
	WS_SESSION_RESUME_RESP     = 210
	WS_SESSION_FINISH          = 211
	WS_SESSION_FINISH_RESP     = 212
	WS_SESSION_BREAK           = 213
	WS_SESSION_BREAK_RESP      = 214
	WS_SESSION_SYNC            = 215
	WS_SESSION_EXPIRE          = 217
	WS_SESSION_EXPIRE_RESP     = 218
	WS_SESSION_CANCEL          = 219
	WS_SESSION_CANCEL_RESP     = 220
	WS_SESSION_RECOVER_STU     = 221
	WS_SESSION_RECOVER_TEACHER = 222
	WS_SESSION_INSTANT_ALERT   = 223
	WS_SESSION_INSTANT_START   = 225
)

func NewPOIWSMessage(msgId string, userId int64, oprCode int64) POIWSMessage {
	timestampNano := time.Now().UnixNano()
	timestampMillis := timestampNano / 1000
	timestamp := float64(timestampMillis) / 1000000.0

	if msgId == "" {
		msgId = uuid.NewV4().String()
	}
	return POIWSMessage{
		MessageId:     msgId,
		UserId:        userId,
		OperationCode: oprCode,
		Timestamp:     timestamp,
		Attribute:     make(map[string]string),
	}
}
