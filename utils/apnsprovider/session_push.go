package apnsprovider

import (
	"errors"

	"github.com/anachronistic/apns"

	"WolaiWebservice/models"
)

func PushSessionInstantStart(deviceToken string, sessionId int64) error {
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
	payload.Alert = "你有一条新的上课消息"
	payload.Badge = 1

	pn := apns.NewPushNotification()
	pn.DeviceToken = deviceToken
	pn.AddPayload(payload)
	pn.Set("type", "session_instant_start")
	pn.Set("sessionId", session.Id)
	pn.Set("studentId", session.Creator)
	pn.Set("teacherId", session.Tutor)
	if order.Type == models.ORDER_TYPE_COURSE_INSTANT {
		pn.Set("courseId", order.CourseId)
	}

	resp := apnsClient.Send(pn)
	if !resp.Success {
		return errors.New("推送失败")
	}

	return nil
}

func PushSessionResume(deviceToken string, sessionId int64) error {
	var err error

	session, err := models.ReadSession(sessionId)
	if err != nil {
		return err
	}

	payload := apns.NewPayload()
	payload.Alert = "你有一条新的上课消息"
	payload.Badge = 1

	pn := apns.NewPushNotification()
	pn.DeviceToken = deviceToken
	pn.AddPayload(payload)
	pn.Set("type", "session_resume")
	pn.Set("sessionId", session.Id)
	pn.Set("teacherId", session.Tutor)

	resp := apnsClient.Send(pn)
	if !resp.Success {
		return errors.New("推送失败")
	}

	return nil
}
