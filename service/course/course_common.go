package course

import (
	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
)

func QueryCourseBanners() ([]*models.CourseBanners, error) {
	o := orm.NewOrm()

	var banners []*models.CourseBanners
	_, err := o.QueryTable("course_banners").
		OrderBy("rank").All(&banners)
	return banners, err
}

func GetCourseStudentCount(courseId int64) int64 {
	course, err := models.ReadCourse(courseId)
	if err != nil {
		return 0
	}
	o := orm.NewOrm()
	var studentCount int64
	if course.Type == models.COURSE_TYPE_DELUXE {
		studentCount, _ = o.QueryTable(new(models.CoursePurchaseRecord).TableName()).
			Filter("course_id", courseId).
			Count()
	} else if course.Type == models.COURSE_TYPE_AUDITION {
		studentCount, _ = o.QueryTable(new(models.CourseAuditionRecord).TableName()).
			Filter("course_id", courseId).
			Count()
	}
	return studentCount
}

func GetConversationCourseCount(studentId, teacherId int64) int64 {
	o := orm.NewOrm()
	count, _ := o.QueryTable("course_purchase_record").
		Filter("user_id", studentId).
		Filter("teacher_id", teacherId).Count()
	return count
}

func GetConversationAuditonCourseCount(studentId, teacherId int64) int64 {
	o := orm.NewOrm()
	count, _ := o.QueryTable("course_audition_record").
		Exclude("status", models.AUDITION_RECORD_STATUS_COMPLETE).
		Filter("user_id", studentId).
		Filter("teacher_id", teacherId).Count()
	return count
}

func GetSessionIdByChapter(chapterId int64) int64 {
	var sessionId int64
	o := orm.NewOrm()
	var order models.Order
	o.QueryTable(new(models.Order).TableName()).
		Filter("status", models.ORDER_STATUS_CONFIRMED).
		Filter("chapter_id", chapterId).
		OrderBy("-create_time").Limit(1).One(&order)
	if order.Id != 0 {
		var session models.Session
		o.QueryTable(new(models.Session).TableName()).
			Filter("order_id", order.Id).One(&session)
		if session.Id != 0 {
			sessionId = session.Id
		}
	}
	return sessionId
}

func QueryCourseContentIntros(courseId int64) ([]models.CourseContentIntro, error) {
	o := orm.NewOrm()
	intros := make([]models.CourseContentIntro, 0)
	_, err := o.QueryTable(new(models.CourseContentIntro).TableName()).
		Filter("course_id", courseId).OrderBy("rank").All(&intros)
	return intros, err
}

func QueryCourseTeachers(courseId int64) ([]*models.CourseToTeacher, error) {
	o := orm.NewOrm()
	var courseTeachers []*models.CourseToTeacher
	_, err := o.QueryTable("course_to_teachers").
		Filter("course_id", courseId).
		All(&courseTeachers)
	return courseTeachers, err
}

func QueryModuleCourses(moduleId int64) ([]*models.CourseToModule, error) {
	o := orm.NewOrm()
	var moduleCourses []*models.CourseToModule
	_, err := o.QueryTable("course_to_module").
		Filter("module_id", moduleId).
		Filter("recommend", 1).
		OrderBy("rank").All(&moduleCourses)
	return moduleCourses, err
}

func QueryCourseModules(moduleId, page, count int64) ([]*models.CourseToModule, error) {
	o := orm.NewOrm()
	var courseModules []*models.CourseToModule
	_, err := o.QueryTable("course_to_module").
		Filter("module_id", moduleId).
		OrderBy("rank").
		Offset(page * count).Limit(count).
		All(&courseModules)
	return courseModules, err
}
