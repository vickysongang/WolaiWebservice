package course

import (
	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
	"WolaiWebservice/websocket"
)

func HandleDeluxeCourseActionQuickbuy(userId int64, courseId int64) (int64, *actionProceedResponse) {
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
		var auditionRecord models.CourseAuditionRecord
		o.QueryTable(new(models.CourseAuditionRecord).TableName()).Filter("source_course_id", courseId).Filter("user_id", userId).
			One(&auditionRecord)
		var teacherId, priceHourly, salaryHourly, priceTotal int64
		var status string
		if auditionRecord.Id != 0 && auditionRecord.TeacherId != 0 {
			teacherId = auditionRecord.TeacherId
			priceHourly = auditionRecord.PriceHourly
			salaryHourly = auditionRecord.SalaryHourly
			priceTotal = priceHourly * currentRecord.ChapterCount
			status = models.PURCHASE_RECORD_STATUS_WAITING
			migerateCourseChapter(userId, auditionRecord.TeacherId, courseId)
		} else {
			status = models.PURCHASE_RECORD_STATUS_APPLY
		}

		chapterCount := courseService.GetCourseChapterCount(courseId)
		// 如果用户没有购买过，创建购买记录
		newRecord := models.CoursePurchaseRecord{
			CourseId:       courseId,
			UserId:         userId,
			TeacherId:      teacherId,
			PriceHourly:    priceHourly,
			SalaryHourly:   salaryHourly,
			PriceTotal:     priceTotal,
			AuditionStatus: models.PURCHASE_RECORD_STATUS_IDLE,
			PurchaseStatus: status,
			TraceStatus:    models.PURCHASE_RECORD_TRACE_STATUS_IDLE,
			ChapterCount:   chapterCount,
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

	case record.PurchaseStatus == models.PURCHASE_RECORD_STATUS_PAID:
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
			err = createDeluxeCourseOrder(record.Id, true)

			if err != nil {
				response = actionProceedResponse{
					Action:  ACTION_PROCEED_NULL,
					Message: "导师还没准备好下一节课时，请耐心等待哦",
					Extra:   nullObject{},
				}
			}
		}

	case record.PurchaseStatus == models.PURCHASE_RECORD_STATUS_COMPLETE:
		response = actionProceedResponse{
			Action:  ACTION_PROCEED_RENEW,
			Message: "课时已经全部上完啦！可以根据导师的课程计划续课喔",
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

func migerateCourseChapter(userId, teacherId, courseId int64) {
	courseChapters, err := courseService.QueryCourseChapters(courseId)
	if err != nil {
		return
	}
	var relId int64
	oldCourseRelation, _ := courseService.QueryCourseRelation(courseId, userId, teacherId)
	if oldCourseRelation.Id == 0 {
		courseRelation := models.CourseRelation{
			CourseId:  courseId,
			UserId:    userId,
			TeacherId: teacherId,
		}
		relId, err = models.InsertCourseRelation(&courseRelation)
	} else {
		relId = oldCourseRelation.Id
	}

	if err != nil {
		return
	}

	for _, courseChapter := range courseChapters {
		oldCustomChapter, _ := courseService.QueryCourseCustomChapter(courseId, courseChapter.Period, userId)
		if oldCustomChapter.Id == 0 {
			customChapter := models.CourseCustomChapter{
				CourseId:  courseChapter.CourseId,
				Title:     courseChapter.Title,
				Abstract:  courseChapter.Abstract,
				Period:    courseChapter.Period,
				UserId:    userId,
				TeacherId: teacherId,
				AttachId:  courseChapter.AttachId,
				RelId:     relId,
			}
			models.InsertCourseCustomChapter(&customChapter)
		} else {
			customChapterInfo := map[string]interface{}{
				"Title":     courseChapter.Title,
				"Abstract":  courseChapter.Abstract,
				"Period":    courseChapter.Period,
				"TeacherId": teacherId,
				"AttachId":  courseChapter.AttachId,
				"RelId":     relId,
			}
			models.UpdateCourseCustomChapter(oldCustomChapter.Id, customChapterInfo)
		}
	}
}
