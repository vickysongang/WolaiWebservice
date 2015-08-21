package main

import (
	"fmt"
	"math"
)

func GetTeacherRecommendationList(pageNum, pageCount int) POITeachers {
	teachers := QueryTeacherList(pageNum, pageCount)
	for i := range teachers {
		teachers[i].LabelList = QueryTeacherLabelByUserId(teachers[i].UserId)
	}
	return teachers
}

func GetTeacherProfile(userId, teacherId int64) *POITeacherProfile {
	teacherProfile := QueryTeacherProfile(teacherId)
	fmt.Println(teacherProfile.Rating)
	mod := math.Mod(float64(teacherId), 50)
	teacherProfile.Rating = float64(50-mod) / 10.0

	if RedisManager.redisError == nil {
		teacherProfile.HasFollowed = RedisManager.HasFollowedUser(userId, teacherId)
	}
	return teacherProfile
}
