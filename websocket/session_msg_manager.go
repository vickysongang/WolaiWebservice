// session_msg_handler
package websocket

import (
	sessionController "WolaiWebservice/controllers/session"
	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
	"WolaiWebservice/service/push"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/cihub/seelog"
)

var ErrUserChanClose = errors.New("user chan closes")

func SendBreakMsgToStudent(studentId, teacherId, sessionId int64) error {
	sessionIdStr := strconv.FormatInt(sessionId, 10)
	breakMsg := NewWSMessage("", studentId, WS_SESSION_BREAK)
	breakMsg.Attribute["sessionId"] = sessionIdStr
	breakMsg.Attribute["studentId"] = strconv.FormatInt(studentId, 10)
	breakMsg.Attribute["teacherId"] = strconv.FormatInt(teacherId, 10)
	length, _ := SessionManager.GetSessionLength(sessionId)
	breakMsg.Attribute["timer"] = strconv.FormatInt(length, 10)
	breakMsg.Attribute["sessionStatus"] = SESSION_STATUS_BREAKED
	if UserManager.HasUserChan(breakMsg.UserId) {
		breakChan := UserManager.GetUserChan(breakMsg.UserId)
		breakChan <- breakMsg
		return nil
	}
	return ErrUserChanClose
}

func SendBreakMsgToTeacher(studentId, teacherId, sessionId int64) error {
	sessionIdStr := strconv.FormatInt(sessionId, 10)
	breakMsg := NewWSMessage("", teacherId, WS_SESSION_BREAK)
	breakMsg.Attribute["sessionId"] = sessionIdStr
	breakMsg.Attribute["studentId"] = strconv.FormatInt(studentId, 10)
	breakMsg.Attribute["teacherId"] = strconv.FormatInt(teacherId, 10)
	length, _ := SessionManager.GetSessionLength(sessionId)
	breakMsg.Attribute["timer"] = strconv.FormatInt(length, 10)
	breakMsg.Attribute["sessionStatus"] = SESSION_STATUS_BREAKED
	if UserManager.HasUserChan(breakMsg.UserId) {
		breakChan := UserManager.GetUserChan(breakMsg.UserId)
		breakChan <- breakMsg
		return nil
	}
	return ErrUserChanClose
}

func SendPauseRespMsgToTeacherOnError(msgId string, teacherId int64) error {
	pauseResp := NewWSMessage(msgId, teacherId, WS_SESSION_PAUSE_RESP)
	pauseResp.Attribute["errCode"] = "2"
	if UserManager.HasUserChan(teacherId) {
		userChan := UserManager.GetUserChan(teacherId)
		userChan <- pauseResp
		return nil
	}
	return ErrUserChanClose
}

func SendPauseRespMsgToTeacher(msgId string, teacherId, sessionId int64) error {
	sessionIdStr := strconv.FormatInt(sessionId, 10)
	pauseResp := NewWSMessage(msgId, teacherId, WS_SESSION_PAUSE_RESP)
	pauseResp.Attribute["errCode"] = "0"
	pauseResp.Attribute["sessionId"] = sessionIdStr
	pauseResp.Attribute["sessionStatus"] = SESSION_STATUS_PAUSED
	if UserManager.HasUserChan(teacherId) {
		userChan := UserManager.GetUserChan(teacherId)
		userChan <- pauseResp
		return nil
	} else {
		seelog.Debugf("session pause when start sessionId: %d, tutor userChan closes userId: %d", sessionId, teacherId)
	}
	return ErrUserChanClose
}

func SendPauseMsgToStudent(studentId, teacherId, sessionId, length int64) error {
	sessionIdStr := strconv.FormatInt(sessionId, 10)
	pauseMsg := NewWSMessage("", studentId, WS_SESSION_PAUSE)
	pauseMsg.Attribute["sessionId"] = sessionIdStr
	pauseMsg.Attribute["teacherId"] = strconv.FormatInt(teacherId, 10)
	pauseMsg.Attribute["timer"] = strconv.FormatInt(length, 10)
	pauseMsg.Attribute["sessionStatus"] = SESSION_STATUS_PAUSED
	if UserManager.HasUserChan(studentId) {
		studentChan := UserManager.GetUserChan(studentId)
		studentChan <- pauseMsg
		return nil
	} else {
		seelog.Debugf("session pause when start sessionId: %d, student userChan closes userId: %d", sessionId, studentId)
	}
	return ErrUserChanClose
}

