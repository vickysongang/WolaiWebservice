// course_renew
package course

import (
	"WolaiWebservice/models"

	"github.com/astaxie/beego/orm"
)

func GetCourseRenewWaitingRecord(userId, courseId int64) *models.CourseRenewRecord {
	o := orm.NewOrm()
	var renewRecord models.CourseRenewRecord
	o.QueryTable(new(models.CourseRenewRecord).TableName()).
		Filter("course_id", courseId).
		Filter("user_id", userId).
		Filter("status", models.COURSE_RENEW_STATUS_WAITING).
		OrderBy("create_time").Limit(1).One(&renewRecord)
	if renewRecord.Id == 0 {
		return nil
	}
	return &renewRecord
}
