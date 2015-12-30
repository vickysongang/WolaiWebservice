package course

import (
	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
)

func HandleCourseActionQuickbuy(userId int64, courseId int64) (int64, *actionProceedResponse) {
	o := orm.NewOrm()

	course, err := models.ReadCourse(courseId)
	if err != nil {
		return 2, nil
	}

	// 先查询该用户是否有购买（或试图购买）过这个课程
	var currentRecord models.CoursePurchaseRecord
	var record *models.CoursePurchaseRecord
	err = o.QueryTable("course_purchase_record").Filter("course_id", courseId).Filter("user_id", userId).
		One(&currentRecord)

	if err == orm.ErrNoRows {

		// 如果用户没有购买过，创建购买记录
		newRecord := models.CoursePurchaseRecord{
			CourseId:       courseId,
			UserId:         userId,
			AuditionStatus: models.PURCHASE_RECORD_STATUS_IDLE,
			PurchaseStatus: models.PURCHASE_RECORD_STATUS_APPLY,
		}

		_, err = models.CreateCoursePurchaseRecord(&newRecord)
		if err != nil {
			return 2, nil
		}

		response := actionProceedResponse{
			Action:  ACTION_PROCEED_REFRESH,
			Message: "申请成功，助教会在30分钟内与你取得联系，请保持电话畅通哦",
			Extra:   nullObject{},
		}

		return 0, &response
	} else if err != nil {

		// 如果到了这里说明数据库报错了...
		return 2, nil
	}
	record = &currentRecord

	var response actionProceedResponse

	switch {
	case record.PurchaseStatus == models.PURCHASE_RECORD_STATUS_IDLE &&
		record.AuditionStatus == models.PURCHASE_RECORD_STATUS_IDLE,
		record.PurchaseStatus == models.PURCHASE_RECORD_STATUS_IDLE &&
			record.AuditionStatus == models.PURCHASE_RECORD_STATUS_APPLY:

		// 学生在还没有被指派导师的时候申请试听
		recordInfo := map[string]interface{}{
			"purchase_status": models.PURCHASE_RECORD_STATUS_APPLY,
			//"last_update_time": "NOW()",
		}

		record, err = models.UpdateCoursePurchaseRecord(record.Id, recordInfo)
		if err != nil {
			return 2, nil
		}

		response = actionProceedResponse{
			Action:  ACTION_PROCEED_REFRESH,
			Message: "购买申请提交成功，助教会马上联系你哦",
			Extra:   nullObject{},
		}

	case record.PurchaseStatus == models.PURCHASE_RECORD_STATUS_IDLE &&
		record.AuditionStatus == models.PURCHASE_RECORD_STATUS_WAITING,
		record.PurchaseStatus == models.PURCHASE_RECORD_STATUS_IDLE &&
			record.AuditionStatus == models.PURCHASE_RECORD_STATUS_PAID,
		record.PurchaseStatus == models.PURCHASE_RECORD_STATUS_IDLE &&
			record.AuditionStatus == models.PURCHASE_RECORD_STATUS_COMPLETE:

		recordInfo := map[string]interface{}{
			"purchase_status": models.PURCHASE_RECORD_STATUS_WAITING,
			//"last_update_time": "NOW()",
		}

		record, err = models.UpdateCoursePurchaseRecord(record.Id, recordInfo)
		if err != nil {
			return 2, nil
		}

		// 客服已经给学生匹配导师，直接支付购买费用
		payment := paymentInfo{
			Title:   PAYMENT_TITLE_PREFIX_PURCHASE + course.Name,
			Price:   record.PriceTotal,
			Comment: PAYMENT_COMMENT_PURCHASE,
			Type:    PAYMENT_TYPE_PURCHASE,
		}

		response = actionProceedResponse{
			Action:  ACTION_PROCEED_PAY,
			Message: "",
			Extra:   payment,
		}

	case record.PurchaseStatus == models.PURCHASE_RECORD_STATUS_APPLY:

		// 学生已经申请过购买，但是客服还没有为其指派老师
		response = actionProceedResponse{
			Action:  ACTION_PROCEED_NULL,
			Message: "稍等一下，助教正在为你匹配最优秀的导师哦",
			Extra:   nullObject{},
		}

	case record.PurchaseStatus == models.PURCHASE_RECORD_STATUS_WAITING:

		// 客服已经给学生匹配导师，直接支付购买费用
		payment := paymentInfo{
			Title:   PAYMENT_TITLE_PREFIX_PURCHASE + course.Name,
			Price:   record.PriceTotal,
			Comment: PAYMENT_COMMENT_PURCHASE,
			Type:    PAYMENT_TYPE_PURCHASE,
		}

		response = actionProceedResponse{
			Action:  ACTION_PROCEED_PAY,
			Message: "",
			Extra:   payment,
		}

	case record.PurchaseStatus == models.PURCHASE_RECORD_STATUS_PAID,
		record.PurchaseStatus == models.PURCHASE_RECORD_STATUS_COMPLETE:

		// 付过款以后这里的按钮就不可以点了，如果APP没处理好让他点了的话也什么都不会发生
		response = actionProceedResponse{
			Action:  ACTION_PROCEED_NULL,
			Message: "",
			Extra:   nullObject{},
		}

	default:
		response = actionProceedResponse{
			Action:  ACTION_PROCEED_NULL,
			Message: "购买操作异常，请联系助教",
			Extra:   nullObject{},
		}
	}

	return 0, &response
}
