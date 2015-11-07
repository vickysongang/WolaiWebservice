package websocket

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

	WS_SESSION_ALERT                   = 201
	WS_SESSION_ALERT_RESP              = 202
	WS_SESSION_START                   = 203
	WS_SESSION_START_RESP              = 204
	WS_SESSION_ACCEPT                  = 205
	WS_SESSION_ACCEPT_RESP             = 206
	WS_SESSION_PAUSE                   = 207
	WS_SESSION_PAUSE_RESP              = 208
	WS_SESSION_RESUME                  = 209
	WS_SESSION_RESUME_RESP             = 210
	WS_SESSION_FINISH                  = 211
	WS_SESSION_FINISH_RESP             = 212
	WS_SESSION_BREAK                   = 213
	WS_SESSION_BREAK_RESP              = 214
	WS_SESSION_SYNC                    = 215
	WS_SESSION_EXPIRE                  = 217
	WS_SESSION_EXPIRE_RESP             = 218
	WS_SESSION_CANCEL                  = 219
	WS_SESSION_CANCEL_RESP             = 220
	WS_SESSION_RECOVER_STU             = 221
	WS_SESSION_RECOVER_TEACHER         = 222
	WS_SESSION_INSTANT_ALERT           = 223
	WS_SESSION_INSTANT_START           = 225
	WS_SESSION_RESUME_CANCEL           = 227
	WS_SESSION_RESUME_CANCEL_RESP      = 228
	WS_SESSION_RESUME_ACCEPT           = 229
	WS_SESSION_RESUME_ACCEPT_RESP      = 230
	WS_SESSION_BREAK_RECONNECT_SUCCESS = 231

	WS_ORDER2_TEACHER_ONLINE         = 131
	WS_ORDER2_TEACHER_ONLINE_RESP    = 132
	WS_ORDER2_TEACHER_OFFLINE        = 133
	WS_ORDER2_TEACHER_OFFLINE_RESP   = 134
	WS_ORDER2_TEACHER_ASSIGNON       = 135
	WS_ORDER2_TEACHER_ASSIGNON_RESP  = 136
	WS_ORDER2_TEACHER_ASSIGNOFF      = 137
	WS_ORDER2_TEACHER_ASSIGNOFF_RESP = 138
	WS_ORDER2_RESULT                 = 139
	WS_ORDER2_RESULT_RESP            = 140
	WS_ORDER2_CREATE                 = 141
	WS_ORDER2_CREATE_RESP            = 142
	WS_ORDER2_CANCEL                 = 143
	WS_ORDER2_CANCEL_RESP            = 144
	WS_ORDER2_DISPATCH               = 145
	WS_ORDER2_DISPATCH_RESP          = 146
	WS_ORDER2_ACCEPT                 = 147
	WS_ORDER2_ACCEPT_RESP            = 148
	WS_ORDER2_EXPIRE                 = 149
	WS_ORDER2_EXPIRE_RESP            = 150
	WS_ORDER2_ASSIGN                 = 151
	WS_ORDER2_ASSIGN_RESP            = 152
	WS_ORDER2_ASSIGN_ACCEPT          = 153
	WS_ORDER2_ASSIGN_ACCEPT_RESP     = 154
	WS_ORDER2_ASSIGN_EXPIRE          = 155
	WS_ORDER2_ASSIGN_EXPIRE_RESP     = 156
	WS_ORDER2_PERSONAL_NOTIFY        = 161
	WS_ORDER2_PERSONAL_NOTIFY_RESP   = 162
	WS_ORDER2_PERSONAL_CHECK         = 163
	WS_ORDER2_PERSONAL_CHECK_RESP    = 164
	WS_ORDER2_PERSONAL_REPLY         = 165
	WS_ORDER2_PERSONAL_REPLY_RESP    = 166
	WS_ORDER2_RECOVER_DISPATCH       = 171
	WS_ORDER2_RECOVER_ASSIGN         = 173
	WS_ORDER2_RECOVER_CREATE         = 175
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
