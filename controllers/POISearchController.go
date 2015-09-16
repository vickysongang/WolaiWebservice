// POISearchController
package controllers

import (
	"github.com/cihub/seelog"
	"github.com/tmhenry/POIWolaiWebService/managers"
	"github.com/tmhenry/POIWolaiWebService/models"
)

func SearchTeacher(userId int64, keyword string, pageNum, pageCount int64) (*models.POITeachers, error) {
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
				UserId:   teacherModel.Id,
				Nickname: teacherModel.Nickname,
				Phone:    teacherModel.Phone,
				Avatar:   teacherModel.Avatar,
				Gender:   teacherModel.Gender},
			ServiceTime:      teacherModel.ServiceTime,
			School:           teacherModel.SchoolName,
			Department:       teacherModel.DeptName,
			PricePerHour:     teacherModel.PricePerHour,
			RealPricePerHour: teacherModel.RealPricePerHour}
		lablList := models.QueryTeacherLabelByUserId(teacherModel.Id)
		teacher.LabelList = lablList
		if managers.RedisManager.RedisError == nil {
			teacher.HasFollowed = managers.RedisManager.HasFollowedUser(0, teacherModel.Id)
		}
		teachers = append(teachers, teacher)
	}
	return &teachers, nil
}
