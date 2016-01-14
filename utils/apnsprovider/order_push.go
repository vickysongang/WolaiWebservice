package apnsprovider

import (
	"encoding/json"
	"errors"

	"github.com/anachronistic/apns"

	"WolaiWebservice/config/settings"
	"WolaiWebservice/models"
)

func PushNewOrderDispatch(deviceToken string, orderId int64) error {
	info := getOrderInfo(orderId)
	infoByte, _ := json.Marshal(info)

	payload := apns.NewPayload()
	payload.Alert = "你有一条新的上课请求"
	payload.Badge = 1

	pn := apns.NewPushNotification()
	pn.DeviceToken = deviceToken
	pn.AddPayload(payload)
	pn.Set("type", "order_dispatch")
	pn.Set("orderInfo", string(infoByte))

	resp := apnsClient.Send(pn)
	if !resp.Success {
		return errors.New("推送失败")
	}

	return nil
}

func PushNewOrderAssign(deviceToken string, orderId int64) error {
	info := getOrderInfo(orderId)
	infoByte, _ := json.Marshal(info)
	orderAssignCountdown := settings.OrderAssignCountdown()

	payload := apns.NewPayload()
	payload.Alert = "你有一条新的指派订单"
	payload.Badge = 1
	payload.Sound = "iOS_new_orde_assign.aif"

	pn := apns.NewPushNotification()
	pn.DeviceToken = deviceToken
	pn.AddPayload(payload)
	pn.Set("type", "order_assign")
	pn.Set("orderInfo", string(infoByte))
	pn.Set("countdown", orderAssignCountdown)

	resp := apnsClient.Send(pn)
	if !resp.Success {
		return errors.New("推送失败")
	}

	return nil
}

func PushOrderAccept(deviceToken string, orderId, teacherId int64) error {
	info := getOrderInfo(orderId)
	orderSessionCountdown := settings.OrderSessionCountdown()
	teacher, err := models.ReadUser(teacherId)
	if err != nil {
		return err
	}
	teacherByte, _ := json.Marshal(teacher)

	payload := apns.NewPayload()
	payload.Alert = "有导师接受了你的提问，快来上课吧"
	payload.Badge = 1

	pn := apns.NewPushNotification()
	pn.DeviceToken = deviceToken
	pn.AddPayload(payload)
	pn.Set("type", "order_accept")
	pn.Set("orderId", orderId)
	pn.Set("countdown", orderSessionCountdown)
	pn.Set("teacherInfo", string(teacherByte))
	pn.Set("title", info.Title)

	resp := apnsClient.Send(pn)
	if !resp.Success {
		return errors.New("推送失败")
	}

	return nil
}

func PushOrderPersonalAccept(deviceToken string, orderId, teacherId int64) error {
	info := getOrderInfo(orderId)
	orderSessionCountdown := settings.OrderSessionCountdown()
	teacher, err := models.ReadUser(teacherId)
	if err != nil {
		return err
	}
	teacherByte, _ := json.Marshal(teacher)

	payload := apns.NewPayload()
	payload.Alert = "导师接受了你的提问，快来上课吧"
	payload.Badge = 1

	pn := apns.NewPushNotification()
	pn.DeviceToken = deviceToken
	pn.AddPayload(payload)
	pn.Set("type", "order_accept")
	pn.Set("orderId", orderId)
	pn.Set("countdown", orderSessionCountdown)
	pn.Set("teacherInfo", string(teacherByte))
	pn.Set("title", info.Title)

	resp := apnsClient.Send(pn)
	if !resp.Success {
		return errors.New("推送失败")
	}

	return nil
}

type orderInfo struct {
	Id          int64        `json:"id"`
	CreatorInfo *models.User `json:"creatorInfo"`
	Title       string       `json:"title"`
}

func getOrderInfo(orderId int64) *orderInfo {
	order, _ := models.ReadOrder(orderId)
	user, _ := models.ReadUser(order.Creator)

	var title string
	if order.Type == models.ORDER_TYPE_PERSONAL_INSTANT ||
		order.Type == models.ORDER_TYPE_GENERAL_INSTANT {
		grade, err1 := models.ReadGrade(order.GradeId)
		subject, err2 := models.ReadSubject(order.SubjectId)

		if err1 == nil && err2 == nil {
			title = grade.Name + subject.Name
		} else {
			title = "实时课堂"
		}
	} else if order.Type == models.ORDER_TYPE_COURSE_INSTANT {
		course, _ := models.ReadCourse(order.CourseId)

		title = course.Name
	}

	info := orderInfo{
		Id:          order.Id,
		CreatorInfo: user,
		Title:       title,
	}

	return &info
}
