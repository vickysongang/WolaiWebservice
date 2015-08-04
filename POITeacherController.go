package main

import (
	"math"
)

func GetTeacherRecommendationList(page int64) POITeachers {
	teachers := DbManager.QueryTeacherList()

	for i := range teachers {
		teachers[i].LabelList = DbManager.QueryTeacherLabelById(teachers[i].UserId)
	}

	return teachers
}

func GetTeacherProfile(userId, teacherId int64) POITeacherProfile {
	teacherProfile := DbManager.QueryTeacherProfile(teacherId)

	teacherProfile.LabelList = DbManager.QueryTeacherLabelById(teacherId)

	teacherProfile.SubjectList = DbManager.QueryTeacherSubjectById(teacherId)

	teacherProfile.EducationList = DbManager.QueryTeacherResumeById(teacherId)

	mod := math.Mod(float64(teacherId), 50)

	teacherProfile.Rating = float64(50-mod) / 10.0

	teacherProfile.HasFollowed = RedisManager.HasFollowedUser(userId, teacherId)

	return teacherProfile
}
