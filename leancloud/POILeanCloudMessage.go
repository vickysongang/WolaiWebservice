// POILeanCloudMessage
package leancloud

const (
	LC_MSG_TEXT        = -1
	LC_MSG_IMAGE       = 2
	LC_MSG_VOICE       = 3
	LC_MSG_DISCOVER    = 4
	LC_MSG_SESSION     = 5
	LC_MSG_SESSION_SYS = 6
	LC_MSG_WHITEBOARD  = 7
	LC_MSG_TRADE       = 8
	LC_MSG_AD          = 9

	LC_DISCOVER_TYPE_COMMENT = "0"
	LC_DISCOVER_TYPE_LIKE    = "1"

	LC_SESSION_REJECT   = "-1"
	LC_SESSION_PERSONAL = "1"
	LC_SESSION_CONFIRM  = "2"
	LC_SESSION_REMINDER = "3"
	LC_SESSION_CANCEL   = "4"
	LC_SESSION_REPORT   = "5"

	LC_TRADE_TYPE_SYSTEM    = "0"
	LC_TRADE_TYPE_TEACHER   = "1"
	LC_TRADE_TYPE_STUDENT   = "2"
	LC_TRADE_STATUS_INCOME  = "1"
	LC_TRADE_STATUS_EXPENSE = "2"
)

type LCMessage struct {
	SendId         string `json:"from_peer"`
	ConversationId string `json:"conv_id"`
	Message        string `json:"message"`
	Transient      bool   `json:"transient"`
}

type LCTypedMessage struct {
	Type      int64             `json:"_lctype"`
	Text      string            `json:"_lctext"`
	Attribute map[string]string `json:"_lcattrs,omitempty"`
}

type POIConversationParticipant struct {
	ConversationId string `json:"convId"`
	Participant    string `json:"participant"`
}

type POIConversationParticipants []POIConversationParticipant
