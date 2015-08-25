package main

import "math"

func GetTeacherRecommendationList(userId, pageNum, pageCount int64) (POITeachers, error) {
	teachers, err := QueryTeacherList(pageNum, pageCount)
	if err != nil {
		return nil, err
	}
	for i := range teachers {
		teachers[i].LabelList = QueryTeacherLabelByUserId(teachers[i].UserId)
		teachers[i].HasFollowed = RedisManager.HasFollowedUser(userId, teachers[i].UserId)
	}
	return teachers, nil
}

func GetTeacherProfile(userId, teacherId int64) (*POITeacherProfile, error) {
	teacherProfile, err := QueryTeacherProfile(teacherId)
	if err != nil {
		return nil, err
	}
	mod := math.Mod(float64(teacherId), 50)
	teacherProfile.Rating = float64(50-mod) / 10.0

	if RedisManager.redisError == nil {
		teacherProfile.HasFollowed = RedisManager.HasFollowedUser(userId, teacherId)
	}
	return teacherProfile, err
}
