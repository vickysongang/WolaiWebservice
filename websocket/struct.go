package websocket

import (
	"WolaiWebservice/models"
)

type orderInfo struct {
	Id          int64        `json:"id"`
	CreatorInfo *models.User `json:"creatorInfo"`
	Title       string       `json:"title"`
}

func GetOrderInfo(orderId int64) *orderInfo {
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
