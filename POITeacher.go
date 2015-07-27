package main

import (
	_ "encoding/json"
)

type POITeacher struct {
	POIUser
	School      string   `json:"school"`
	Department  string   `json:"department"`
	ServiceTime int64    `json:"serviceTime"`
	LabelList   []string `json:"labelList"`
}
type POITeachers []POITeacher

type POITeacherSubject struct {
	SubjectId   int64  `json:"subjectId"`
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
	SpecialtyList POITeacherSubjects `json:"specialtyList"`
	EducationList POITeacherResumes  `json:"eduList"`
	Intro         string             `json:"intro"`
}

func GetTeacherRecommendationList() POITeachers {
	teachers := DbManager.QueryTeacherList()

	for i := range teachers {
		teachers[i].LabelList = DbManager.QueryTeacherLabelById(teachers[i].UserId)
	}

	return teachers
}
