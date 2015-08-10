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

func NewCloseMessage(userId int64) POIWSMessage {
	timestampNano := time.Now().UnixNano()
	timestampMillis := timestampNano / 1000
	timestamp := float64(timestampMillis) / 1000000.0

	return POIWSMessage{
		MessageId:     uuid.NewV4().String(),
		UserId:        userId,
		OperationCode: -1,
		Timestamp:     timestamp,
		Attribute:     make(map[string]string),
	}
}

func NewType1Message() POIWSMessage {
	timestampNano := time.Now().UnixNano()
	timestampMillis := timestampNano / 1000
	timestamp := float64(timestampMillis) / 1000000.0

	return POIWSMessage{
		MessageId:     uuid.NewV4().String(),
		UserId:        10001,
		OperationCode: 1,
		Timestamp:     timestamp,
		Attribute:     make(map[string]string),
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
		Attribute:     make(map[string]string),
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
		Attribute:     make(map[string]string),
	}
}

func NewType5Message() POIWSMessage {
	timestampNano := time.Now().UnixNano()
	timestampMillis := timestampNano / 1000
	timestamp := float64(timestampMillis) / 1000000.0

	return POIWSMessage{
		MessageId:     uuid.NewV4().String(),
		UserId:        10001,
		OperationCode: 5,
		Timestamp:     timestamp,
		Attribute:     make(map[string]string),
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
		Attribute:     make(map[string]string),
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
		Attribute:     make(map[string]string),
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
		Attribute:     make(map[string]string),
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
		Attribute:     make(map[string]string),
	}
}
