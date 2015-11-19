package controllers

import (
	"encoding/json"
	"fmt"

	"WolaiWebService/models"
	"WolaiWebService/redis"

	"github.com/cihub/seelog"
)

func GetTeacherRecommendationList(userId, pageNum, pageCount int64) (models.POITeachers, error) {
	teachers, err := models.QueryTeacherList(pageNum, pageCount)
	if err != nil {
		return nil, err
	}
	for i := range teachers {
		teachers[i].LabelList = models.QueryTeacherLabelByUserId(teachers[i].UserId)
		teachers[i].HasFollowed = redis.RedisManager.HasFollowedUser(userId, teachers[i].UserId)
	}
	return teachers, nil
}

func GetSupportAndTeacherList(userId, pageNum, pageCount int64) (models.POITeachers, error) {
	teachers, err := models.QueryTeacherList(pageNum, pageCount)
	if err != nil {
		return nil, err
	}
	if pageNum == 0 {
		supports, _ := models.QuerySupportList()
		teachers = append(supports, teachers...)
	}
	for i := range teachers {
		teachers[i].LabelList = models.QueryTeacherLabelByUserId(teachers[i].UserId)
		teachers[i].HasFollowed = redis.RedisManager.HasFollowedUser(userId, teachers[i].UserId)
	}
	return teachers, nil
}

func GetTeacherProfile(userId, teacherId int64) (*models.POITeacherProfile, error) {
	teacherProfile, err := models.QueryTeacherProfile(teacherId)
	if err != nil {
		return nil, err
	}
	teacherProfile.Rating = 5.0

	if redis.RedisManager.RedisError == nil {
		teacherProfile.HasFollowed = redis.RedisManager.HasFollowedUser(userId, teacherId)
	}
	return teacherProfile, nil
}

//teacherInfo为json格式
func InsertTeacher(teacherInfo string) (models.POITeacherInfos, error) {
	fmt.Println(teacherInfo)
	var teachers models.POITeacherInfos
	err := json.Unmarshal([]byte(teacherInfo), &teachers)
	if err != nil {
		seelog.Error("teacherInfo:", teacherInfo, " ", err.Error())
		return teachers, err
	}
	for i := range teachers {
		teacher := teachers[i]
		var userId int64
		oldUser := models.QueryUserByPhone(teacher.Phone)
		if oldUser == nil {
			//插入用户基本信息
			user := models.POIUser{}
			user.AccessRight = 2
			user.Status = 1
			user.Avatar = teacher.Avatar
			user.Gender = teacher.Gender
			user.Nickname = teacher.Nickname
			user.Phone = teacher.Phone

			userId, err = models.InsertPOIUser(&user)
			if err != nil {
				return nil, err
			}
		} else {
			if oldUser.AccessRight == 3 {
				userInfo := map[string]interface{}{
					"AccessRight": 2,
					"Status":      1,
					"Avatar":      teacher.Avatar,
					"Gender":      teacher.Gender,
					"Nickname":    teacher.Nickname,
				}
				userId = oldUser.UserId
				models.UpdateUserInfo(userId, userInfo)
			}
		}
		teacher.POIUser.UserId = userId
		//处理Label信息
		labelList := teacher.LabelList
		for _, label := range labelList {
			teacherLabel := models.QueryTeacherLabelByName(label)
			var labelId int64
			//如果Label已经存在则直接使用，否则先将Label插入数据库后再使用
			if teacherLabel == nil {
				labelId = models.InsertTeacherLabel(label)
			} else {
				labelId = teacherLabel.Id
			}
			teacherToLabel := models.POITeacherToLabel{
				UserId:  userId,
				LabelId: labelId}
			models.InsertTeacherToLabel(&teacherToLabel)
		}
		//处理科目信息
		subjectList := teacher.SubjectInfo
		for _, subject := range subjectList {
			teacherSubject := models.POITeacherToSubject{
				UserId:      userId,
				SubjectId:   subject.SubjectId,
				Description: subject.Description}
			models.InsertTeacherToSubject(&teacherSubject)
		}
		//处理简历信息
		resumeList := teacher.ResumeInfo
		for _, resume := range resumeList {
			teacherResume := models.POITeacherResume{
				UserId: userId,
				Start:  resume.Start,
				Stop:   resume.Stop,
				Name:   resume.Name}
			models.InsertTeacherToResume(&teacherResume)
		}
		//处理Profile信息
		teacherProfile := models.POITeacherProfileModel{
			UserId:           userId,
			SchoolId:         teacher.SchoolId,
			DepartmentId:     teacher.DepartmentId,
			Intro:            teacher.Intro,
			PricePerHour:     teacher.PricePerHour,
			RealPricePerHour: teacher.RealPricePerHour,
			ServiceTime:      teacher.ServiceTime}
		models.InsertTeacherProfile(&teacherProfile)
	}
	return teachers, nil
}
