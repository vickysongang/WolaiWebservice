package user

import (
	"WolaiWebservice/models"

	"github.com/astaxie/beego/orm"
)

type teacherItem struct {
	Id           int64    `json:"id"`
	Nickname     string   `json:"nickname"`
	Avatar       string   `json:"avatar"`
	Gender       int64    `json:"gender"`
	AccessRight  int64    `json:"accessRight"`
	School       string   `json:"school"`
	SubjectList  []string `json:"subjectList"`
	OnlineStatus string   `json:"onlineStatus"`
}

func GetTeacherRecommendation(userId int64, page int64, count int64) (int64, []teacherItem) {
	o := orm.NewOrm()

	var teachers []*models.TeacherProfile
	_, err := o.QueryTable("teacher_profile").OrderBy("-service_time").
		Offset(page * count).Limit(count).All(&teachers)
	if err != nil {
		return 2, nil
	}

	subjectDummy := []string{
		"数学",
		"英语",
		"物理",
	}

	result := make([]teacherItem, 0)
	for _, teacher := range teachers {
		user, _ := models.ReadUser(teacher.UserId)
		item := teacherItem{
			Id:           teacher.UserId,
			Nickname:     user.Nickname,
			Avatar:       user.Avatar,
			Gender:       user.Gender,
			AccessRight:  user.AccessRight,
			School:       "湖南大学",
			SubjectList:  subjectDummy,
			OnlineStatus: "online",
		}
		result = append(result, item)
	}

	return 0, result
}
