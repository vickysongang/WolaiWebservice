package models

import (
	"WolaiWebservice/utils"

	"github.com/astaxie/beego/orm"
	seelog "github.com/cihub/seelog"
	_ "github.com/go-sql-driver/mysql"
)

const (
	USER_WOLAI_TEAM = 1003
)

type POITeacher struct {
	POIUser
	School           string   `json:"school"`
	Department       string   `json:"department"`
	ServiceTime      int64    `json:"serviceTime"`
	LabelList        []string `json:"labelList,omitempty"`
	PricePerHour     int64    `json:"pricePerHour"`
	RealPricePerHour int64    `json:"realPricePerHour"`
	HasFollowed      bool     `json:"hasFollowed"`
}
type POITeachers []POITeacher

//老师科目信息结构体
type POITeacherSubject struct {
	SubjectName string `json:"subjectName" orm:"column(name)"`
	Description string `json:"description"`
}
type POITeacherSubjects []POITeacherSubject

//老师简历结构体
type POITeacherResume struct {
	Id     int64  `json:"-" orm:"pk"`
	UserId int64  `json:"-"`
	Start  int64  `json:"start"`
	Stop   int64  `json:"stop"`
	Name   string `json:"name"`
}
type POITeacherResumes []POITeacherResume

//老师详细信息结构体
type POITeacherProfile struct {
	UserId        int64 `json:"-" orm:"pk"`
	POITeacher    `orm:"-"`
	Rating        float64            `json:"rating"`
	SubjectList   POITeacherSubjects `json:"subjectList" orm:"-"`
	EducationList POITeacherResumes  `json:"eduList" orm:"-"`
	Intro         string             `json:"intro"`
	Extra         string             `json:"extra"`
	ServiceTime   int64              `json:"-"`
}

//老师详细信息结构体，字段完全与数据库对应
type POITeacherProfileModel struct {
	UserId           int64   `json:"-" orm:"pk"`
	SchoolId         int64   `json:"schoolId"`
	DepartmentId     int64   `json:"departmentId"`
	Intro            string  `json:"intro"`
	Extra            string  `json:"extra"`
	PricePerHour     int64   `json:"pricePerHour"`
	RealPricePerHour int64   `json:"realPricePerHour"`
	ServiceTime      int64   `json:"serviceTime"`
	Rating           float64 `json:"-"`
}

//老师标签结构体，用户读取和维护老师标签信息
type POITeacherLabel struct {
	Id   int64  `json:"id" orm:"pk"`
	Name string `json:"name"`
}

//老师与标签对应关系的结构体，用于读取和维护老师与标签的对应关系
type POITeacherToLabel struct {
	Id      int64 `json:"id" orm:"pk"`
	UserId  int64 `json:"userId"`
	LabelId int64 `json:"labelId"`
}

//老师与科目对应关系的结构体，用户读取和维护老师与科目的对应关系
type POITeacherToSubject struct {
	Id          int64  `json:"-" orm:"pk"`
	UserId      int64  `json:"-"`
	SubjectId   int64  `json:"subjectId"`
	Description string `json:"description"`
}

//老师的完整信息结构体，用于解析维护老师信息时从客户端传过来的json字符串
type POITeacherInfo struct {
	POIUser                `json:"teacherInfo"`
	LabelList              []string              `json:"labelList,omitempty"`
	ResumeInfo             []POITeacherResume    `json:"resumeInfo"`
	SubjectInfo            []POITeacherToSubject `json:"subjectInfo"`
	POITeacherProfileModel `json:"profileInfo"`
}

//老师信息结构体，用于多表关联查询时接受各字段的值
type POITeacherModel struct {
	Id               int64
	Nickname         string
	Avatar           string
	Gender           int64
	ServiceTime      int64
	PricePerHour     int64
	RealPricePerHour int64
	SchoolName       string
	DeptName         string
	Intro            string
	Extra            string
	Rating           float64
	Phone            string
	AccessRight      int64
}

type POITeacherModels []POITeacherModel

type POITeacherInfos []POITeacherInfo

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

/*
 * 分页查询老师列表
 */
