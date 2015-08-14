package main

import (
	"math"
)

func GetTeacherRecommendationList(page int64) POITeachers {
	teachers := QueryTeacherList()

	for i := range teachers {
		teachers[i].LabelList = QueryTeacherLabelById(teachers[i].UserId)
	}

	return teachers
}

func GetTeacherProfile(userId, teacherId int64) POITeacherProfile {
	teacherProfile := QueryTeacherProfile(teacherId)

	teacherProfile.LabelList = QueryTeacherLabelById(teacherId)

	teacherProfile.SubjectList = QueryTeacherSubjectById(teacherId)

	teacherProfile.EducationList = QueryTeacherResumeById(teacherId)

	mod := math.Mod(float64(teacherId), 50)

	teacherProfile.Rating = float64(50-mod) / 10.0
	
	if RedisManager.redisError == nil {
		teacherProfile.HasFollowed = RedisManager.HasFollowedUser(userId, teacherId)
	}
	return teacherProfile
}
