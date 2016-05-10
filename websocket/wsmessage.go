package websocket

import (
	"time"

	"github.com/satori/go.uuid"
)

type WSMessage struct {
	MessageId     string            `json:"msgId"`
	UserId        int64             `json:"userId"`
	OperationCode int64             `json:"oprCode"`
	Timestamp     float64           `json:"timestamp"`
	Attribute     map[string]string `json:"attr"`
}

const (
	WS_FORCE_QUIT       = -1
	SIGNAL_ORDER_QUIT   = -2
	SIGNAL_SESSION_QUIT = -3

	WS_PING = 0
	WS_PONG = 1

	WS_LOGIN          = 11
	WS_LOGIN_RESP     = 12
	WS_LOGOUT         = 13
	WS_LOGOUT_RESP    = 14
	WS_FORCE_LOGOUT   = 16
	WS_RECONNECT      = 17
	WS_RECONNECT_RESP = 18

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
	WS_SESSION_RECOVER_STU             = 221
	WS_SESSION_RECOVER_TEACHER         = 222
	WS_SESSION_INSTANT_ALERT           = 223
	WS_SESSION_INSTANT_START           = 225
	WS_SESSION_RESUME_CANCEL           = 227
	WS_SESSION_RESUME_CANCEL_RESP      = 228
	WS_SESSION_RESUME_ACCEPT           = 229
	WS_SESSION_RESUME_ACCEPT_RESP      = 230
	WS_SESSION_BREAK_RECONNECT_SUCCESS = 231
	WS_SESSION_STATUS_SYNC             = 233

	WS_SESSION_ASK_FINISH             = 235 // 学生C端主动发起结束，S收到后回复WS_SESSION_ASK_FINISH_RESP，并向老师C端发WS_SESSION_ASK_FINISH
	WS_SESSION_ASK_FINISH_RESP        = 236
	WS_SESSION_ASK_FINISH_REJECT      = 237 // 老师C端拒绝下课请求，S端收到后向学生C端再发WS_SESSION_ASK_FINISH_REJECT，并且回复老师WS_SESSION_ASK_FINISH_REJECT_RESP；如果同意C端发WS_SESSION_FINISH
	WS_SESSION_ASK_FINISH_REJECT_RESP = 238

	WS_SESSION_CONTINUE      = 239 //老师C端从主动暂停中点击恢复，直接开始重新计时
	WS_SESSION_CONTINUE_RESP = 240

	WS_SESSION_QAPKG_TIME_END     = 241
	WS_SESSION_AUTO_FINISH_TIP    = 243
	WS_SESSION_NOT_EVALUATION_TIP = 245

	WS_SESSION_REPORT = 251
)

func NewWSMessage(msgId string, userId int64, oprCode int64) WSMessage {
	timestampNano := time.Now().UnixNano()
	timestampMillis := timestampNano / 1000
	timestamp := float64(timestampMillis) / 1000000.0

	if msgId == "" {
		msgId = uuid.NewV4().String()
	}
	return WSMessage{
		MessageId:     msgId,
		UserId:        userId,
		OperationCode: oprCode,
		Timestamp:     timestamp,
		Attribute:     make(map[string]string),
	}
}