func QueryTeacherList(pageNum, pageCount int64) (POITeachers, error) {
	start := pageNum * pageCount
	teachers := make(POITeachers, 0)
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("users.id, users.nickname, users.avatar, users.gender,users.access_right,teacher_profile.service_time,teacher_profile.price_per_hour," +
		"teacher_profile.real_price_per_hour,school.name school_name, department.name dept_name").
		From("users").LeftJoin("teacher_profile").On("users.id = teacher_profile.user_id").LeftJoin("school").
		On("teacher_profile.school_id = school.id").LeftJoin("department").On("teacher_profile.department_id = department.id").
		Where("users.access_right = 2 and users.status = 0").OrderBy("teacher_profile.service_time").Desc().
		Limit(int(pageCount)).Offset(int(start))
	sql := qb.String()
	o := orm.NewOrm()
	var teacherModels POITeacherModels
	_, err := o.Raw(sql).QueryRows(&teacherModels)
	if err != nil {
		seelog.Error(err.Error())
		return nil, err
	}
	for i := range teacherModels {
		teacher := teacherModels[i]
		teachers = append(teachers, POITeacher{
			POIUser: POIUser{
				UserId:      teacher.Id,
				Nickname:    teacher.Nickname,
				Avatar:      teacher.Avatar,
				AccessRight: teacher.AccessRight,
				Gender:      teacher.Gender},
			ServiceTime:      teacher.ServiceTime,
			School:           teacher.SchoolName,
			Department:       teacher.DeptName,
			PricePerHour:     teacher.PricePerHour,
			RealPricePerHour: teacher.RealPricePerHour})
	}
	return teachers, nil
}

/*
 * 查询我来客服和我来团队
 */
func QuerySupportList() (POITeachers, error) {
	teachers := make(POITeachers, 0)
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("users.id, users.nickname, users.avatar, users.gender,users.access_right,teacher_profile.service_time,teacher_profile.price_per_hour," +
		"teacher_profile.real_price_per_hour,school.name school_name, department.name dept_name").
		From("users").LeftJoin("teacher_profile").On("users.id = teacher_profile.user_id").LeftJoin("school").
		On("teacher_profile.school_id = school.id").LeftJoin("department").On("teacher_profile.department_id = department.id").
		Where("users.id = ? or users.access_right = 1")
	sql := qb.String()
	o := orm.NewOrm()
	var teacherModels POITeacherModels
	_, err := o.Raw(sql, USER_WOLAI_TEAM).QueryRows(&teacherModels)
	if err != nil {
		seelog.Error(err.Error())
		return nil, err
	}
	for i := range teacherModels {
		teacher := teacherModels[i]
		teachers = append(teachers, POITeacher{
			POIUser: POIUser{
				UserId:      teacher.Id,
				Nickname:    teacher.Nickname,
				Avatar:      teacher.Avatar,
				AccessRight: teacher.AccessRight,
				Gender:      teacher.Gender},
			ServiceTime:      teacher.ServiceTime,
			School:           teacher.SchoolName,
			Department:       teacher.DeptName,
			PricePerHour:     teacher.PricePerHour,
			RealPricePerHour: teacher.RealPricePerHour})
	}
	return teachers, nil
}

func QueryTeacher(userId int64) *POITeacher {
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("users.id,users.nickname,users.avatar,users.access_right,users.gender,teacher_profile.service_time, teacher_profile.price_per_hour,teacher_profile.real_price_per_hour,school.name school_name,department.name dept_name").
		From("users").LeftJoin("teacher_profile").On("users.id = teacher_profile.user_id").
		LeftJoin("school").On("teacher_profile.school_id = school.id").
		LeftJoin("department").On("teacher_profile.department_id = department.id").
		Where("users.id = ?")
	sql := qb.String()
	o := orm.NewOrm()
	var teacherModel POITeacherModel
	err := o.Raw(sql, userId).QueryRow(&teacherModel)
	if err != nil {
		seelog.Error("userId:", userId, " ", err.Error())
		return nil
	}
	teacher := POITeacher{
		POIUser: POIUser{
			UserId:      teacherModel.Id,
			Nickname:    teacherModel.Nickname,
			Avatar:      teacherModel.Avatar,
			Gender:      teacherModel.Gender,
			AccessRight: teacherModel.AccessRight},
		ServiceTime:      teacherModel.ServiceTime,
		School:           teacherModel.SchoolName,
		Department:       teacherModel.DeptName,
		PricePerHour:     teacherModel.PricePerHour,
		RealPricePerHour: teacherModel.RealPricePerHour}
	return &teacher
}

