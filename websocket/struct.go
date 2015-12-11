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

type teacherInfo struct {
	Id          int64  `json:"id"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	Gender      int64  `json:"Gender"`
	AccessRight int64  `json:"accessRight"`
	School      string `json:"school"`
	Major       string `json:"major"`
	ServiceTime int64  `json:"serviceTime"`
}

func GetTeacherInfo(teacherId int64) *teacherInfo {
	teacher, _ := models.ReadTeacherProfile(teacherId)
	user, _ := models.ReadUser(teacherId)

	info := teacherInfo{
		Id:          user.Id,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Gender:      user.Gender,
		AccessRight: user.AccessRight,
		School:      "湖南大学",
		Major:       teacher.Major,
		ServiceTime: teacher.ServiceTime,
	}

	return &info
}
