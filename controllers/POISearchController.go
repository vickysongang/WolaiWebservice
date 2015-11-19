// POISearchController
package controllers

import (
	"WolaiWebservice/models"
	"WolaiWebservice/redis"

	"github.com/cihub/seelog"
)

func SearchTeachers(userId int64, keyword string, pageNum, pageCount int64) (*models.POITeachers, error) {
	teacherModels, err := models.QueryTeachersByCond(userId, keyword, pageNum, pageCount)
	if err != nil {
		seelog.Error("keyword:", keyword, " ", err.Error())
		return nil, err
	}
	teachers := make(models.POITeachers, 0)
	for i := range teacherModels {
		teacherModel := teacherModels[i]
		teacher := models.POITeacher{
			POIUser: models.POIUser{
				UserId:      teacherModel.Id,
				Nickname:    teacherModel.Nickname,
				Phone:       teacherModel.Phone,
				Avatar:      teacherModel.Avatar,
				AccessRight: teacherModel.AccessRight,
				Gender:      teacherModel.Gender},
			ServiceTime:      teacherModel.ServiceTime,
			School:           teacherModel.SchoolName,
			Department:       teacherModel.DeptName,
			PricePerHour:     teacherModel.PricePerHour,
			RealPricePerHour: teacherModel.RealPricePerHour}
		lablList := models.QueryTeacherLabelByUserId(teacherModel.Id)
		teacher.LabelList = lablList
		if redis.RedisManager.RedisError == nil {
			teacher.HasFollowed = redis.RedisManager.HasFollowedUser(0, teacherModel.Id)
		}
		teachers = append(teachers, teacher)
	}
	return &teachers, nil
}

func SearchUsers(userId int64, keyword string, pageNum, pageCount int64) (*models.POITeachers, error) {
	teacherModels, err := models.QueryUsersByCond(userId, keyword, pageNum, pageCount)
	if err != nil {
		seelog.Error("keyword:", keyword, " ", err.Error())
		return nil, err
	}
	teachers := make(models.POITeachers, 0)
	for i := range teacherModels {
		teacherModel := teacherModels[i]
		teacher := models.POITeacher{
			POIUser: models.POIUser{
				UserId:      teacherModel.Id,
				Nickname:    teacherModel.Nickname,
				Phone:       teacherModel.Phone,
				Avatar:      teacherModel.Avatar,
				AccessRight: teacherModel.AccessRight,
				Gender:      teacherModel.Gender},
			ServiceTime:      teacherModel.ServiceTime,
			School:           teacherModel.SchoolName,
			Department:       teacherModel.DeptName,
			PricePerHour:     teacherModel.PricePerHour,
			RealPricePerHour: teacherModel.RealPricePerHour}
		lablList := models.QueryTeacherLabelByUserId(teacherModel.Id)
		teacher.LabelList = lablList
		if redis.RedisManager.RedisError == nil {
			teacher.HasFollowed = redis.RedisManager.HasFollowedUser(0, teacherModel.Id)
		}
		teachers = append(teachers, teacher)
	}
	return &teachers, nil
}