func QueryTeacherProfile(userId int64) (*POITeacherProfile, error) {
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("users.id,users.nickname,users.avatar,users.access_right,users.gender,teacher_profile.service_time," +
		"teacher_profile.intro,teacher_profile.extra,teacher_profile.price_per_hour,teacher_profile.real_price_per_hour," +
		"school.name school_name, department.name dept_name,teacher_profile.rating").
		From("users").LeftJoin("teacher_profile").On("users.id = teacher_profile.user_id").
		LeftJoin("school").On("teacher_profile.school_id = school.id").
		LeftJoin("department").On("teacher_profile.department_id = department.id").
		Where("users.id = ?")
	sql := qb.String()
	o := orm.NewOrm()
	var teacherModel POITeacherModel
	err := o.Raw(sql, userId).QueryRow(&teacherModel)
	if err != nil {
		seelog.Error("userId:", userId, " ", err.Error())
		return nil, err
	}
	teacherProfile := POITeacherProfile{
		POITeacher: POITeacher{
			POIUser: POIUser{
				UserId:      teacherModel.Id,
				Nickname:    teacherModel.Nickname,
				AccessRight: teacherModel.AccessRight,
				Avatar:      teacherModel.Avatar,
				Gender:      teacherModel.Gender},
			ServiceTime:      teacherModel.ServiceTime,
			School:           teacherModel.SchoolName,
			Department:       teacherModel.DeptName,
			PricePerHour:     teacherModel.PricePerHour,
			RealPricePerHour: teacherModel.RealPricePerHour},
		Intro:  teacherModel.Intro,
		Extra:  teacherModel.Extra,
		Rating: teacherModel.Rating}
	teacherProfile.LabelList = QueryTeacherLabelByUserId(teacherModel.Id)
	teacherProfile.SubjectList = QueryTeacherSubjectByUserId(teacherModel.Id)
	teacherProfile.EducationList = QueryTeacherResumeByUserId(teacherModel.Id)
	return &teacherProfile, nil
}

func QueryTeacherProfileByUserId(userId int64) *POITeacherProfileModel {
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("user_id,school_id,department_id,service_time,rating,intro,extra,price_per_hour,real_price_per_hour").
		From("teacher_profile").Where("user_id = ?")
	sql := qb.String()
	o := orm.NewOrm()
	profile := POITeacherProfileModel{}
	err := o.Raw(sql, userId).QueryRow(&profile)
	if err != nil {
		seelog.Error("userId:", userId, " ", err.Error())
		return nil
	}
	return &profile
}

func QueryTeacherLabelByUserId(userId int64) []string {
	labels := make([]string, 0)
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("teacher_label.name").From("teacher_label").InnerJoin("teacher_to_label").
		On("teacher_to_label.label_id = teacher_label.id").Where("teacher_to_label.user_id = ?")
	o := orm.NewOrm()
	sql := qb.String()
	_, err := o.Raw(sql, userId).QueryRows(&labels)
	if err != nil {
		seelog.Error("userId:", userId, " ", err.Error())
		return nil
	}
	return labels
}

func QueryTeacherSubjectByUserId(userId int64) POITeacherSubjects {
	subjects := make(POITeacherSubjects, 0)
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("subject.name,teacher_to_subject.description").From("teacher_to_subject").
		InnerJoin("subject").On("teacher_to_subject.subject_id = subject.id").Where("teacher_to_subject.user_id = ?")
	sql := qb.String()
	o := orm.NewOrm()
	_, err := o.Raw(sql, userId).QueryRows(&subjects)
	if err != nil {
		seelog.Error("userId:", userId, " ", err.Error())
		return nil
	}
	return subjects
}

func QueryTeacherResumeByUserId(userId int64) POITeacherResumes {
	resumes := make(POITeacherResumes, 0)
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("id,user_id,start,stop,name").From("teacher_to_resume").Where("user_id = ?")
	sql := qb.String()
	o := orm.NewOrm()
	_, err := o.Raw(sql, userId).QueryRows(&resumes)
	if err != nil {
		seelog.Error("userId:", userId, " ", err.Error())
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
		seelog.Error("userId:", userId, " length:", length, " ", err.Error())
	}
}

func QueryTeacherLabelByName(name string) *POITeacherLabel {
	o := orm.NewOrm()
	teacherLabel := POITeacherLabel{}
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("id,name").From("teacher_label").Where("name = ?")
	sql := qb.String()
	err := o.Raw(sql, name).QueryRow(&teacherLabel)
	if err != nil {
		seelog.Error("name:", name, " ", err.Error())
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
		seelog.Error("name:", name, " ", err.Error())
		return 0
	}
	return id
}

