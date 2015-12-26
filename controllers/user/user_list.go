package user

import (
	"github.com/astaxie/beego/orm"

	"WolaiWebservice/config"
	"WolaiWebservice/models"
	"WolaiWebservice/websocket"
)

type teacherItem struct {
	Id           int64    `json:"id"`
	Nickname     string   `json:"nickname"`
	Avatar       string   `json:"avatar"`
	Gender       int64    `json:"gender"`
	AccessRight  int64    `json:"accessRight"`
	School       string   `json:"school"`
	Major        string   `json:"major"`
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
		var schoolStr string
		if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER {
			profile, err := models.ReadTeacherProfile(user.Id)
			if err != nil {
				continue
			}

			school, err := models.ReadSchool(profile.SchoolId)
			if err == nil {
				schoolStr = school.Name
			}
		}

		item := teacherItem{
			Id:           user.Id,
			Nickname:     user.Nickname,
			Avatar:       user.Avatar,
			Gender:       user.Gender,
			AccessRight:  user.AccessRight,
			School:       schoolStr,
			SubjectList:  nil,
			OnlineStatus: "",
		}

		if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER {
			profile, _ := models.ReadTeacherProfile(user.Id)
			item.Major = profile.Major
		}

		result = append(result, item)
	}

	return 0, result
}

func GetTeacherRecent(userId int64, page int64, count int64) (int64, []teacherItem) {
	o := orm.NewOrm()

	type sessionTeacher struct {
		Tutor int64
	}
	var teacherIds []sessionTeacher

	qb, _ := orm.NewQueryBuilder(config.Env.Database.Type)
	qb.Select("distinct tutor").From("sessions").Where("creator = ?").
		Limit(int(count)).Offset(int(page * count))
	sql := qb.String()
	o.Raw(sql, userId).QueryRows(&teacherIds)

	result := make([]teacherItem, 0)
	for _, teacherId := range teacherIds {
		user, err := models.ReadUser(teacherId.Tutor)
		if err != nil {
			continue
		}

		profile, err := models.ReadTeacherProfile(teacherId.Tutor)
		if err != nil {
			continue
		}

		var schoolStr string
		school, err := models.ReadSchool(profile.SchoolId)
		if err == nil {
			schoolStr = school.Name
		}

		subjects := GetTeacherSubject(profile.UserId)
		var subjectNames []string
		if subjects != nil {
			subjectNames = ParseSubjectNameSlice(subjects)
		} else {
			subjectNames = make([]string, 0)
		}

		item := teacherItem{
			Id:           profile.UserId,
			Nickname:     user.Nickname,
			Avatar:       user.Avatar,
			Gender:       user.Gender,
			AccessRight:  user.AccessRight,
			School:       schoolStr,
			SubjectList:  subjectNames,
			OnlineStatus: websocket.WsManager.GetUserStatus(user.Id),
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

	result := make([]teacherItem, 0)
	for _, teacher := range teachers {
		user, err := models.ReadUser(teacher.UserId)
		if err != nil {
			continue
		}

		var schoolStr string
		school, err := models.ReadSchool(teacher.SchoolId)
		if err == nil {
			schoolStr = school.Name
		}

		subjects := GetTeacherSubject(teacher.UserId)
		var subjectNames []string
		if subjects != nil {
			subjectNames = ParseSubjectNameSlice(subjects)
		} else {
			subjectNames = make([]string, 0)
		}

		item := teacherItem{
			Id:           teacher.UserId,
			Nickname:     user.Nickname,
			Avatar:       user.Avatar,
			Gender:       user.Gender,
			AccessRight:  user.AccessRight,
			School:       schoolStr,
			SubjectList:  subjectNames,
			OnlineStatus: websocket.WsManager.GetUserStatus(user.Id),
		}
		result = append(result, item)
	}

	return 0, result
}

func GetContactRecommendation(userId int64, page int64, count int64) (int64, []teacherItem) {
	o := orm.NewOrm()

	result := make([]teacherItem, 0)

	if page == 0 {
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
	}

	var teachers []*models.TeacherProfile
	_, err := o.QueryTable("teacher_profile").OrderBy("-service_time").
		Offset(page * count).Limit(count).All(&teachers)
	if err != nil {
		return 2, nil
	}

	for _, teacher := range teachers {
		user, err := models.ReadUser(teacher.UserId)
		if err != nil {
			continue
		}

		var schoolStr string
		school, err := models.ReadSchool(teacher.SchoolId)
		if err == nil {
			schoolStr = school.Name
		}

		item := teacherItem{
			Id:           teacher.UserId,
			Nickname:     user.Nickname,
			Avatar:       user.Avatar,
			Gender:       user.Gender,
			AccessRight:  user.AccessRight,
			School:       schoolStr,
			SubjectList:  nil,
			OnlineStatus: "",
		}
		result = append(result, item)
	}

	return 0, result
}
