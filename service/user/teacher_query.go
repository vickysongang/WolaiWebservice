package user

import (
	"errors"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"

	"WolaiWebservice/config"
	"WolaiWebservice/models"
)

func QueryTeacherBySessionFreq(userId, page, count int64) ([]int64, error) {
	var err error

	o := orm.NewOrm()

	result := make([]int64, 0)

	type sessionTeacher struct {
		Tutor int64
	}
	var teacherIds []sessionTeacher

	qb, err := orm.NewQueryBuilder(config.Env.Database.Type)
	if err != nil {
		seelog.Errorf("%s", err.Error())
		return nil, errors.New("数据库查询失败")
	}

	qb.Select("tutor").From(new(models.Session).TableName()).
		Where("creator = ?").
		GroupBy("tutor").
		OrderBy("count(id) DESC, create_time DESC").
		Limit(int(count)).Offset(int(page * count))
	sql := qb.String()
	o.Raw(sql, userId).QueryRows(&teacherIds)

	for _, teacherId := range teacherIds {
		result = append(result, teacherId.Tutor)
	}

	return result, nil
}

func QueryTeacherRecommendation(userId, page, count int64) ([]int64, error) {
	var err error

	o := orm.NewOrm()

	var teachers []*models.TeacherProfile
	num, err := o.QueryTable(new(models.TeacherProfile).TableName()).
		OrderBy("-service_time").
		Limit(count).Offset(page * count).
		All(&teachers)
	if err != nil {
		return nil, errors.New("导师资料异常")
	}

	result := make([]int64, num)

	for i, teacher := range teachers {
		result[i] = teacher.UserId
	}

	return result, nil
}