func QueryTeacherToLabel(userId, labelId int64) *POITeacherToLabel {
	o := orm.NewOrm()
	ttl := POITeacherToLabel{}
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("id,user_id,label_id").From("teacher_to_lable").Where("user_id = ? and label_id = ?")
	sql := qb.String()
	err := o.Raw(sql, userId, labelId).QueryRow(&ttl)
	if err != nil {
		seelog.Error("userId:", userId, " labelId:", labelId, " ", err.Error())
		return nil
	}
	return &ttl
}

func InsertTeacherToLabel(teacherLabel *POITeacherToLabel) int64 {
	o := orm.NewOrm()
	id, err := o.Insert(teacherLabel)
	if err != nil {
		seelog.Error("teacherLabel:", teacherLabel, " ", err.Error())
		return 0
	}
	return id
}

func QueryTeacherToSubject(userId, subjectId int64) *POITeacherToSubject {
	o := orm.NewOrm()
	tts := POITeacherToSubject{}
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("id,user_id,subject_id,description").From("teacher_to_subject").Where("user_id = ? and subject_id = ?")
	sql := qb.String()
	err := o.Raw(sql, userId, subjectId).QueryRow(&tts)
	if err != nil {
		seelog.Error("userId:", userId, " subjectId:", subjectId, " ", err.Error())
		return nil
	}
	return &tts
}

func InsertTeacherToSubject(teacherSubject *POITeacherToSubject) int64 {
	o := orm.NewOrm()
	id, err := o.Insert(teacherSubject)
	if err != nil {
		seelog.Error("teacherSubject:", teacherSubject, " ", err.Error())
		return 0
	}
	return id
}

func InsertTeacherToResume(resume *POITeacherResume) int64 {
	o := orm.NewOrm()
	id, err := o.Insert(resume)
	if err != nil {
		seelog.Error("resume:", resume, " ", err.Error())
		return 0
	}
	return id
}

func InsertTeacherProfile(profile *POITeacherProfileModel) int64 {
	o := orm.NewOrm()
	id, err := o.Insert(profile)
	if err != nil {
		seelog.Error("profile:", profile, " ", err.Error())
		return 0
	}
	return id
}

//搜索老师
func QueryTeachersByCond(userId int64, keyword string, pageNum, pageCount int64) (POITeacherModels, error) {
	start := pageNum * pageCount
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("users.id,users.nickname,users.phone,users.access_right,users.avatar,users.gender,teacher_profile.service_time, " +
		"teacher_profile.price_per_hour,teacher_profile.real_price_per_hour,school.name school_name,department.name dept_name").
		From("users").InnerJoin("teacher_profile").On("users.id = teacher_profile.user_id").
		InnerJoin("school").On("teacher_profile.school_id = school.id").
		InnerJoin("department").On("teacher_profile.department_id = department.id").
		Where("users.access_right = 2 and users.status = 0 and (users.nickname like ? or users.phone like ?)").Limit(int(pageCount)).Offset(int(start))
	sql := qb.String()
	o := orm.NewOrm()
	var teacherModels POITeacherModels
	_, err := o.Raw(sql, "%"+keyword+"%", "%"+keyword+"%").QueryRows(&teacherModels)
	if err != nil {
		seelog.Error("keyword:", keyword, " ", err.Error())
		return teacherModels, err
	}
	return teacherModels, nil
}

//搜索老师和学生
func QueryUsersByCond(userId int64, keyword string, pageNum, pageCount int64) (POITeacherModels, error) {
	start := pageNum * pageCount
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("users.id,users.nickname,users.phone,users.avatar,users.gender,users.access_right,teacher_profile.service_time," +
		"teacher_profile.price_per_hour,teacher_profile.real_price_per_hour,school.name school_name,department.name dept_name").
		From("users").LeftJoin("teacher_profile").On("users.id = teacher_profile.user_id").
		LeftJoin("school").On("teacher_profile.school_id = school.id").
		LeftJoin("department").On("teacher_profile.department_id = department.id").
		Where("users.access_right in (2,3) and users.status = 0 and (users.nickname like ? or users.phone like ?)").Limit(int(pageCount)).Offset(int(start))
	sql := qb.String()
	o := orm.NewOrm()
	var teacherModels POITeacherModels
	_, err := o.Raw(sql, "%"+keyword+"%", "%"+keyword+"%").QueryRows(&teacherModels)
	if err != nil {
		seelog.Error("keyword:", keyword, " ", err.Error())
		return teacherModels, err
	}
	return teacherModels, nil
}
