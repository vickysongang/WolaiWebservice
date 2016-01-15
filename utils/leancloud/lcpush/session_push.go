package lcpush

import (
	"WolaiWebservice/models"
	"WolaiWebservice/utils/leancloud"
)

func PushSessionInstantStart(objectId string, sessionId int64) error {
	var err error

	session, err := models.ReadSession(sessionId)
	if err != nil {
		return err
	}

	order, err := models.ReadOrder(session.OrderId)
	if err != nil {
		return err
	}

	var alert string
	alert = "上课开始了，快回到课堂吧"

	lcReq := map[string]interface{}{
		"where": map[string]interface{}{
			"objectId": objectId,
		},
		"data": map[string]interface{}{
			"alert":     alert,
			"title":     "我来",
			"action":    "poi.push.REQUEST_BACK_TO_SESSION",
			"sessionId": sessionId,
			"studentId": session.Creator,
			"teacherId": session.Tutor,
			"courseId":  order.CourseId,
		},
	}

	go leancloud.LCPushNotification(&lcReq)

	return nil
}

func PushSessionResume(objectId string, sessionId int64) error {
	var err error

	session, err := models.ReadSession(sessionId)
	if err != nil {
		return err
	}

	var alert string
	alert = "导师正在邀请你进入课堂"

	lcReq := map[string]interface{}{
		"where": map[string]interface{}{
			"objectId": objectId,
		},
		"data": map[string]interface{}{
			"alert":     alert,
			"title":     "我来",
			"action":    "poi.push.RESUME_SESSION",
			"sessionId": sessionId,
			"teacherId": session.Tutor,
		},
	}

	go leancloud.LCPushNotification(&lcReq)

	return nil
}
