// course_renew_detail
package course

import (
	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
	"errors"

	"github.com/astaxie/beego/orm"
)

const (
	COURSE_RENEW_TYPE_MANUAL = "manual"
	COURSE_RENEW_TYPE_AUTO   = "auto"
)

type CourseRenewDetail struct {
	CourseId    int64  `json:"courseId"`
	UserId      int64  `json:"userId"`
	TeacherId   int64  `json:"teacherId"`
	PriceHourly int64  `json:"priceHourly"`
	PriceTotal  int64  `json:"priceTotal"`
	RenewCount  int64  `json:"renewCount"`
	Type        string `json:"type"`
}

func GetCourseRenewDetail(courseId, userId int64) (int64, *CourseRenewDetail, error) {
	var detail CourseRenewDetail
	detail.CourseId = courseId
	detail.UserId = userId

	record := courseService.GetCourseRenewWaitingRecord(userId, courseId)
	if record != nil {
		detail.TeacherId = record.TeacherId
		detail.PriceHourly = record.PriceHourly
		detail.PriceTotal = record.PriceTotal
		detail.RenewCount = record.RenewCount
		detail.Type = COURSE_RENEW_TYPE_AUTO
	} else {
		o := orm.NewOrm()
		var purchaseRecord models.CoursePurchaseRecord
		o.QueryTable("course_purchase_record").Filter("course_id", courseId).Filter("user_id", userId).
			One(&purchaseRecord)
		if purchaseRecord.Id != 0 {
			detail.TeacherId = purchaseRecord.TeacherId
			detail.PriceHourly = purchaseRecord.PriceHourly
			detail.Type = COURSE_RENEW_TYPE_MANUAL
		} else {
			return 2, nil, errors.New("购买记录异常")
		}
	}
	return 0, &detail, nil
}
