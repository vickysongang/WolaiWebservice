package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
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

func (dbm *POIDBManager) QueryTeacherList() POITeachers {
	stmtQuery, err := dbm.dbClient.Prepare(
		`
		SELECT users.id, users.nickname, users.avatar, users.gender, 
			teacher_profile.service_time, school.name, department.name
		FROM users, teacher_profile, school, department
		WHERE users.access_right = 2 
			AND users.id = teacher_profile.user_id 
			AND teacher_profile.school_id = school.id 
			AND teacher_profile.department_id = department.id`)
	if err != nil {
		panic(err.Error())
	}
	defer stmtQuery.Close()

	rows, err := stmtQuery.Query()
	if err != nil {
		panic(err.Error())
	}

	var userId int64
	var nicknameNS sql.NullString
	var avatarNS sql.NullString
	var gender int64
	var serviceTime int64
	var school string
	var department string
	teachers := make(POITeachers, 0)

	for rows.Next() {
		err = rows.Scan(&userId, &nicknameNS, &avatarNS, &gender, &serviceTime, &school, &department)

		nickname := ""
		if nicknameNS.Valid {
			nickname = nicknameNS.String
		}

		avatar := ""
		if avatarNS.Valid {
			avatar = avatarNS.String
		}

		teachers = append(teachers, POITeacher{POIUser: POIUser{UserId: userId, Nickname: nickname, Avatar: avatar, Gender: gender},
			ServiceTime: serviceTime, School: school, Department: department})
	}

	return teachers
}

func (dbm *POIDBManager) QueryTeacher(userId int64) *POITeacher {
	stmtQuery, err := dbm.dbClient.Prepare(
		`
		SELECT users.nickname, users.avatar, users.gender, 
			teacher_profile.service_time, teacher_profile.price_per_hour, 
			school.name, department.name
		FROM users, teacher_profile, school, department
		WHERE users.id = ?
			AND users.id = teacher_profile.user_id
			AND teacher_profile.school_id = school.id 
			AND teacher_profile.department_id = department.id`)
	if err != nil {
		panic(err.Error())
	}
	defer stmtQuery.Close()

	row := stmtQuery.QueryRow(userId)
	var nicknameNS sql.NullString
	var avatarNS sql.NullString
	var gender int64
	var serviceTime int64
	var price int64
	var school string
	var department string

	err = row.Scan(&nicknameNS, &avatarNS, &gender, &serviceTime, &price, &school, &department)
	if err != nil {
		panic(err.Error())
	}

	nickname := ""
	if nicknameNS.Valid {
		nickname = nicknameNS.String
	}

	avatar := ""
	if avatarNS.Valid {
		avatar = avatarNS.String
	}

	teacher := POITeacher{POIUser: POIUser{UserId: userId, Nickname: nickname, Avatar: avatar, Gender: gender},
		ServiceTime: serviceTime, School: school, Department: department, PricePerHour: price}

	return &teacher
}

func (dbm *POIDBManager) QueryTeacherProfile(userId int64) POITeacherProfile {
	stmtQuery, err := dbm.dbClient.Prepare(
		`
		SELECT users.nickname, users.avatar, users.gender, 
			teacher_profile.service_time, teacher_profile.intro, teacher_profile.price_per_hour, 
			school.name, department.name
		FROM users, teacher_profile, school, department
		WHERE users.id = ?
			AND users.id = teacher_profile.user_id
			AND teacher_profile.school_id = school.id 
			AND teacher_profile.department_id = department.id`)
	if err != nil {
		panic(err.Error())
	}
	defer stmtQuery.Close()

	row := stmtQuery.QueryRow(userId)
	var nicknameNS sql.NullString
	var avatarNS sql.NullString
	var gender int64
	var serviceTime int64
	var intro string
	var price int64
	var school string
	var department string

	err = row.Scan(&nicknameNS, &avatarNS, &gender, &serviceTime, &intro, &price, &school, &department)
	if err != nil {
		panic(err.Error())
	}

	nickname := ""
	if nicknameNS.Valid {
		nickname = nicknameNS.String
	}

	avatar := ""
	if avatarNS.Valid {
		avatar = avatarNS.String
	}

	teacherProfile := POITeacherProfile{POITeacher: POITeacher{POIUser: POIUser{UserId: userId, Nickname: nickname, Avatar: avatar, Gender: gender},
		ServiceTime: serviceTime, School: school, Department: department, PricePerHour: price}, Intro: intro}

	return teacherProfile
}

func (dbm *POIDBManager) QueryTeacherLabelById(userId int64) []string {
	stmtQuery, err := dbm.dbClient.Prepare(
		`
		SELECT teacher_label.name FROM teacher_to_label, teacher_label
		WHERE teacher_to_label.user_id = ? AND teacher_to_label.label_id = teacher_label.id`)
	if err != nil {
		panic(err.Error())
	}
	defer stmtQuery.Close()

	rows, err := stmtQuery.Query(userId)
	if err != nil {
		panic(err.Error())
	}

	var label string
	labels := make([]string, 0)

	for rows.Next() {
		err = rows.Scan(&label)
		if err != nil {
			panic(err.Error())
		}

		labels = append(labels, label)
	}

	return labels
}

func (dbm *POIDBManager) QueryTeacherSubjectById(userId int64) POITeacherSubjects {
	stmtQuery, err := dbm.dbClient.Prepare(
		`
		SELECT subject.name, teacher_to_subject.description FROM teacher_to_subject, subject
		WHERE teacher_to_subject.user_id = ? AND teacher_to_subject.subject_id = subject.id`)
	if err != nil {
		panic(err.Error())
	}
	defer stmtQuery.Close()

	rows, err := stmtQuery.Query(userId)
	if err != nil {
		panic(err.Error())
	}

	var name string
	var description string
	subjects := make(POITeacherSubjects, 0)

	for rows.Next() {
		err = rows.Scan(&name, &description)
		if err != nil {
			panic(err.Error())
		}

		subjects = append(subjects, POITeacherSubject{SubjectName: name, Description: description})
	}

	return subjects
}
