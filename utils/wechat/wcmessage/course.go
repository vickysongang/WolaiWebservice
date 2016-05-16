package wcmessage

import (
	//	"encoding/json"
	"fmt"
	//"math"

	"WolaiWebservice/models"
	"WolaiWebservice/service/evaluation"
	"WolaiWebservice/utils/wechat"
	"encoding/json"
)

func SendChapterCompleteNotification(recordId, chapterId int64, courseName string, courseType string) {
	var err error
	var studentId, teacherId int64

	if courseType == models.COURSE_TYPE_DELUXE {
		record, err := models.ReadCoursePurchaseRecord(recordId)
		if err != nil {
			return
		}
		studentId = record.UserId
		teacherId = record.TeacherId
	} else if courseType == models.COURSE_TYPE_AUDITION {
		record, err := models.ReadCourseAuditionRecord(recordId)
		if err != nil {
			return
		}
		studentId = record.UserId
		teacherId = record.TeacherId
	}

	student, err := models.ReadUser(studentId)
	if err != nil {
		return
	}

	tutor, err := models.ReadUser(teacherId)
	if err != nil {
		return
	}

	wxBind, err := models.QueryUserWxQauthByUserId(studentId)
	if err != nil {
		return
	}

	wxId := wxBind.WxId

	wcMsg := wechat.WCMessage{}
	wcMsg.ToUser = wxId
	wcMsg.Url = fmt.Sprintf("%s%d/%d", evaluation.GetEvaluationDetailUrlPrefix(), chapterId, recordId)
	wcMsg.MsgType = "COURSESUMMARY"
	wcData := make(map[string]wechat.WCField)
	wcData["first"] = wechat.WCField{"您收到了一条课后评价"}
	wcData["keyword1"] = wechat.WCField{student.Nickname}
	wcData["keyword2"] = wechat.WCField{courseName}
	wcData["keyword3"] = wechat.WCField{tutor.Nickname}
	wcData["remark"] = wechat.WCField{"请点击详情查看课时总结。"}
	wcDataByte, _ := json.Marshal(&wcData)
	wcMsg.Data = string(wcDataByte)

	wechat.WCSendTypedMessage(&wcMsg)

}
