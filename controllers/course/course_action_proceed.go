package course

import (
	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
	"WolaiWebservice/websocket"
)

type nullObject struct{}

type paymentInfo struct {
	Title   string `json:"title"`
	Price   int64  `json:"price"`
	Comment string `json:"comment"`
	Type    string `json:"type"`
}

type sessionInfo struct {
	TeacherId int64 `json:"teacherId"`
}

type actionProceedResponse struct {
	Action  string      `json:"action"`
	Message string      `json:"message"`
	Extra   interface{} `json:"extra,omitempty"`
}

const (
	ACTION_PROCEED_NULL    = "null"
	ACTION_PROCEED_REFRESH = "refresh"
	ACTION_PROCEED_PAY     = "pay"
	ACTION_PROCEED_SERVE   = "serve"

	PAYMENT_TITLE_PREFIX_AUDITION = "课程试听-"
	PAYMENT_TITLE_PREFIX_PURCHASE = "课程购买-"

	PAYMENT_TYPE_AUDITION = "audition"
	PAYMENT_TYPE_PURCHASE = "purchase"

	PAYMENT_COMMENT_AUDITION = "试听支付"
	PAYMENT_COMMENT_PURCHASE = "无"

	PAYMENT_PRICE_AUDITION = 100
)

func HandleCourseActionProceed(userId int64, courseId int64, sourceCourseId int64) (int64, *actionProceedResponse) {
	var course *models.Course
	var err error
	if courseId == 0 { //代表试听课不是从课程进入的
		course = courseService.QueryAuditionCourse()
		if course == nil {
			return 2, nil
		}
	} else {
		course, err = models.ReadCourse(courseId)
		if err != nil {
			return 2, nil
		}
	}
	if course.Type == models.COURSE_TYPE_DELUXE {
		status, response := HandleDeluxeCourseActionProceed(userId, course)
		return status, response
	} else if course.Type == models.COURSE_TYPE_AUDITION {
		status, response := HandleAuditionCourseActionProceed(userId, course, sourceCourseId)
		return status, response
	}
	return 0, nil
}

//处理旧版本的试听
func HandleDeluxeCourseActionProceed(userId int64, course *models.Course) (int64, *actionProceedResponse) {
	var err error
	courseId := course.Id
	o := orm.NewOrm()

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
			AuditionStatus: models.PURCHASE_RECORD_STATUS_APPLY,
			PurchaseStatus: models.PURCHASE_RECORD_STATUS_IDLE,
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
			"audition_status": models.PURCHASE_RECORD_STATUS_APPLY,
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
			Message: "别着急...助教正在定制你的课程并为你匹配合适的导师",
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
			createDeluxeCourseOrder(record.Id)
		}

	case record.AuditionStatus == models.PURCHASE_RECORD_STATUS_COMPLETE &&
		record.PurchaseStatus == models.PURCHASE_RECORD_STATUS_IDLE,
		record.AuditionStatus == models.PURCHASE_RECORD_STATUS_COMPLETE &&
			record.PurchaseStatus == models.PURCHASE_RECORD_STATUS_WAITING:

		// 学生已经完成试听课程，学生须支付课程费用

		recordInfo := map[string]interface{}{
			"purchase_status": models.PURCHASE_RECORD_STATUS_WAITING,
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
			"audition_status": models.PURCHASE_RECORD_STATUS_WAITING,
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
			createDeluxeCourseOrder(record.Id)
		}

	case record.PurchaseStatus == models.PURCHASE_RECORD_STATUS_COMPLETE:

		// 学生的课程已经完成，无法继续操作
		response = actionProceedResponse{
			Action:  ACTION_PROCEED_NULL,
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

//处理新版本的试听
func HandleAuditionCourseActionProceed(userId int64, course *models.Course, sourceCourseId int64) (int64, *actionProceedResponse) {
	var err error
	courseId := course.Id
	o := orm.NewOrm()
	var auditionRecordStatus string
	// 先查询该用户是否有未完成的试听
	var currentRecord models.CourseAuditionRecord
	err = o.QueryTable(new(models.CourseAuditionRecord).TableName()).
		Filter("course_id", courseId).Filter("user_id", userId).
		Exclude("status", models.AUDITION_RECORD_STATUS_COMPLETE).
		One(&currentRecord)
	auditionRecordStatus = currentRecord.Status
	if err == orm.ErrNoRows {
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
	} else if err != nil {
		// 如果到了这里说明数据库报错了...
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
