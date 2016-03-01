// course_relation
package course

import (
	"WolaiWebservice/models"

	"github.com/astaxie/beego/orm"
)

func GetCourseRelation(courseId, studentId, teacherId int64) (*models.CourseRelation, error) {
	var courseRelation models.CourseRelation
	o := orm.NewOrm()
	err := o.QueryTable("course_relation").
		Filter("course_id", courseId).
		Filter("teacher_id", teacherId).
		Filter("user_id", studentId).One(&courseRelation)
	return &courseRelation, err
}