func SendExpireMsg(userId, sessionId int64) error {
	sessionIdStr := strconv.FormatInt(sessionId, 10)
	expireMsg := NewWSMessage("", userId, WS_SESSION_EXPIRE)
	expireMsg.Attribute["sessionId"] = sessionIdStr
	expireMsg.Attribute["sessionStatus"] = SESSION_STATUS_COMPLETE
	if UserManager.HasUserChan(userId) {
		userChan := UserManager.GetUserChan(userId)
		userChan <- expireMsg
		return nil
	}
	return ErrUserChanClose
}

func SendSyncMsg(userId, sessionId, length int64) error {
	sessionIdStr := strconv.FormatInt(sessionId, 10)
	syncMsg := NewWSMessage("", userId, WS_SESSION_SYNC)
	syncMsg.Attribute["sessionId"] = sessionIdStr
	syncMsg.Attribute["timer"] = strconv.FormatInt(length, 10)
	status, err := SessionManager.GetSessionStatus(sessionId)
	if err != nil {
		seelog.Debugf("GetSesssionStatus failed sessionId: %d ,error: %s", sessionId, err.Error())
	}
	syncMsg.Attribute["sessionStatus"] = status
	if UserManager.HasUserChan(userId) {
		userChan := UserManager.GetUserChan(userId)
		userChan <- syncMsg
		return nil
	}
	return ErrUserChanClose
}

func SendFinishMsgToStudent(studentId, sessionId int64) error {
	sessionIdStr := strconv.FormatInt(sessionId, 10)
	finishMsg := NewWSMessage("", studentId, WS_SESSION_FINISH)
	finishMsg.Attribute["sessionId"] = sessionIdStr
	finishMsg.Attribute["sessionStatus"] = SESSION_STATUS_COMPLETE
	if UserManager.HasUserChan(studentId) {
		creatorChan := UserManager.GetUserChan(studentId)
		creatorChan <- finishMsg
		return nil
	}
	return ErrUserChanClose
}

func SendFinishRespMsgToTeacherOnError(msgId string, userId int64) error {
	finishResp := NewWSMessage(msgId, userId, WS_SESSION_FINISH_RESP)
	finishResp.Attribute["errCode"] = "2"
	finishResp.Attribute["errMsg"] = "You are not the teacher of this session"
	if UserManager.HasUserChan(userId) {
		userChan := UserManager.GetUserChan(userId)
		userChan <- finishResp
		return nil
	}
	return ErrUserChanClose
}

func SendFinishRespMsgToTeacher(msgId string, userId, sessionId int64) error {
	sessionIdStr := strconv.FormatInt(sessionId, 10)
	finishResp := NewWSMessage(msgId, userId, WS_SESSION_FINISH_RESP)
	finishResp.Attribute["errCode"] = "0"
	finishResp.Attribute["sessionId"] = sessionIdStr
	finishResp.Attribute["sessionStatus"] = SESSION_STATUS_COMPLETE
	if UserManager.HasUserChan(userId) {
		userChan := UserManager.GetUserChan(userId)
		userChan <- finishResp
		return nil
	} else {
		seelog.Debug("session finish: userChan closes | sessionHandler:", sessionId)
	}
	return ErrUserChanClose
}

func SendRecoverMsgToTeacher(studentId, teacherId, sessionId, length int64, order *models.Order) error {
	sessionIdStr := strconv.FormatInt(sessionId, 10)
	recoverTeacherMsg := NewWSMessage("", teacherId, WS_SESSION_RECOVER_TEACHER)
	recoverTeacherMsg.Attribute["sessionId"] = sessionIdStr
	recoverTeacherMsg.Attribute["studentId"] = strconv.FormatInt(studentId, 10)
	recoverTeacherMsg.Attribute["timer"] = strconv.FormatInt(length, 10)
	if order.Type == models.ORDER_TYPE_COURSE_INSTANT {
		courseRelation, _ := courseService.GetCourseRelation(order.RecordId, models.COURSE_TYPE_DELUXE)
		virturlCourseId := courseRelation.Id
		recoverTeacherMsg.Attribute["courseId"] = strconv.FormatInt(virturlCourseId, 10)
	} else if order.Type == models.ORDER_TYPE_AUDITION_COURSE_INSTANT {
		courseRelation, _ := courseService.GetCourseRelation(order.RecordId, models.COURSE_TYPE_AUDITION)
		virturlCourseId := courseRelation.Id
		recoverTeacherMsg.Attribute["courseId"] = strconv.FormatInt(virturlCourseId, 10)
	}
	sessionStatus, _ := SessionManager.GetSessionStatus(sessionId)
	recoverTeacherMsg.Attribute["sessionStatus"] = sessionStatus
	if UserManager.HasUserChan(teacherId) {
		teacherChan := UserManager.GetUserChan(teacherId)
		teacherChan <- recoverTeacherMsg
		return nil
	}
	return ErrUserChanClose
}

