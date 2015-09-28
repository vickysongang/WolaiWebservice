// POICourse
package models

import (
	"POIWolaiWebService/utils"
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"
)

const (
	COURSE_UNOPEN  = "unopen"
	COURSE_OPENING = "opening"
	COURSE_SERVING = "serving"
	COURSE_EXPIRED = "expired"

	COURSE_JOIN  = "join"
	COURSE_RENEW = "renew"

	COURSE_PURCHASE_PENDING   = "pending"
	COURSE_PURCHASE_COMPLETED = "completed"
)

type POICourse struct {
	Id       int64  `json:"id" orm:"pk"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Intro    string `json:"intro"`
	Price    int64  `json:"price"`
	Banner   string `json:"banner"`
	Type     int64  `json:"-"`
	Length   int64  `json:"length"`
	TimeUnit string `json:"timeUnit"`
}

type POICourses []*POICourse

type POIUserToCourse struct {
	Id         int64     `json:"-" orm:"pk"`
	UserId     int64     `json:"userId"`
	CourseId   int64     `json:"courseId"`
	Status     string    `json:"status"`
	TimeFrom   time.Time `json:"timeFrom"`
	TimeTo     time.Time `json:"timeTo"`
	CreateTime time.Time `json:"-" orm:"auto_now_add;type(datetime)"`
}

type POIUserToCourses []*POIUserToCourse

type POICoursePurchaseRecord struct {
	Id         int64     `json:"-" orm:"pk"`
	UserId     int64     `json:"userId"`
	CourseId   int64     `json:"courseId"`
	CreateTime time.Time `json:"-" orm:"auto_now_add;type(datetime)"`
	TimeTo     time.Time `json:"timeTo"`
	Count      int64     `json:"count"`
	Type       string    `json:"type"`
	Status     string    `json:"status"`
}

type POICourse4User struct {
	Course       *POICourse       `json:"courseBaseInfo"`
	UserToCourse *POIUserToCourse `json:"courseStatusInfo"`
	CurrentTime  time.Time        `json:"currentTime"`
	JoinCount    int64            `json:"joinCount"`
}

type POICourseInfos []*POICourse4User

func (c *POICourse) TableName() string {
	return "courses"
}

func (c *POIUserToCourse) TableName() string {
	return "user_to_course"
}

func (c *POICoursePurchaseRecord) TableName() string {
	return "course_purchase_record"
}

func init() {
	orm.RegisterModel(new(POICourse), new(POIUserToCourse), new(POICoursePurchaseRecord))
}

func InsertCourse(course *POICourse) (*POICourse, error) {
	o := orm.NewOrm()
	id, err := o.Insert(course)
	if err != nil {
		return nil, err
	}
	course.Id = id
	return course, nil
}

func QueryCourses() (POICourses, error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("id,title,subtitle,intro,price,banner,length,time_unit").
		From("courses").Where("type = 1")
	sql := qb.String()
	courses := make(POICourses, 0)
	_, err := o.Raw(sql).QueryRows(&courses)
	if err != nil {
		return nil, err
	}
	return courses, nil
}

func QueryCourseById(courseId int64) (*POICourse, error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("id,title,subtitle,intro,price,banner,length,time_unit").
		From("courses").Where("id = ? and type = 1")
	sql := qb.String()
	course := POICourse{}
	err := o.Raw(sql, courseId).QueryRow(&course)
	if err != nil {
		return nil, err
	}
	return &course, nil
}

func QueryDefaultCourseId() (int64, error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("id").From("courses").Where("type = 1").Limit(1)
	sql := qb.String()
	var courseId int64
	err := o.Raw(sql).QueryRow(&courseId)
	if err != nil {
		return 0, err
	}
	return courseId, nil
}

func QueryCourse4User(userId int64) (*POIUserToCourse, error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	course4User := POIUserToCourse{}
	qb.Select("course_id,user_id,status,time_from,time_to").From("user_to_course").Where("user_id = ? and status = 'serving'").Limit(1)
	sql := qb.String()
	err := o.Raw(sql, userId).QueryRow(&course4User)
	if err != nil {
		return nil, err
	}
	return &course4User, nil
}

/*
 * 查询赠送课程
 */
func QueryGiveCourse() (*POICourse, error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("id,title,subtitle,intro,price,banner,length,time_unit").
		From("courses").Where("type = 0").Limit(1)
	sql := qb.String()
	course := POICourse{}
	err := o.Raw(sql).QueryRow(&course)
	if err != nil {
		return nil, err
	}
	return &course, nil
}

func InsertUserToCourse(userToCourse *POIUserToCourse) (*POIUserToCourse, error) {
	o := orm.NewOrm()
	id, err := o.Insert(userToCourse)
	if err != nil {
		return nil, err
	}
	userToCourse.Id = id
	return userToCourse, nil
}

func QueryUserToCourse(courseId, userId int64) (*POIUserToCourse, error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("id,user_id,course_id,status,time_from,time_to,create_time").
		From("user_to_course").Where("course_id = ? and user_id = ?")
	sql := qb.String()
	userToCourse := POIUserToCourse{}
	err := o.Raw(sql, courseId, userId).QueryRow(&userToCourse)
	if err != nil {
		return nil, err
	}
	return &userToCourse, nil
}

func UpdateUserCourseInfo(userId int64, courseId int64, updateInfo map[string]interface{}) {
	o := orm.NewOrm()
	var params orm.Params = make(orm.Params)
	for k, v := range updateInfo {
		params[k] = v
	}
	_, err := o.QueryTable("user_to_course").Filter("user_id", userId).Filter("course_id", courseId).Update(params)
	if err != nil {
		seelog.Error("userId:", userId, " updateInfo:", updateInfo, " ", err.Error())
	}
}

func UpdateUserCourseInfoById(userToCourseId int64, updateInfo map[string]interface{}) {
	o := orm.NewOrm()
	var params orm.Params = make(orm.Params)
	for k, v := range updateInfo {
		params[k] = v
	}
	_, err := o.QueryTable("user_to_course").Filter("id", userToCourseId).Update(params)
	if err != nil {
		seelog.Error("userToCourseId:", userToCourseId, " updateInfo:", updateInfo, " ", err.Error())
	}
}

func QueryExpiredCourses(processTime string) (POIUserToCourses, error) {
	fmt.Println("processTime:", processTime)
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("id,user_id,course_id,status,time_from,time_to").From("user_to_course").Where("status = 'serving' and time_to < ?")
	sql := qb.String()
	userToCourses := make(POIUserToCourses, 0)
	_, err := o.Raw(sql, processTime).QueryRows(&userToCourses)
	if err != nil {
		return nil, err
	}
	return userToCourses, nil
}

func InsertCoursePurchaseRecord(purchaseRecord *POICoursePurchaseRecord) (*POICoursePurchaseRecord, error) {
	o := orm.NewOrm()
	id, err := o.Insert(purchaseRecord)
	if err != nil {
		return nil, err
	}
	purchaseRecord.Id = id
	return purchaseRecord, nil
}

func QueryPendingPurchaseRecord(userId, courseId int64) (*POICoursePurchaseRecord, error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("id,user_id,course_id,time_to,create_time,count,type,status").
		From("course_purchase_record").Where("course_id = ? and user_id = ? and status = 'pending'")
	sql := qb.String()
	purchaseRecord := POICoursePurchaseRecord{}
	err := o.Raw(sql, courseId, userId).QueryRow(&purchaseRecord)
	if err != nil {
		return nil, err
	}
	return &purchaseRecord, nil
}

func UpdatePurchaseRecord(userId int64, courseId int64, updateInfo map[string]interface{}) {
	o := orm.NewOrm()
	var params orm.Params = make(orm.Params)
	for k, v := range updateInfo {
		params[k] = v
	}
	_, err := o.QueryTable("course_purchase_record").Filter("user_id", userId).Filter("course_id", courseId).
		Filter("status", COURSE_PURCHASE_PENDING).Update(params)
	if err != nil {
		seelog.Error("userId:", userId, " updateInfo:", updateInfo, " ", err.Error())
	}
}

func IsUserFree4Session(userId int64, currTime string) bool {
	o := orm.NewOrm()
	count, err := o.QueryTable("user_to_course").Filter("user_id", userId).Filter("status", "serving").Filter("time_to__gte", currTime).Count()
	if err != nil {
		return false
	}
	if count > 0 {
		return true
	}
	return false
}

func GetCourseJoinCount(courseId int64) int64 {
	o := orm.NewOrm()
	count, err := o.QueryTable("user_to_course").Filter("course_id", courseId).Count()
	if err != nil {
		return 0
	}
	return count
}
