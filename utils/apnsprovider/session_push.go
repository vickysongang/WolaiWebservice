package apnsprovider

import (
	"encoding/json"

	"github.com/anachronistic/apns"

	"WolaiWebservice/models"
)

func PushSessionInstantStart(deviceToken, deviceProfile string, sessionId int64) error {
	var err error

	session, err := models.ReadSession(sessionId)
	if err != nil {
		return err
	}

	order, err := models.ReadOrder(session.OrderId)
	if err != nil {
		return err
	}

	payload := apns.NewPayload()
	payload.Alert = "上课开始了，快回到课堂吧"
	payload.Badge = 1

	pn := apns.NewPushNotification()
	pn.DeviceToken = deviceToken
	pn.AddPayload(payload)
	pn.Set("type", "session_instant_start")
	pn.Set("sessionId", session.Id)
	pn.Set("studentId", session.Creator)
	pn.Set("teacherId", session.Tutor)
	if order.Type == models.ORDER_TYPE_COURSE_INSTANT || order.Type == models.ORDER_TYPE_AUDITION_COURSE_INSTANT {
		pn.Set("courseId", order.CourseId)
	}

	return send(pn, deviceProfile)
}

func PushSessionResume(deviceToken, deviceProfile string, sessionId int64) error {
	var err error

	session, err := models.ReadSession(sessionId)
	if err != nil {
		return err
	}

	teacher, err := models.ReadUser(session.Tutor)
	if err != nil {
		return err
	}
	teacherByte, _ := json.Marshal(teacher)

	payload := apns.NewPayload()
	payload.Alert = "导师正在邀请你进入课堂"
	payload.Badge = 1

	pn := apns.NewPushNotification()
	pn.DeviceToken = deviceToken
	pn.AddPayload(payload)
	pn.Set("type", "session_resume")
	pn.Set("sessionId", session.Id)
	pn.Set("teacherId", session.Tutor)
	pn.Set("teacherInfo", string(teacherByte))

	return send(pn, deviceProfile)
}