func SendRecoverMsgToStudent(studentId, teacherId, sessionId, length int64, order *models.Order) error {
	sessionIdStr := strconv.FormatInt(sessionId, 10)
	recoverStuMsg := NewWSMessage("", studentId, WS_SESSION_RECOVER_STU)
	recoverStuMsg.Attribute["sessionId"] = sessionIdStr
	recoverStuMsg.Attribute["teacherId"] = strconv.FormatInt(teacherId, 10)
	recoverStuMsg.Attribute["timer"] = strconv.FormatInt(length, 10)
	if order.Type == models.ORDER_TYPE_COURSE_INSTANT || order.Type == models.ORDER_TYPE_AUDITION_COURSE_INSTANT {
		recoverStuMsg.Attribute["courseId"] = strconv.FormatInt(order.CourseId, 10)
	}
	sessionStatus, _ := SessionManager.GetSessionStatus(sessionId)
	recoverStuMsg.Attribute["sessionStatus"] = sessionStatus
	if UserManager.HasUserChan(studentId) {
		studentChan := UserManager.GetUserChan(studentId)
		studentChan <- recoverStuMsg
		return nil
	}
	return ErrUserChanClose
}

func SendBreakReconnectSuccessMsgToTeacher(studentId, teacherId, sessionId, length int64) error {
	sessionIdStr := strconv.FormatInt(sessionId, 10)
	sessionStatusMsg := NewWSMessage("", teacherId, WS_SESSION_BREAK_RECONNECT_SUCCESS)
	sessionStatusMsg.Attribute["sessionId"] = sessionIdStr
	sessionStatusMsg.Attribute["studentId"] = strconv.FormatInt(studentId, 10)
	sessionStatusMsg.Attribute["teacherId"] = strconv.FormatInt(teacherId, 10)
	sessionStatusMsg.Attribute["timer"] = strconv.FormatInt(length, 10)
	sessionStatus, _ := SessionManager.GetSessionStatus(sessionId)
	sessionStatusMsg.Attribute["sessionStatus"] = sessionStatus
	if UserManager.HasUserChan(teacherId) {
		teacherChan := UserManager.GetUserChan(teacherId)
		teacherChan <- sessionStatusMsg
		return nil
	}
	return ErrUserChanClose
}

func SendBreakReconnectSuccessMsgToStudent(studentId, teacherId, sessionId, length int64) error {
	sessionIdStr := strconv.FormatInt(sessionId, 10)
	sessionStatusMsg := NewWSMessage("", studentId, WS_SESSION_BREAK_RECONNECT_SUCCESS)
	sessionStatusMsg.Attribute["sessionId"] = sessionIdStr
	sessionStatusMsg.Attribute["studentId"] = strconv.FormatInt(studentId, 10)
	sessionStatusMsg.Attribute["teacherId"] = strconv.FormatInt(teacherId, 10)
	sessionStatusMsg.Attribute["timer"] = strconv.FormatInt(length, 10)
	sessionStatus, _ := SessionManager.GetSessionStatus(sessionId)
	sessionStatusMsg.Attribute["sessionStatus"] = sessionStatus
	if UserManager.HasUserChan(studentId) {
		studentChan := UserManager.GetUserChan(studentId)
		studentChan <- sessionStatusMsg
		return nil
	}
	return ErrUserChanClose
}

func SendStatusSyncMsg(userId, sessionId int64) error {
	syncMsg := NewWSMessage("", userId, WS_SESSION_STATUS_SYNC)
	syncMsg.Attribute["errCode"] = "0"
	sessionStatus, _ := SessionManager.GetSessionStatus(sessionId)
	syncMsg.Attribute["sessionStatus"] = sessionStatus
	_, userInfo := sessionController.GetSessionInfo(sessionId, userId)
	userInfoByte, _ := json.Marshal(userInfo)
	syncMsg.Attribute["sessionInfo"] = string(userInfoByte)

	if UserManager.HasUserChan(userId) {
		userChan := UserManager.GetUserChan(userId)
		userChan <- syncMsg
		return nil
	}
	return ErrUserChanClose
}

