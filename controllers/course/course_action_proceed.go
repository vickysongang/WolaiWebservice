package course

import (
	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
	"WolaiWebservice/websocket"
)

func HandleCourseActionProceed(userId int64, courseId int64) (int64, *actionProceedResponse) {
	var err error

	course, err := models.ReadCourse(courseId)
	if err != nil {
		return 2, nil
	}

	// 先查询该用户是否有购买（或试图购买）过这个课程
	var record *models.CoursePurchaseRecord
	currentRecord, err := courseService.GetCoursePurchaseRecordByUserId(courseId, userId)
	if err == orm.ErrNoRows {
		chaperCount := courseService.GetCourseChapterCount(courseId)
		// 如果用户没有购买过，创建购买记录
		newRecord := models.CoursePurchaseRecord{
			CourseId:       courseId,
			UserId:         userId,
			AuditionStatus: models.PURCHASE_RECORD_STATUS_APPLY,
			PurchaseStatus: models.PURCHASE_RECORD_STATUS_IDLE,
			TraceStatus:    models.PURCHASE_RECORD_TRACE_STATUS_IDLE,
			ChapterCount:   chaperCount,
			PurchaseCount:  chaperCount,
		}

		_, err = models.CreateCoursePurchaseRecord(&newRecord)
		if err != nil {
			return 2, nil
		}

		response := actionProceedResponse{
			Action:  ACTION_PROCEED_REFRESH,
			Message: "助教会在30分钟内与你取得联系，请保持电话畅通哦",
			Extra:   nullObject{},
		}

		return 0, &response
	} else if err != nil {

		// 如果到了这里说明数据库报错了...
		return 2, nil
	}
	record = &currentRecord

	var response actionProceedResponse
	//好了，我们拿到了用户的购买记录,现在玩一个游戏叫排列组合...
	switch {
	case record.AuditionStatus == models.PURCHASE_RECORD_STATUS_IDLE &&
		record.PurchaseStatus == models.PURCHASE_RECORD_STATUS_IDLE,
		record.AuditionStatus == models.PURCHASE_RECORD_STATUS_IDLE &&
			record.PurchaseStatus == models.PURCHASE_RECORD_STATUS_APPLY:

		// 学生在还没有被指派导师的时候申请试听
		recordInfo := map[string]interface{}{
			"AuditionStatus": models.PURCHASE_RECORD_STATUS_APPLY,
		}

		record, err = models.UpdateCoursePurchaseRecord(record.Id, recordInfo)
		if err != nil {
			return 2, nil
		}

		response = actionProceedResponse{
			Action:  ACTION_PROCEED_REFRESH,
			Message: "助教会在30分钟内与你取得联系，请保持电话畅通哦",
			Extra:   nullObject{},
		}

	case record.AuditionStatus == models.PURCHASE_RECORD_STATUS_APPLY &&
		record.PurchaseStatus == models.PURCHASE_RECORD_STATUS_IDLE,
		record.AuditionStatus == models.PURCHASE_RECORD_STATUS_APPLY &&
			record.PurchaseStatus == models.PURCHASE_RECORD_STATUS_APPLY:

		// 学生已经申请过试听，但是客服还没有为其指派老师
		response = actionProceedResponse{
			Action:  ACTION_PROCEED_NULL,
			Message: "稍等一下，助教正在为你匹配最优秀的导师哦",
			Extra:   nullObject{},
		}

	case record.AuditionStatus == models.PURCHASE_RECORD_STATUS_WAITING &&
		record.PurchaseStatus == models.PURCHASE_RECORD_STATUS_IDLE,
		record.AuditionStatus == models.PURCHASE_RECORD_STATUS_WAITING &&
			record.PurchaseStatus == models.PURCHASE_RECORD_STATUS_WAITING:

		// 客服已经为学生指派导师，学生支付试听押金
		payment := paymentInfo{
			Title:   PAYMENT_TITLE_PREFIX_AUDITION + course.Name,
			Price:   PAYMENT_PRICE_AUDITION,
			Comment: PAYMENT_COMMENT_AUDITION,
			Type:    PAYMENT_TYPE_AUDITION,
		}

		response = actionProceedResponse{
			Action:  ACTION_PROCEED_PAY,
			Message: "",
			Extra:   payment,
		}

	case record.AuditionStatus == models.PURCHASE_RECORD_STATUS_PAID &&
		record.PurchaseStatus == models.PURCHASE_RECORD_STATUS_IDLE,
		record.AuditionStatus == models.PURCHASE_RECORD_STATUS_PAID &&
			record.PurchaseStatus == models.PURCHASE_RECORD_STATUS_WAITING:

		// 学生已经支付了试听押金，开始上课！
		session := sessionInfo{
			TeacherId: record.TeacherId,
		}

		if websocket.OrderManager.HasOrderOnline(userId, record.TeacherId) {
			response = actionProceedResponse{
				Action:  ACTION_PROCEED_SERVE,
				Message: "你已经向该导师发过一条上课请求了，请耐心等待回复哦",
				Extra:   nullObject{},
			}
		} else {
			response = actionProceedResponse{
				Action:  ACTION_PROCEED_SERVE,
				Message: "",
				Extra:   session,
			}
			createDeluxeCourseOrder(record.Id, false)
		}

	case record.AuditionStatus == models.PURCHASE_RECORD_STATUS_COMPLETE &&
		record.PurchaseStatus == models.PURCHASE_RECORD_STATUS_IDLE,
		record.AuditionStatus == models.PURCHASE_RECORD_STATUS_COMPLETE &&
			record.PurchaseStatus == models.PURCHASE_RECORD_STATUS_WAITING:

		// 学生已经完成试听课程，学生须支付课程费用

		recordInfo := map[string]interface{}{
			"PurchaseStatus": models.PURCHASE_RECORD_STATUS_WAITING,
		}

		record, err = models.UpdateCoursePurchaseRecord(record.Id, recordInfo)
		if err != nil {
			return 2, nil
		}

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

	case record.AuditionStatus == models.PURCHASE_RECORD_STATUS_IDLE &&
		record.PurchaseStatus == models.PURCHASE_RECORD_STATUS_WAITING:

		// 学生在还没有被指派导师的时候申请试听
		recordInfo := map[string]interface{}{
			"AuditionStatus": models.PURCHASE_RECORD_STATUS_WAITING,
		}

		record, err = models.UpdateCoursePurchaseRecord(record.Id, recordInfo)
		if err != nil {
			return 2, nil
		}

		payment := paymentInfo{
			Title:   PAYMENT_TITLE_PREFIX_AUDITION + course.Name,
			Price:   PAYMENT_PRICE_AUDITION,
			Comment: PAYMENT_COMMENT_AUDITION,
			Type:    PAYMENT_TYPE_AUDITION,
		}

		response = actionProceedResponse{
			Action:  ACTION_PROCEED_PAY,
			Message: "",
			Extra:   payment,
		}

	case record.AuditionStatus == models.PURCHASE_RECORD_STATUS_PAID &&
		record.PurchaseStatus == models.PURCHASE_RECORD_STATUS_PAID,
		record.AuditionStatus == models.PURCHASE_RECORD_STATUS_COMPLETE &&
			record.PurchaseStatus == models.PURCHASE_RECORD_STATUS_PAID:

		// 学生已经完成试听，并且支付课程包费用，开始上课
		session := sessionInfo{
			TeacherId: record.TeacherId,
		}

		if websocket.OrderManager.HasOrderOnline(userId, record.TeacherId) {
			response = actionProceedResponse{
				Action:  ACTION_PROCEED_SERVE,
				Message: "你已经向该导师发过一条上课请求了，请耐心等待回复哦",
				Extra:   nullObject{},
			}

		} else {
			response = actionProceedResponse{
				Action:  ACTION_PROCEED_SERVE,
				Message: "",
				Extra:   session,
			}

			createDeluxeCourseOrder(record.Id, false)
		}

	case record.PurchaseStatus == models.PURCHASE_RECORD_STATUS_COMPLETE:

		// 学生的课程已经完成，无法继续操作
		response = actionProceedResponse{
			Action:  ACTION_PROCEED_NULL,
			Message: "您的课程已经完成，欢迎继续选购其他课程",
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
