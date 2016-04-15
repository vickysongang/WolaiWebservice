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
