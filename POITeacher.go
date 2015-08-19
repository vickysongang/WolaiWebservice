package main

import (
	"encoding/json"
	"fmt"
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
	SubjectName string `json:"subjectName" orm:"column(name)"`
	Description string `json:"description"`
}
type POITeacherSubjects []POITeacherSubject

type POITeacherResume struct {
	Id     int64  `json:"-" orm:"pk"`
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

type POITeacherProfileModel struct {
	UserId           int64  `json:"-" orm:"pk"`
	SchoolId         int64  `json:"schoolId"`
	DepartmentId     int64  `json:"departmentId"`
	Intro            string `json:"intro"`
	PricePerHour     int64  `json:"pricePerHour"`
	RealPricePerHour int64  `json:"realPricePerHour"`
}

type POITeacherLabel struct {
	Id   int64  `json:"id" orm:"pk"`
	Name string `json:"name"`
}

type POITeacherToLabel struct {
	Id      int64 `json:"id" orm:"pk"`
	UserId  int64 `json:"userId"`
	LabelId int64 `json:"labelId"`
}

type POITeacherToSubject struct {
	Id          int64  `json:"-" orm:"pk"`
	UserId      int64  `json:"-"`
	SubjectId   int64  `json:"subjectId"`
	Description string `json:"description"`
}

type POITeacherInfo struct {
	POIUser                `json:"teacherInfo"`
	LabelList              []string `json:"labelList,omitempty"`
	POITeacherResume       `json:"resumeInfo"`
	POITeacherToSubject    `json:"subjectInfo"`
	POITeacherProfileModel `json:"profileInfo"`
}

func (r *POITeacherResume) TableName() string {
	return "teacher_to_resume"
}

func (p *POITeacherProfileModel) TableName() string {
	return "teacher_profile"
}

func (ttl *POITeacherLabel) TableName() string {
	return "teacher_label"
}

func (ttl *POITeacherToLabel) TableName() string {
	return "teacher_to_label"
}

func (tts *POITeacherToSubject) TableName() string {
	return "teacher_to_subject"
}

func init() {
	orm.RegisterModel(new(POITeacherResume), new(POITeacherLabel), new(POITeacherToLabel), new(POITeacherToSubject), new(POITeacherProfileModel))
}

func QueryTeacherList(pageNum, pageCount int) POITeachers {
	start := pageNum * pageCount
	teachers := make(POITeachers, 0)
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("users.id, users.nickname, users.avatar, users.gender,teacher_profile.service_time, school.name school_name, department.name dept_name").
		From("users").InnerJoin("teacher_profile").On("users.id = teacher_profile.user_id").InnerJoin("school").
		On("teacher_profile.school_id = school.id").InnerJoin("department").On("teacher_profile.department_id = department.id").
		Where("users.access_right = 2 and users.status = 0").Limit(pageCount).Offset(start)
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

func QueryTeacherLabelByName(name string) *POITeacherLabel {
	o := orm.NewOrm()
	teacherLabel := POITeacherLabel{}
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("id,name").From("teacher_label").Where("name = ?")
	sql := qb.String()
	err := o.Raw(sql, name).QueryRow(&teacherLabel)
	if err != nil {
		return nil
	}
	return &teacherLabel
}

func InsertTeacherLabel(name string) int64 {
	label := QueryTeacherLabelByName(name)
	if label != nil {
		return label.Id
	}
	teacherLabel := POITeacherLabel{Name: name}
	o := orm.NewOrm()
	id, err := o.Insert(&teacherLabel)
	if err != nil {
		return 0
	}
	return id
}

func QueryTeacherToLabel(userId, labelId int64) *POITeacherToLabel {
	o := orm.NewOrm()
	ttl := POITeacherToLabel{}
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("id,user_id,label_id").From("teacher_to_lable").Where("user_id = ? and label_id = ?")
	sql := qb.String()
	err := o.Raw(sql, userId, labelId).QueryRow(&ttl)
	if err != nil {
		return nil
	}
	return &ttl
}

func InsertTeacherToLabel(teacherLabel *POITeacherToLabel) int64 {
	o := orm.NewOrm()
	id, err := o.Insert(teacherLabel)
	if err != nil {
		return 0
	}
	return id
}

func QueryTeacherToSubject(userId, subjectId int64) *POITeacherToSubject {
	o := orm.NewOrm()
	tts := POITeacherToSubject{}
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select("id,user_id,subject_id,description").From("teacher_to_subject").Where("user_id = ? and subject_id = ?")
	sql := qb.String()
	err := o.Raw(sql, userId, subjectId).QueryRow(&tts)
	if err != nil {
		return nil
	}
	return &tts
}

func InsertTeacherToSubject(teacherSubject *POITeacherToSubject) int64 {
	o := orm.NewOrm()
	id, err := o.Insert(teacherSubject)
	if err != nil {
		return 0
	}
	return id
}

func InsertTeacherToResume(resume *POITeacherResume) int64 {
	o := orm.NewOrm()
	id, err := o.Insert(resume)
	if err != nil {
		return 0
	}
	return id
}

func InsertTeacherProfile(profile *POITeacherProfileModel) int64 {
	o := orm.NewOrm()
	id, err := o.Insert(profile)
	if err != nil {
		return 0
	}
	return id
}

func GenerateTeacherJson() string {
	teacher := POITeacherInfo{}
	//user
	teacher.Phone = "17090889054"
	teacher.Avatar = "avataraaaaa"
	teacher.Nickname = "vicky"
	teacher.Gender = 0

	//profile
	teacher.SchoolId = 1
	teacher.DepartmentId = 1001
	teacher.Intro = "我们都是好孩子"
	teacher.PricePerHour = 1000
	teacher.RealPricePerHour = 10000

	//label
	labelList := []string{"1", "2", "3"}
	teacher.LabelList = labelList

	//resume
	teacher.Start = 2005
	teacher.Stop = 2009
	teacher.Name = "不知道这个是啥"

	//subject
	teacher.SubjectId = 11
	teacher.Description = "科目描述"

	teacherInfo, _ := json.Marshal(&teacher)
	return string(teacherInfo)
}

func InsertTeacher(teacherInfo string) *POITeacherInfo {
	teacher := POITeacherInfo{}
	json.Unmarshal([]byte(teacherInfo), &teacher)
	//插入用户基本信息
	user := POIUser{}
	user.AccessRight = 2
	user.Status = 0
	user.Avatar = teacher.Avatar
	user.Gender = teacher.Gender
	user.Nickname = teacher.Nickname
	user.Phone = teacher.Phone
	userId := InsertPOIUser(&user)
	if userId == 0 {
		return nil
	}
	fmt.Println("userId:", userId)
	//处理Label信息
	fmt.Println("labelList:", teacher.LabelList)
	labelList := teacher.LabelList
	for _, label := range labelList {
		teacherLabel := QueryTeacherLabelByName(label)
		var labelId int64
		//如果Label已经存在则直接使用，否则先将Label插入数据库后再使用
		if teacherLabel == nil {
			labelId = InsertTeacherLabel(label)
		} else {
			labelId = teacherLabel.Id
		}
		teacherToLabel := POITeacherToLabel{UserId: userId, LabelId: labelId}
		InsertTeacherToLabel(&teacherToLabel)
	}
	//处理科目信息
	teacherSubject := POITeacherToSubject{UserId: userId, SubjectId: teacher.SubjectId, Description: teacher.Description}
	InsertTeacherToSubject(&teacherSubject)
	//处理简历信息
	teacherResume := POITeacherResume{UserId: userId, Start: teacher.Start, Stop: teacher.Stop, Name: teacher.Name}
	InsertTeacherToResume(&teacherResume)
	//处理Profile信息
	teacherProfile := POITeacherProfileModel{UserId: userId, SchoolId: teacher.SchoolId, DepartmentId: teacher.DepartmentId,
		Intro: teacher.Intro, PricePerHour: teacher.PricePerHour, RealPricePerHour: teacher.RealPricePerHour}
	InsertTeacherProfile(&teacherProfile)
	fmt.Println(teacher.Phone)
	return &teacher
}
