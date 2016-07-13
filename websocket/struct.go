package websocket

import (
	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
	"time"
)

type orderInfo struct {
	Id          int64        `json:"id"`
	CreatorInfo *models.User `json:"creatorInfo"`
	Title       string       `json:"title"`
	Status      string       `json:"status"`
	Type        string       `json:"orderType"`
	CreateTime  time.Time    `json:"createTime"`
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
	} else if order.Type == models.ORDER_TYPE_COURSE_INSTANT || order.Type == models.ORDER_TYPE_AUDITION_COURSE_INSTANT {
		course, _ := models.ReadCourse(order.CourseId)

		title = course.Name
	}

	info := orderInfo{
		Id:          order.Id,
		CreatorInfo: user,
		Title:       title,
		Status:      order.Status,
		Type:        order.Type,
		CreateTime:  order.CreateTime,
	}

	return &info
}

func getTeacherSubject(teacherId int64) []*models.Subject {
	o := orm.NewOrm()

	var teacherSubjects []*models.TeacherSubject
	num, err := o.QueryTable("teacher_to_subject").Filter("user_id", teacherId).All(&teacherSubjects)
	if err != nil {
		return nil
	}

	subjects := make([]*models.Subject, num)
	for i, teacherSubject := range teacherSubjects {
		subject, err := models.ReadSubject(teacherSubject.SubjectId)
		if err != nil {
			continue
		}

		subjects[i] = subject
	}

	return subjects
}
