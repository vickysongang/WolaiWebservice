package course

import (
	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
	"WolaiWebservice/websocket"
)

func HandleAuditionCourseProceed(userId, recordId, courseId, sourceCourseId int64) (int64, *actionProceedResponse) {
	var err error
	var auditionRecordStatus string
	currentRecord, err := models.ReadCourseAuditionRecord(recordId)
	if err != nil {
		// 如果用户没有购买过，创建试听课购买记录
		newRecord := models.CourseAuditionRecord{
			CourseId:       courseId,
			UserId:         userId,
			Status:         models.AUDITION_RECORD_STATUS_APPLY,
			SourceCourseId: sourceCourseId,
			TraceStatus:    models.AUDITION_RECORD_TRACE_STATUS_IDLE,
		}
		_, err = models.CreateCourseAuditionRecord(&newRecord)
		if err != nil {
			return 2, nil
		}
		auditionRecordStatus = newRecord.Status
	} else {
		auditionRecordStatus = currentRecord.Status
	}
	course, err := models.ReadCourse(courseId)
	if err != nil {
		return 2, nil
	}
	var response actionProceedResponse
	switch auditionRecordStatus {
	case models.AUDITION_RECORD_STATUS_APPLY,
		models.AUDITION_RECORD_STATUS_WAITING:

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

	case models.AUDITION_RECORD_STATUS_PAID:
		if currentRecord.TeacherId == 0 {
			response = actionProceedResponse{
				Action:  ACTION_PROCEED_NULL,
				Message: "别着急...助教正在定制你的课程并为你匹配合适的导师哦",
				Extra:   nullObject{},
			}
			return 0, &response
		}
		// 学生已经支付了试听押金，开始上课！
		session := sessionInfo{
			TeacherId: currentRecord.TeacherId,
		}

		if websocket.OrderManager.HasOrderOnline(userId, currentRecord.TeacherId) {
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
			createAuditionCourseOrder(currentRecord.Id)
		}

	case models.AUDITION_RECORD_STATUS_COMPLETE:
		response = actionProceedResponse{
			Action:  ACTION_PROCEED_NULL,
			Message: "该试听课已经上完啦！",
			Extra:   nullObject{},
		}
	}
	return 0, &response
}

func HandleCourseActionAuditionCheck(userId int64) (int64, *actionProceedResponse) {
	var response actionProceedResponse
	auditionRecord := courseService.GetUncompletedAuditionRecord(userId)
	if auditionRecord != nil {
		auditionInfo := map[string]interface{}{
			"auditionCourseId": auditionRecord.CourseId,
			"sourceCourseId":   auditionRecord.SourceCourseId,
			"recordId":         auditionRecord.Id,
			"exist":            true,
		}
		response = actionProceedResponse{
			Action:  ACTION_PROCEED_NULL,
			Message: "你已经申请了一节试听课，建议先上完课哦！你也可以联系助教修改试听内容",
			Extra:   auditionInfo,
		}
	} else {
		auditionInfo := map[string]interface{}{
			"exist": false,
		}
		response = actionProceedResponse{
			Action:  ACTION_PROCEED_NULL,
			Message: "",
			Extra:   auditionInfo,
		}
	}
	return 0, &response
}
