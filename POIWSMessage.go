package main

import (
	_ "encoding/json"
	"github.com/satori/go.uuid"
	"time"
)

type POIWSMessage struct {
	MessageId     string                 `json:"msgId"`
	UserId        int64                  `json:"userId"`
	OperationCode int64                  `json:"oprCode"`
	Timestamp     float64                `json:"timestamp"`
	Attribute     map[string]interface{} `json:"attr"`
}

func NewCloseMessage(userId int64) POIWSMessage {
	timestampNano := time.Now().UnixNano()
	timestampMillis := timestampNano / 1000
	timestamp := float64(timestampMillis) / 1000000.0

	return POIWSMessage{
		MessageId:     uuid.NewV4().String(),
		UserId:        userId,
		OperationCode: -1,
		Timestamp:     timestamp,
		Attribute:     make(map[string]interface{}),
	}
}

func NewType2Message() POIWSMessage {
	timestampNano := time.Now().UnixNano()
	timestampMillis := timestampNano / 1000
	timestamp := float64(timestampMillis) / 1000000.0

	return POIWSMessage{
		MessageId:     uuid.NewV4().String(),
		UserId:        10001,
		OperationCode: 2,
		Timestamp:     timestamp,
		Attribute:     make(map[string]interface{}),
	}
}

func NewType3Message() POIWSMessage {
	timestampNano := time.Now().UnixNano()
	timestampMillis := timestampNano / 1000
	timestamp := float64(timestampMillis) / 1000000.0

	return POIWSMessage{
		MessageId:     uuid.NewV4().String(),
		UserId:        10001,
		OperationCode: 3,
		Timestamp:     timestamp,
		Attribute:     make(map[string]interface{}),
	}
}

func NewType6Message() POIWSMessage {
	timestampNano := time.Now().UnixNano()
	timestampMillis := timestampNano / 1000
	timestamp := float64(timestampMillis) / 1000000.0

	return POIWSMessage{
		MessageId:     uuid.NewV4().String(),
		UserId:        10001,
		OperationCode: 6,
		Timestamp:     timestamp,
		Attribute:     make(map[string]interface{}),
	}
}

func NewType7Message() POIWSMessage {
	timestampNano := time.Now().UnixNano()
	timestampMillis := timestampNano / 1000
	timestamp := float64(timestampMillis) / 1000000.0

	return POIWSMessage{
		MessageId:     uuid.NewV4().String(),
		UserId:        10001,
		OperationCode: 7,
		Timestamp:     timestamp,
		Attribute:     make(map[string]interface{}),
	}
}

func NewType10Message() POIWSMessage {
	timestampNano := time.Now().UnixNano()
	timestampMillis := timestampNano / 1000
	timestamp := float64(timestampMillis) / 1000000.0

	return POIWSMessage{
		MessageId:     uuid.NewV4().String(),
		UserId:        10001,
		OperationCode: 10,
		Timestamp:     timestamp,
		Attribute:     make(map[string]interface{}),
	}
}

func NewType11Message() POIWSMessage {
	timestampNano := time.Now().UnixNano()
	timestampMillis := timestampNano / 1000
	timestamp := float64(timestampMillis) / 1000000.0

	return POIWSMessage{
		MessageId:     uuid.NewV4().String(),
		UserId:        10001,
		OperationCode: 11,
		Timestamp:     timestamp,
		Attribute:     make(map[string]interface{}),
	}
}
