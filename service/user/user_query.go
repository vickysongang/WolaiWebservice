package user

import (
	"errors"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"

	"WolaiWebservice/models"
)

func QueryUserByPhone(phone string) (*models.User, error) {
	var err error

	o := orm.NewOrm()

	var user models.User
	err = o.QueryTable(new(models.User).TableName()).
		Filter("phone", phone).
		One(&user)
	if err != nil {
		seelog.Errorf("%s | Phone: %s", err.Error(), phone)
		return nil, errors.New("该手机号未注册")
	}

	return &user, nil
}

func QueryUserByKeyword(keyword string, page, count int64) ([]*models.User, error) {
	var err error

	o := orm.NewOrm()

	cond := orm.NewCondition()
	cond1 := cond.And("access_right__in", 2, 3).And("status", 0)
	cond2 := cond.Or("nickname__icontains", keyword).Or("phone__icontains", keyword)
	condFin := cond.AndCond(cond1).AndCond(cond2)

	var users []*models.User
	_, err = o.QueryTable(new(models.User).TableName()).
		SetCond(condFin).
		Limit(count).Offset(page * count).
		All(&users)
	if err != nil {
		return nil, errors.New("没有符合条件的查询结果")
	}

	return users, nil
}

func QueryUserByAccessRight(accessRight, page, count int64) ([]*models.User, error) {
	var err error

	o := orm.NewOrm()

	var users []*models.User
	_, err = o.QueryTable(new(models.User).TableName()).
		Filter("access_right", accessRight).
		Limit(count).Offset(page * count).
		All(&users)
	if err != nil {
		return nil, errors.New("没有符合条件的查询结果")
	}

	return users, nil
}

func QueryAssistants() ([]int64, error) {
	var err error
	o := orm.NewOrm()
	var teacherProfiles []models.TeacherProfile
	_, err = o.QueryTable(new(models.TeacherProfile).TableName()).Filter("certify_flag", "Y").OrderBy("user_id").
		All(&teacherProfiles, "UserId")
	userIds := make([]int64, 0)
	if err == nil {
		for _, profile := range teacherProfiles {
			userIds = append(userIds, profile.UserId)
		}
	}
	return userIds, err
}
