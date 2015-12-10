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
	SubjectList  []string `json:"subjectList,omitempty"`
	OnlineStatus string   `json:"onlineStatus,omitempty"`
}

func SearchUser(userId int64, keyword string, page, count int64) (int64, []teacherItem) {
	o := orm.NewOrm()

	cond := orm.NewCondition()
	cond1 := cond.And("access_right__in", 2, 3).And("status", 0)
	cond2 := cond.Or("nickname__icontains", keyword).Or("phone__icontains", keyword)
	condFin := cond.AndCond(cond1).AndCond(cond2)

	var users []*models.User
	_, err := o.QueryTable("users").SetCond(condFin).
		Offset(page * count).Limit(count).All(&users)
	if err != nil {
		return 2, nil
	}

	result := make([]teacherItem, 0)
	for _, user := range users {
		item := teacherItem{
			Id:           user.Id,
			Nickname:     user.Nickname,
			Avatar:       user.Avatar,
			Gender:       user.Gender,
			AccessRight:  user.AccessRight,
			School:       "",
			SubjectList:  nil,
			OnlineStatus: "",
		}
		result = append(result, item)
	}

	return 0, result
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

func GetContactRecommendation(userId int64, page int64, count int64) (int64, []teacherItem) {
	o := orm.NewOrm()

	result := make([]teacherItem, 0)

	wolaiTeam, err := models.ReadUser(models.USER_WOLAI_TEAM)
	wolaiItem := teacherItem{
		Id:           wolaiTeam.Id,
		Nickname:     wolaiTeam.Nickname,
		Avatar:       wolaiTeam.Avatar,
		Gender:       wolaiTeam.Gender,
		AccessRight:  wolaiTeam.AccessRight,
		School:       "",
		SubjectList:  nil,
		OnlineStatus: "",
	}
	result = append(result, wolaiItem)

	var users []*models.User
	_, err = o.QueryTable("users").Filter("access_right", 1).All(&users)
	if err == nil {
		for _, user := range users {
			item := teacherItem{
				Id:           user.Id,
				Nickname:     user.Nickname,
				Avatar:       user.Avatar,
				Gender:       user.Gender,
				AccessRight:  user.AccessRight,
				School:       "",
				SubjectList:  nil,
				OnlineStatus: "",
			}
			result = append(result, item)
		}
	}

	var teachers []*models.TeacherProfile
	_, err = o.QueryTable("teacher_profile").OrderBy("-service_time").
		Offset(page * count).Limit(count).All(&teachers)
	if err != nil {
		return 2, nil
	}

	for _, teacher := range teachers {
		user, _ := models.ReadUser(teacher.UserId)
		item := teacherItem{
			Id:           teacher.UserId,
			Nickname:     user.Nickname,
			Avatar:       user.Avatar,
			Gender:       user.Gender,
			AccessRight:  user.AccessRight,
			School:       "湖南大学",
			SubjectList:  nil,
			OnlineStatus: "",
		}
		result = append(result, item)
	}

	return 0, result
}
