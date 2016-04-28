// course_relation
package course

import (
	"WolaiWebservice/models"

	"github.com/astaxie/beego/orm"
)

func GetCourseRelation(recordId int64, recordType string) (*models.CourseRelation, error) {
	var courseRelation models.CourseRelation
	o := orm.NewOrm()
	err := o.QueryTable("course_relation").
		Filter("record_id", recordId).
		Filter("type", recordType).
		One(&courseRelation)
	return &courseRelation, err
}