func SendResumeRespMsgToTeacherOnError(msgId string, teacherId int64, errMsg string) error {
	resumeResp := NewWSMessage(msgId, teacherId, WS_SESSION_RESUME_RESP)
	resumeResp.Attribute["errCode"] = "2"
	resumeResp.Attribute["errMsg"] = errMsg
	if UserManager.HasUserChan(teacherId) {
		userChan := UserManager.GetUserChan(teacherId)
		userChan <- resumeResp
		return nil
	}
	return ErrUserChanClose
}

func SendResumeRespMsgToTeacher(msgId string, teacherId, sessionId int64) error {
	sessionIdStr := strconv.FormatInt(sessionId, 10)
	resumeResp := NewWSMessage(msgId, teacherId, WS_SESSION_RESUME_RESP)
	resumeResp.Attribute["errCode"] = "0"
	resumeResp.Attribute["sessionId"] = sessionIdStr
	resumeResp.Attribute["sessionStatus"] = SESSION_STATUS_CALLING
	if UserManager.HasUserChan(teacherId) {
		userChan := UserManager.GetUserChan(teacherId)
		userChan <- resumeResp
		return nil
	} else {
		seelog.Debug("session resume: userChan closes | sessionHandler:", sessionId)
	}
	return ErrUserChanClose
}

func SendResumeMsgToStudent(studentId, teacherId, sessionId int64) error {
	sessionIdStr := strconv.FormatInt(sessionId, 10)
	resumeMsg := NewWSMessage("", studentId, WS_SESSION_RESUME)
	resumeMsg.Attribute["sessionId"] = sessionIdStr
	resumeMsg.Attribute["teacherId"] = strconv.FormatInt(teacherId, 10)
	resumeMsg.Attribute["sessionStatus"] = SESSION_STATUS_CALLING
	if UserManager.HasUserChan(studentId) {
		studentChan := UserManager.GetUserChan(studentId)
		studentChan <- resumeMsg
	} else {
		push.PushSessionResume(studentId, sessionId)
	}
	return nil
}

func SendResumeCancelRespMsgToTeacherOnError(msgId string, teacherId int64) error {
	resCancelResp := NewWSMessage(msgId, teacherId, WS_SESSION_RESUME_CANCEL_RESP)
	resCancelResp.Attribute["errCode"] = "2"
	resCancelResp.Attribute["errMsg"] = "nobody is calling"
	if UserManager.HasUserChan(teacherId) {
		userChan := UserManager.GetUserChan(teacherId)
		userChan <- resCancelResp
		return nil
	}
	return ErrUserChanClose
}

func SendResumeCancelRespMsgToTeacher(msgId string, teacherId, sessionId int64, sessionStatus string) error {
	sessionIdStr := strconv.FormatInt(sessionId, 10)
	resCancelResp := NewWSMessage(msgId, teacherId, WS_SESSION_RESUME_CANCEL_RESP)
	resCancelResp.Attribute["errCode"] = "0"
	resCancelResp.Attribute["sessionId"] = sessionIdStr
	resCancelResp.Attribute["sessionStatus"] = sessionStatus
	if UserManager.HasUserChan(teacherId) {
		userChan := UserManager.GetUserChan(teacherId)
		userChan <- resCancelResp
		return nil
	} else {
		seelog.Debug("session resume cancel: userChan closes | sessionHandler:", sessionId)
	}
	return ErrUserChanClose
}

func SendResumeCancelRespMsgToStudent(studentId, teacherId, sessionId int64, sessionStatus string) error {
	sessionIdStr := strconv.FormatInt(sessionId, 10)
	resCancelMsg := NewWSMessage("", studentId, WS_SESSION_RESUME_CANCEL)
	resCancelMsg.Attribute["sessionId"] = sessionIdStr
	resCancelMsg.Attribute["teacherId"] = strconv.FormatInt(teacherId, 10)
	resCancelMsg.Attribute["sessionStatus"] = sessionStatus
	if UserManager.HasUserChan(studentId) {
		studentChan := UserManager.GetUserChan(studentId)
		studentChan <- resCancelMsg
		return nil
	}
	return ErrUserChanClose
}

func SendResumeAcceptRespMsgToStudentOnError(msgId string, studentId int64, errMsg string) error {
	resAcceptResp := NewWSMessage(msgId, studentId, WS_SESSION_RESUME_ACCEPT_RESP)
	resAcceptResp.Attribute["errCode"] = "2"
	resAcceptResp.Attribute["errMsg"] = errMsg
	if UserManager.HasUserChan(studentId) {
		userChan := UserManager.GetUserChan(studentId)
		userChan <- resAcceptResp
		return nil
	}
	return ErrUserChanClose
}

