package main

import (
	"strconv"

	"github.com/astaxie/beego/orm"
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
	Id     int64  `json:"-"`
	UserId int64  `json:"-"`
	Start  int64  `json:"start"`
	Stop   int64  `json:"stop"`
	Name   string `json:"name"`
}
type POITeacherResumes []POITeacherResume

type POITeacherProfile struct {
	UserId        int64 `json:"-" orm:"pk"`
	POITeacher    `orm:"-"`
	Rating        float64            `json:"rating"`
	SubjectList   POITeacherSubjects `json:"subjectList" orm:"-"`
	EducationList POITeacherResumes  `json:"eduList" orm:"-"`
	Intro         string             `json:"intro"`
	HasFollowed   bool               `json:"hasFollowed"`
	ServiceTime   int64              `json:"-"`
}

func (r *POITeacherResume) TableName() string {
	return "teacher_to_resume"
}

func (p *POITeacherProfile) TableName() string {
	return "teacher_profile"
}

func init() {
	orm.RegisterModel(new(POITeacherResume), new(POITeacherProfile))
}

func QueryTeacherList() POITeachers {
	teachers := make(POITeachers, 0)
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("users.id, users.nickname, users.avatar, users.gender,teacher_profile.service_time, school.name school_name, department.name dept_name").
		From("users").InnerJoin("teacher_profile").On("users.id = teacher_profile.user_id").InnerJoin("school").
		On("teacher_profile.school_id = school.id").InnerJoin("department").On("teacher_profile.department_id = department.id").
		Where("users.access_right = 2")
	sql := qb.String()
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw(sql).Values(&maps)
	if err == orm.ErrNoRows {
		return nil
	}
	for i := range maps {
		teacher := maps[i]
		userIdStr, _ := teacher["id"].(string)
		nickname, _ := teacher["nickname"].(string)
		avatar, _ := teacher["avatar"].(string)
		genderStr, _ := teacher["gender"].(string)
		serviceTimeStr, _ := teacher["service_time"].(string)
		schoolName, _ := teacher["school_name"].(string)
		deptName, _ := teacher["dept_name"].(string)
		userId, _ := strconv.ParseInt(userIdStr, 10, 64)
		gender, _ := strconv.ParseInt(genderStr, 10, 64)
		serviceTime, _ := strconv.ParseInt(serviceTimeStr, 10, 64)
		teachers = append(teachers, POITeacher{POIUser: POIUser{UserId: userId, Nickname: nickname,
			Avatar: avatar, Gender: gender}, ServiceTime: serviceTime, School: schoolName,
			Department: deptName})
	}
	return teachers
}

func QueryTeacher(userId int64) *POITeacher {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("users.nickname, users.avatar, users.gender,teacher_profile.service_time, teacher_profile.price_per_hour,school.name school_name,department.name dept_name").
		From("users").InnerJoin("teacher_profile").On("users.id = teacher_profile.user_id").
		InnerJoin("school").On("teacher_profile.school_id = school.id").
		InnerJoin("department").On("teacher_profile.department_id = department.id").
		Where("users.id = ?")
	sql := qb.String()
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw(sql, userId).Values(&maps)
	if err == orm.ErrNoRows {
		return nil
	}
	teacherInfo := maps[0]
	nickname, _ := teacherInfo["nickname"].(string)
	avatar, _ := teacherInfo["avatar"].(string)
	genderStr, _ := teacherInfo["gender"].(string)
	serviceTimeStr, _ := teacherInfo["service_time"].(string)
	schoolName, _ := teacherInfo["school_name"].(string)
	deptName, _ := teacherInfo["dept_name"].(string)
	gender, _ := strconv.ParseInt(genderStr, 10, 64)
	serviceTime, _ := strconv.ParseInt(serviceTimeStr, 10, 64)
	pricePerHourStr, _ := teacherInfo["price_per_hour"].(string)
	pricePerHour, _ := strconv.ParseInt(pricePerHourStr, 10, 64)
	teacher := POITeacher{POIUser: POIUser{UserId: userId, Nickname: nickname, Avatar: avatar, Gender: gender},
		ServiceTime: serviceTime, School: schoolName, Department: deptName, PricePerHour: pricePerHour}
	return &teacher
}

func QueryTeacherProfile(userId int64) POITeacherProfile {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("users.nickname, users.avatar, users.gender,teacher_profile.service_time, teacher_profile.intro, teacher_profile.price_per_hour,school.name school_name, department.name dept_name").
		From("users").InnerJoin("teacher_profile").On("users.id = teacher_profile.user_id").
		InnerJoin("school").On("teacher_profile.school_id = school.id").
		InnerJoin("department").On("teacher_profile.department_id = department.id").
		Where("users.id = ?")
	sql := qb.String()
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw(sql, userId).Values(&maps)
	if err == orm.ErrNoRows {
		panic(err.Error())
	}
	profile := maps[0]
	nickname, _ := profile["nickname"].(string)
	avatar, _ := profile["avatar"].(string)
	genderStr, _ := profile["gender"].(string)
	serviceTimeStr, _ := profile["service_time"].(string)
	schoolName, _ := profile["school_name"].(string)
	intro, _ := profile["intro"].(string)
	deptName, _ := profile["dept_name"].(string)
	gender, _ := strconv.ParseInt(genderStr, 10, 64)
	serviceTime, _ := strconv.ParseInt(serviceTimeStr, 10, 64)
	pricePerHourStr, _ := profile["price_per_hour"].(string)
	pricePerHour, _ := strconv.ParseInt(pricePerHourStr, 10, 64)
	teacherProfile := POITeacherProfile{POITeacher: POITeacher{POIUser: POIUser{UserId: userId, Nickname: nickname, Avatar: avatar, Gender: gender},
		ServiceTime: serviceTime, School: schoolName, Department: deptName, PricePerHour: pricePerHour}, Intro: intro}
	return teacherProfile
}

func QueryTeacherLabelById(userId int64) []string {
	labels := make([]string, 0)
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("teacher_label.name").From("teacher_label").InnerJoin("teacher_to_label").
		On("teacher_to_label.label_id = teacher_label.id").Where("teacher_to_label.user_id = ?")
	o := orm.NewOrm()
	sql := qb.String()
	_, err := o.Raw(sql, userId).QueryRows(&labels)
	if err == orm.ErrNoRows {
		return nil
	}
	return labels
}

func QueryTeacherSubjectById(userId int64) POITeacherSubjects {
	subjects := make(POITeacherSubjects, 0)
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("subject.name,teacher_to_subject.description").From("teacher_to_subject").
		InnerJoin("subject").On("teacher_to_subject.subject_id = subject.id").Where("teacher_to_subject.user_id = ?")
	sql := qb.String()
	o := orm.NewOrm()
	_, err := o.Raw(sql, userId).QueryRows(&subjects)
	if err == orm.ErrNoRows {
		return nil
	}
	return subjects
}

func QueryTeacherResumeById(userId int64) POITeacherResumes {
	resumes := make(POITeacherResumes, 0)
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("id,user_id,start,stop,name").From("teacher_to_resume").Where("user_id = ?")
	sql := qb.String()
	o := orm.NewOrm()
	_, err := o.Raw(sql, userId).QueryRows(&resumes)
	if err == orm.ErrNoRows {
		return nil
	}
	return resumes
}

func UpdateTeacherServiceTime(userId int64, length int64) {
	o := orm.NewOrm()
	_, err := o.QueryTable("teacher_profile").Filter("user_id", userId).Update(orm.Params{
		"service_time": orm.ColValue(orm.Col_Add, length),
	})
	if err != nil {
		panic(err.Error())
	}
}

