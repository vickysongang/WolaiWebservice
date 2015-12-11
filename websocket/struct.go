package websocket

import (
	"WolaiWebservice/models"
)

type orderInfo struct {
	OrderId     int64        `json:"orderId"`
	CreatorInfo *models.User `json:"creatorInfo"`
	Title       string       `json:"title"`
}

func GetOrderInfo(orderId int64) *orderInfo {
	order, _ := models.ReadOrder(orderId)
	user, _ := models.ReadUser(order.Creator)

	grade, err1 := models.ReadGrade(order.GradeId)
	subject, err2 := models.ReadSubject(order.SubjectId)

	var title string
	if err1 == nil && err2 == nil {
		title = grade.Name + subject.Name
	} else {
		title = "实时课堂"
	}

	info := orderInfo{
		OrderId:     order.Id,
		CreatorInfo: user,
		Title:       title,
	}

	return &info
}
