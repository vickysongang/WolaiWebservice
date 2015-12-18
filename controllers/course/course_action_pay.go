package course

import (
	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
)

func HandleCourseActionPay(userId int64, courseId int64, payType string) int64 {
	o := orm.NewOrm()

	_, err := models.ReadCourse(courseId)
	if err != nil {
		return 2
	}

	// 先查询该用户是否有购买（或试图购买）过这个课程
	var currentRecord models.CoursePurchaseRecord
	var record *models.CoursePurchaseRecord
	err = o.QueryTable("course_purchase_record").Filter("course_id", courseId).Filter("user_id", userId).
		One(&currentRecord)
	if err != nil {
		return 2
	}

	record = &currentRecord

	switch payType {
	case PAYMENT_TYPE_AUDITION:
		if record.AuditionStatus != models.PURCHASE_RECORD_STATUS_WAITING {
			return 2
		}

		recordInfo := map[string]interface{}{
			"audition_status": models.PURCHASE_RECORD_STATUS_PAID,
			//"last_update_time": "NOW()",
		}

		record, err = models.UpdateCoursePurchaseRecord(record.Id, recordInfo)
		if err != nil {
			return 2
		}

	case PAYMENT_TYPE_PURCHASE:
		if record.PurchaseStatus != models.PURCHASE_RECORD_STATUS_WAITING {
			return 2
		}

		recordInfo := map[string]interface{}{
			"purchase_status": models.PURCHASE_RECORD_STATUS_PAID,
			//"last_update_time": "NOW()",
		}

		record, err = models.UpdateCoursePurchaseRecord(record.Id, recordInfo)
		if err != nil {
			return 2
		}
	}

	return 0
}
