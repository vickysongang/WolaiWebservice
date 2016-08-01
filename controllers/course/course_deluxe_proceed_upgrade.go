package course

import (
	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
	"WolaiWebservice/websocket"
)

func HandleDeluxeCourseActionQuickbuy(userId int64, courseId int64) (int64, *actionProceedResponse) {
	course, err := models.ReadCourse(courseId)
	if err != nil {
		return 2, nil
	}

	// 先查询该用户是否有购买（或试图购买）过这个课程
	var record *models.CoursePurchaseRecord
	currentRecord, err := courseService.GetCoursePurchaseRecordByUserId(courseId, userId)
	if err == orm.ErrNoRows {

		chapterCount := courseService.GetCourseChapterCount(courseId)
		// 如果用户没有购买过，创建购买记录
		newRecord := models.CoursePurchaseRecord{
			CourseId:       courseId,
			UserId:         userId,
			AuditionStatus: models.PURCHASE_RECORD_STATUS_IDLE,
			PurchaseStatus: models.PURCHASE_RECORD_STATUS_APPLY,
			TraceStatus:    models.PURCHASE_RECORD_TRACE_STATUS_IDLE,
			ChapterCount:   chapterCount,
			PurchaseCount:  chapterCount,
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
	case record.PurchaseStatus == models.PURCHASE_RECORD_STATUS_IDLE:

		recordInfo := map[string]interface{}{
			"PurchaseStatus": models.PURCHASE_RECORD_STATUS_APPLY,
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

	case record.PurchaseStatus == models.PURCHASE_RECORD_STATUS_APPLY:

		// 学生已经申请过购买，但是客服还没有为其指派老师
		response = actionProceedResponse{
			Action:  ACTION_PROCEED_NULL,
			Message: "别着急...助教正在定制你的课程并为你匹配合适的导师哦",
			Extra:   nullObject{},
		}

	case record.PurchaseStatus == models.PURCHASE_RECORD_STATUS_WAITING:

		// 客服已经给学生匹配导师，直接支付购买费用
		payment := paymentInfo{
			Title:        PAYMENT_TITLE_PREFIX_PURCHASE + course.Name,
			Price:        record.PriceTotal,
			Comment:      PAYMENT_COMMENT_PURCHASE,
			Type:         PAYMENT_TYPE_PURCHASE,
			ChapterCount: record.ChapterCount,
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
