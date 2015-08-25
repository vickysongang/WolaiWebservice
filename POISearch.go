// POISearch
package main

import (
	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"
)

func SearchTeacher(userId int64, keyword string, pageNum, pageCount int64) (*POITeachers, error) {
	start := pageNum * pageCount
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("users.id,users.nickname,users.phone,users.avatar, users.gender,teacher_profile.service_time, teacher_profile.price_per_hour,teacher_profile.real_price_per_hour,school.name school_name,department.name dept_name").
		From("users").InnerJoin("teacher_profile").On("users.id = teacher_profile.user_id").
		InnerJoin("school").On("teacher_profile.school_id = school.id").
		InnerJoin("department").On("teacher_profile.department_id = department.id").
		Where("users.access_right = 2 and (users.nickname like ? or users.phone like ?)").Limit(int(pageCount)).Offset(int(start))
	sql := qb.String()
	o := orm.NewOrm()
	var teacherModels POITeacherModels
	_, err := o.Raw(sql, "%"+keyword+"%", "%"+keyword+"%").QueryRows(&teacherModels)
	if err != nil {
		seelog.Error("keyword:", keyword, " ", err.Error())
		return nil, err
	}
	teachers := make(POITeachers, 0)
	for i := range teacherModels {
		teacherModel := teacherModels[i]
		teacher := POITeacher{POIUser: POIUser{UserId: teacherModel.Id, Nickname: teacherModel.Nickname, Phone: teacherModel.Phone, Avatar: teacherModel.Avatar, Gender: teacherModel.Gender},
			ServiceTime: teacherModel.ServiceTime, School: teacherModel.SchoolName, Department: teacherModel.DeptName, PricePerHour: teacherModel.PricePerHour, RealPricePerHour: teacherModel.RealPricePerHour}
		lablList := QueryTeacherLabelByUserId(teacherModel.Id)
		teacher.LabelList = lablList
		if RedisManager.redisError == nil {
			teacher.HasFollowed = RedisManager.HasFollowedUser(0, teacherModel.Id)
		}
		teachers = append(teachers, teacher)
	}
	return &teachers, nil
}
