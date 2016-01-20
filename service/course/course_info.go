package course

import (
	"github.com/astaxie/beego/orm"
	//"github.com/cihub/seelog"

	"WolaiWebservice/models"
)

func GetCourseStudentCount(courseId int64) int64 {
	var err error

	o := orm.NewOrm()

	studentCount, err := o.QueryTable(new(models.CoursePurchaseRecord).TableName()).
		Filter("course_id", courseId).
		Count()
	if err != nil {
		return 0
	}

	return studentCount
}

func GetCourseChapterCount(courseId int64) int64 {
	var err error

	o := orm.NewOrm()

	chapterCount, err := o.QueryTable(new(models.CourseChapter).TableName()).
		Filter("course_id", courseId).
		Count()
	if err != nil {
		return 0
	}

	return chapterCount
}
