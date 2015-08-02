package main

import (
	_ "encoding/json"
	"math"
)

type POITeacher struct {
	POIUser
	School       string   `json:"school"`
	Department   string   `json:"department"`
	ServiceTime  int64    `json:"serviceTime"`
	LabelList    []string `json:"labelList,omitempty"`
	PricePerHour int64    `json:"pricePerHour"`
}
type POITeachers []POITeacher

type POITeacherSubject struct {
	//SubjectId   int64  `json:"subjectId"`
	SubjectName string `json:"subjectName"`
	Description string `json:"description"`
}
type POITeacherSubjects []POITeacherSubject

type POITeacherResume struct {
	Start int64  `json:"start"`
	Stop  int64  `json:"stop"`
	Name  string `json:"name"`
}
type POITeacherResumes []POITeacherResume

type POITeacherProfile struct {
	POITeacher
	Rating        float64            `json:"rating"`
	SubjectList   POITeacherSubjects `json:"subjectList"`
	EducationList POITeacherResumes  `json:"eduList"`
	Intro         string             `json:"intro"`
	HasFollowed   bool               `json:"hasFollowed"`
}

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
