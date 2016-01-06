package course

import (
	"errors"

	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
	"WolaiWebservice/service/trade"
)

func HandleCourseActionPay(userId int64, courseId int64, payType string) (int64, error) {
	var err error
	o := orm.NewOrm()

	_, err = models.ReadCourse(courseId)
	if err != nil {
		return 2, errors.New("课程信息异常")
	}

	// 先查询该用户是否有购买（或试图购买）过这个课程
	var currentRecord models.CoursePurchaseRecord
	var record *models.CoursePurchaseRecord
	err = o.QueryTable("course_purchase_record").Filter("course_id", courseId).Filter("user_id", userId).
		One(&currentRecord)
	if err != nil {
		return 2, errors.New("购买记录异常")
	}

	record = &currentRecord
	switch payType {
	case PAYMENT_TYPE_AUDITION:
		if record.AuditionStatus != models.PURCHASE_RECORD_STATUS_WAITING {
			return 2, errors.New("购买记录异常")
		}

		err = trade.HandleCourseAudition(record.Id, PAYMENT_PRICE_AUDITION)
		if err != nil {
			return 2, err
		}

		recordInfo := map[string]interface{}{
			"audition_status": models.PURCHASE_RECORD_STATUS_PAID,
		}

		record, err = models.UpdateCoursePurchaseRecord(record.Id, recordInfo)
		if err != nil {
			return 2, errors.New("购买记录异常")
		}

	case PAYMENT_TYPE_PURCHASE:
		if record.PurchaseStatus != models.PURCHASE_RECORD_STATUS_WAITING {
			return 2, errors.New("购买记录异常")
		}

		err = trade.HandleCoursePurchase(record.Id)
		if err != nil {
			return 2, err
		}

		recordInfo := map[string]interface{}{
			"purchase_status": models.PURCHASE_RECORD_STATUS_PAID,
		}

		if record.AuditionStatus == models.PURCHASE_RECORD_STATUS_IDLE ||
			record.AuditionStatus == models.PURCHASE_RECORD_STATUS_APPLY ||
			record.AuditionStatus == models.PURCHASE_RECORD_STATUS_WAITING {
			recordInfo["audition_status"] = models.PURCHASE_RECORD_STATUS_PAID
		}

		record, err = models.UpdateCoursePurchaseRecord(record.Id, recordInfo)
		if err != nil {
			return 2, errors.New("购买记录异常")
		}
	}

	return 0, nil
}
