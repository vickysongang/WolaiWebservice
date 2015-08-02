package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const DB_URL_DEV = "poi:public11223@tcp(poianalytics.mysql.rds.aliyuncs.com:3306)/wolai"

type POIDBManager struct {
	dbClient *sql.DB
}

func NewPOIDBManager() POIDBManager {
	dbClient, _ := sql.Open("mysql", DB_URL_DEV)
	return POIDBManager{dbClient: dbClient}
}

func (dbm *POIDBManager) GetUserById(userId int64) *POIUser {
	var nicknameNS sql.NullString
	var avatarNS sql.NullString
	var gender int64
	var accessRight int64

	stmtQuery, err := dbm.dbClient.Prepare(
		`SELECT nickname, avatar, gender, access_right FROM users WHERE id = ?`)
	defer stmtQuery.Close()

	if err != nil {
		panic(err.Error())
	}

	rowUser := stmtQuery.QueryRow(userId)
	err = rowUser.Scan(&nicknameNS, &avatarNS, &gender, &accessRight)

	if err == sql.ErrNoRows {
		return nil
	}

	nickname := ""
	if nicknameNS.Valid {
		nickname = nicknameNS.String
	}

	avatar := ""
	if avatarNS.Valid {
		avatar = avatarNS.String
	}

	user := NewPOIUser(int64(userId), nickname, avatar, gender, accessRight)
	return &user
}

func (dbm *POIDBManager) GetUserByPhone(phone string) *POIUser {
	var userId int64
	var nicknameNS sql.NullString
	var avatarNS sql.NullString
	var gender int64
	var accessRight int64

	stmtQuery, err := dbm.dbClient.Prepare(
		`SELECT id, nickname, avatar, gender, access_right FROM users WHERE phone = ?`)
	defer stmtQuery.Close()

	if err != nil {
		panic(err.Error())
	}

	rowUser := stmtQuery.QueryRow(phone)
	err = rowUser.Scan(&userId, &nicknameNS, &avatarNS, &gender, &accessRight)
	if err == sql.ErrNoRows {
		return nil
	}

	nickname := ""
	if nicknameNS.Valid {
		nickname = nicknameNS.String
	}

	avatar := ""
	if avatarNS.Valid {
		avatar = avatarNS.String
	}

	user := NewPOIUser(userId, nickname, avatar, gender, accessRight)
	return &user
}

func (dbm *POIDBManager) InsertUser(phone string) int64 {
	stmtInsert, err := dbm.dbClient.Prepare(
		`INSERT INTO users(phone) VALUES(?)`)
	defer stmtInsert.Close()
	if err != nil {
		panic(err.Error())
	}

	result, err := stmtInsert.Exec(phone)
	if err != nil {
		panic(err.Error())
	}

	id, _ := result.LastInsertId()

	return id
}

func (dbm *POIDBManager) UpdateUserInfo(userId int64, nickname string, avatar string, gender int64) *POIUser {
	stmtUpdate, err := dbm.dbClient.Prepare(
		`UPDATE users SET nickname = ?, avatar = ?, gender = ? WHERE id = ?`)
	defer stmtUpdate.Close()

	if err != nil {
		panic(err.Error())
	}

	_, err = stmtUpdate.Exec(nickname, avatar, gender, userId)
	if err != nil {
		panic(err.Error())
	}

	user := NewPOIUser(userId, nickname, avatar, gender, 3)
	return &user
}

func (dbm *POIDBManager) InsertUserOauth(userId int64, qqOpenId string) {
	stmtInsert, err := dbm.dbClient.Prepare(
		`INSERT INTO user_oauth(user_id, open_id_qq) VALUES(?, ?)`)
	defer stmtInsert.Close()
	if err != nil {
		panic(err.Error())
	}

	_, err = stmtInsert.Exec(userId, qqOpenId)
	if err != nil {
		panic(err.Error())
	}
}

func (dbm *POIDBManager) QueryUserByQQOpenId(qqOpenId string) int64 {
	stmtQuery, err := dbm.dbClient.Prepare(
		`SELECT user_id FROM user_oauth WHERE open_id_qq = ?`)
	defer stmtQuery.Close()
	if err != nil {
		panic(err.Error())
	}

	userRow := stmtQuery.QueryRow(qqOpenId)
	var userId int64
	err = userRow.Scan(&userId)
	if err == sql.ErrNoRows {
		return -1
	}

	return userId
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
