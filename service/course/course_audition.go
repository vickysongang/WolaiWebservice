// course_audition
package course

import (
	"WolaiWebservice/models"

	"github.com/astaxie/beego/orm"
)

func GetAuditionCourseStudentCount(courseId int64) int64 {
	o := orm.NewOrm()
	auditionCount, _ := o.QueryTable(new(models.CourseAuditionRecord).TableName()).Count()
	return auditionCount
}

func QueryAuditionCourse() *models.Course {
	o := orm.NewOrm()
	var course models.Course
	o.QueryTable(new(models.Course).TableName()).Filter("type", models.COURSE_TYPE_AUDITION).One(&course)
	if course.Id == 0 {
		return nil
	}
	return &course
}

func GetUncompletedAuditionRecord(userId int64) *models.CourseAuditionRecord {
	o := orm.NewOrm()
	var record models.CourseAuditionRecord
	o.QueryTable(new(models.CourseAuditionRecord).TableName()).
		Filter("user_id", userId).Exclude("status", models.AUDITION_RECORD_STATUS_COMPLETE).One(&record)
	if record.Id == 0 {
		return nil
	}
	return &record
}
