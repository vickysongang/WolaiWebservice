package user

import (
	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
	"WolaiWebservice/redis"
)

func GetUserInfo(userId int64) (int64, *models.User) {
	user, err := models.ReadUser(userId)
	if err != nil {
		return 2, nil
	}

	return 0, user
}

func UpdateUserInfo(userId int64, nickname string, avatar string, gender int64) (int64, *models.User) {
	user, err := models.UpdateUserInfo(userId, nickname, avatar, gender)
	if err != nil {
		return 2, nil
	}

	return 0, user
}

func UserLaunch(userId int64, objectId, address, ip, userAgent string) (int64, interface{}) {
	info := models.UserLoginInfo{
		UserId:    userId,
		ObjectId:  objectId,
		Address:   address,
		IP:        ip,
		UserAgent: userAgent,
	}

	models.CreateUserLoginInfo(&info)

	return 0, map[string]string{
		"websocket": redis.RedisManager.GetConfigStr(redis.CONFIG_GENERAL,
			redis.CONFIG_KEY_GENERAL_WEBSOCKET),
		"kamailio": redis.RedisManager.GetConfigStr(redis.CONFIG_GENERAL,
			redis.CONFIG_KEY_GENERAL_KAMAILIO),
	}
}

////////////////////////////////////////////////////////////////////////////////
///
///
////////////////////////////////////////////////////////////////////////////////

type teacherProfile struct {
	Id          int64                   `json:"id"`
	Nickname    string                  `json:"nickname"`
	Avatar      string                  `json:"avatar"`
	Gender      int64                   `json:"gender"`
	AccessRight int64                   `json:"accessRight"`
	School      string                  `json:"school"`
	Major       string                  `json:"major"`
	ServiceTime int64                   `json:"serviceTime"`
	SubjectList []string                `json:"subjectList,omitempty"`
	Intro       string                  `json:"intro"`
	Resume      []*models.TeacherResume `json:"resume,omitempty"`
}

func GetTeacherProfile(userId int64, teacherId int64) (int64, *teacherProfile) {
	o := orm.NewOrm()

	teacher, err := models.ReadTeacherProfile(teacherId)
	if err != nil {
		return 2, nil
	}

	school, err := models.ReadSchool(teacher.SchoolId)
	if err != nil {
		return 2, nil
	}

	user, err := models.ReadUser(teacherId)
	if err != nil {
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
		School:      school.Name,
		Major:       teacher.Major,
		ServiceTime: teacher.ServiceTime,
		SubjectList: subjectDummy,
		Intro:       teacher.Intro,
		Resume:      teacherResumes,
	}

	return 0, &profile
}
