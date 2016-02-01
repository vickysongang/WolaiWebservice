package lcmessage

import (
	"WolaiWebservice/utils/leancloud"
)

func SendWelcomeMessageTeacher(userId int64) {
	attr := make(map[string]string)

	msg := leancloud.LCTypedMessage{
		Type:      LC_MSG_TEXT,
		Text:      "Hi~ 欢迎你加入“我来”导师家族，陪伴学弟学妹们成长！\n你是百里挑一的精英学霸，你是闪闪发光的榜样力量~\n现在点击首页“开始答疑”，马上开启你的“超人之旅”！",
		Attribute: attr,
	}
	leancloud.LCSendTypedMessage(USER_WOLAI_SUPPORT, userId, &msg)
}

func SendWelcomeMessageStudent(userId int64) {
	attr := make(map[string]string)

	msg := leancloud.LCTypedMessage{
		Type:      LC_MSG_TEXT,
		Text:      "Hi~你终于来了，欢迎加入最温暖的“我来”学院~\n我来团队携手全国86所顶尖高校学霸导师，与你共度学习的美好时光！\n现在就回到首页去开启你的“我来奇妙之旅”吧！",
		Attribute: attr,
	}
	leancloud.LCSendTypedMessage(USER_WOLAI_SUPPORT, userId, &msg)
}
