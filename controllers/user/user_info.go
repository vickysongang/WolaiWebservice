package user

import (
	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
)

func GetUserInfo(userId int64) (int64, *models.User) {
	user, err := models.ReadUser(userId)
	if err != nil {
		return 2, nil
	}

	return 0, user
}

func UpdateUserInfo(userId int64, nickname string, avatar string, gender int64) (int64, *models.User) {
	user, err := models.UpdateUser(userId, nickname, avatar, gender)
	if err != nil {
		return 2, nil
	}

	return 0, user
}

type teacherProfile struct {
	Id          int64                   `json:"id"`
	Nickname    string                  `json:"nickname"`
	Avatar      string                  `json:"avatar"`
	Gender      int64                   `json:"gender"`
	AccessRight int64                   `json:"accessRight"`
	School      string                  `json:"school"`
	Major       string                  `json:"major"`
	SubjectList []string                `json:"subjectList,omitempty"`
	Intro       string                  `json:"intro"`
	Resume      []*models.TeacherResume `json:"resume,omitempty"`
}

func GetTeacherProfile(userId int64, teacherId int64) (int64, *teacherProfile) {
	o := orm.NewOrm()

	teacher := models.TeacherProfile{UserId: teacherId}
	err := o.Read(&teacher)
	if err != nil {
		println(err.Error())
		return 2, nil
	}

	user, err := models.ReadUser(teacherId)
	if err != nil {
		println(err.Error())
		return 2, nil
	}

	subjectDummy := []string{
		"数学",
		"英语",
		"物理",
	}

	var teacherResumes []*models.TeacherResume
	_, err = o.QueryTable("teacher_to_resume").Filter("user_id", teacherId).All(&teacherResumes)
	if err != nil {
		println(err.Error())
		return 2, nil
	}

	profile := teacherProfile{
		Id:          user.Id,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Gender:      user.Gender,
		AccessRight: user.AccessRight,
		School:      "湖南大学",
		Major:       "化学系",
		SubjectList: subjectDummy,
		Intro:       teacher.Intro,
		Resume:      teacherResumes,
	}

	return 0, &profile
}
