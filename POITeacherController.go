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

	mod := math.Mod(float64(teacherId), 50)

	teacherProfile.Rating = float64(50-mod) / 10.0

	resumes := make(POITeacherResumes, 2)
	resumes[0] = POITeacherResume{Start: 2008, Stop: -1, Name: "电线杆子科技大学"}
	resumes[1] = POITeacherResume{Start: 2005, Stop: 2008, Name: "马路牙子高级中学"}

	teacherProfile.EducationList = resumes

	teacherProfile.HasFollowed = RedisManager.HasFollowedUser(userId, teacherId)

	return teacherProfile
}
