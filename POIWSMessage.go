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

	WS_LOGIN        = 11
	WS_LOGIN_RESP   = 12
	WS_LOGOUT       = 13
	WS_LOGOUT_RESP  = 14
	WS_FORCE_LOGOUT = 16
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
