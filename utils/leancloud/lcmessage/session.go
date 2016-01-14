package lcmessage

import (
	"WolaiWebservice/models"
	"WolaiWebservice/utils/leancloud"
)

/*
 * 各种系统消息
 */
func SendSessionStartMsg(sessionId int64) {
	session, err := models.ReadSession(sessionId)
	if err != nil {
		return
	}

	attr := make(map[string]string)
	attr["accessRight"] = "[2, 3]"

	lcTMsg := leancloud.LCTypedMessage{
		Type:      LC_MSG_SYSTEM,
		Text:      "已进入课堂",
		Attribute: attr,
	}

	leancloud.LCSendSystemMessage(USER_SYSTEM_MESSAGE, session.Creator, session.Tutor, &lcTMsg)
}

func SendSessionFinishMsg(sessionId int64) {
	session, err := models.ReadSession(sessionId)
	if err != nil {
		return
	}

	attr := make(map[string]string)
	attr["accessRight"] = "[2, 3]"

	lcTMsg := leancloud.LCTypedMessage{
		Type:      LC_MSG_SYSTEM,
		Text:      "上课结束，别忘了留下评价哦",
		Attribute: attr,
	}

	leancloud.LCSendSystemMessage(USER_SYSTEM_MESSAGE, session.Creator, session.Tutor, &lcTMsg)
}

func SendSessionExpireMsg(sessionId int64) {
	session, err := models.ReadSession(sessionId)
	if err != nil {
		return
	}

	attr := make(map[string]string)
	attr["accessRight"] = "[2, 3]"

	lcTMsg := leancloud.LCTypedMessage{
		Type:      LC_MSG_SYSTEM,
		Text:      "上课中断，建议沟通后继续上课哦",
		Attribute: attr,
	}

	leancloud.LCSendSystemMessage(USER_SYSTEM_MESSAGE, session.Creator, session.Tutor, &lcTMsg)
}

func SendSessionBreakMsg(sessionId int64) {
	session, err := models.ReadSession(sessionId)
	if err != nil {
		return
	}

	attr := make(map[string]string)
	attr["accessRight"] = "[2, 3]"

	lcTMsg := leancloud.LCTypedMessage{
		Type:      LC_MSG_SYSTEM,
		Text:      "上课暂时中断，需要静静重连一下",
		Attribute: attr,
	}

	leancloud.LCSendSystemMessage(USER_SYSTEM_MESSAGE, session.Creator, session.Tutor, &lcTMsg)
}

func SendSessionResumeMsg(sessionId int64) {
	session, err := models.ReadSession(sessionId)
	if err != nil {
		return
	}

	attr := make(map[string]string)
	attr["accessRight"] = "[2, 3]"

	lcTMsg := leancloud.LCTypedMessage{
		Type:      LC_MSG_SYSTEM,
		Text:      "静静说可以继续上课啦",
		Attribute: attr,
	}

	leancloud.LCSendSystemMessage(USER_SYSTEM_MESSAGE, session.Creator, session.Tutor, &lcTMsg)
}