func SendResumeAcceptRespMsgToStudent(msgId string, studentId, sessionId int64, sessionStatus string) error {
	sessionIdStr := strconv.FormatInt(sessionId, 10)
	resAcceptResp := NewWSMessage(msgId, studentId, WS_SESSION_RESUME_ACCEPT_RESP)
	resAcceptResp.Attribute["errCode"] = "0"
	resAcceptResp.Attribute["sessionId"] = sessionIdStr
	resAcceptResp.Attribute["sessionStatus"] = sessionStatus
	if UserManager.HasUserChan(studentId) {
		userChan := UserManager.GetUserChan(studentId)
		userChan <- resAcceptResp
		return nil
	} else {
		seelog.Debug("session resume accept: userChan closes | sessionHandler:", sessionId)
	}
	return ErrUserChanClose
}

func SendResumeAcceptMsgToTeacher(teacherId, sessionId int64, acceptStr string, sessionStatus string) error {
	sessionIdStr := strconv.FormatInt(sessionId, 10)
	resAcceptMsg := NewWSMessage("", teacherId, WS_SESSION_RESUME_ACCEPT)
	resAcceptMsg.Attribute["sessionId"] = sessionIdStr
	resAcceptMsg.Attribute["accept"] = acceptStr
	resAcceptMsg.Attribute["sessionStatus"] = sessionStatus
	if UserManager.HasUserChan(teacherId) {
		teacherChan := UserManager.GetUserChan(teacherId)
		teacherChan <- resAcceptMsg
		return nil
	}
	return ErrUserChanClose
}

func SendQaPkgTimeEndMsgToStudent(studentId, sessionId int64) error {
	sessionIdStr := strconv.FormatInt(sessionId, 10)
	qaPkgTimeEndMsg := NewWSMessage("", studentId, WS_SESSION_QAPKG_TIME_END)
	qaPkgTimeEndMsg.Attribute["sessionId"] = sessionIdStr
	qaPkgTimeEndMsg.Attribute["comment"] = "答疑时间用完啦，本次上课已经换到账户余额支付"
	qaPkgTimeEndMsg.Attribute["sessionStatus"] = SESSION_STATUS_SERVING
	if UserManager.HasUserChan(studentId) {
		userChan := UserManager.GetUserChan(studentId)
		userChan <- qaPkgTimeEndMsg
		return nil
	}
	return ErrUserChanClose
}

func SendAutoFinishTipMsgToStudent(studentId, sessionId, autoFinishLimit int64) error {
	sessionIdStr := strconv.FormatInt(sessionId, 10)
	autoFinishTipMsg := NewWSMessage("", studentId, WS_SESSION_AUTO_FINISH_TIP)
	autoFinishTipMsg.Attribute["sessionId"] = sessionIdStr
	autoFinishTipMsg.Attribute["comment"] = fmt.Sprintf("%s%d%s", "哎呀！账户里的钱都用完了", autoFinishLimit, "分钟后将自动下课哦")
	autoFinishTipMsg.Attribute["sessionStatus"] = SESSION_STATUS_SERVING
	if UserManager.HasUserChan(studentId) {
		userChan := UserManager.GetUserChan(studentId)
		userChan <- autoFinishTipMsg
		return nil
	}
	return ErrUserChanClose
}

func SendAutoFinishTipMsgToTeacher(teacherId, sessionId, autoFinishLimit int64) error {
	sessionIdStr := strconv.FormatInt(sessionId, 10)
	autoFinishTipMsg := NewWSMessage("", teacherId, WS_SESSION_AUTO_FINISH_TIP)
	autoFinishTipMsg.Attribute["sessionId"] = sessionIdStr
	autoFinishTipMsg.Attribute["comment"] = "学生余额不足"
	autoFinishTipMsg.Attribute["sessionStatus"] = SESSION_STATUS_SERVING
	if UserManager.HasUserChan(teacherId) {
		userChan := UserManager.GetUserChan(teacherId)
		userChan <- autoFinishTipMsg
		return nil
	}
	return ErrUserChanClose
}

func SendAssignOffMsgToTeacher(teacherId int64) error {
	assignOffMsg := NewWSMessage("", teacherId, WS_ORDER2_TEACHER_ASSIGNOFF_RESP)
	assignOffMsg.Attribute["errCode"] = "0"
	if UserManager.HasUserChan(teacherId) {
		tutorChan := UserManager.GetUserChan(teacherId)
		tutorChan <- assignOffMsg
		return nil
	}
	return ErrUserChanClose
}
