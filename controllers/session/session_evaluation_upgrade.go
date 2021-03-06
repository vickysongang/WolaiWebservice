// session_evaluation_upgrade
package session

import (
	"WolaiWebservice/models"
	evaluationService "WolaiWebservice/service/evaluation"
	"strings"
)

func CreateEvaluationUpgrade(userId, sessionId, chapterId, recordId int64, evaluationType, evaluationContent string) (*models.Evaluation, error) {
	user, _ := models.ReadUser(userId)
	if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER { //导师插入评价
		switch evaluationType {
		case "course":
			apply, _ := evaluationService.GetEvaluationApply(userId, chapterId, recordId)
			chapter, _ := models.ReadCourseCustomChapter(chapterId)
			if apply.Id == 0 {
				evaluationApply := models.EvaluationApply{
					UserId:    userId,
					SessionId: sessionId,
					CourseId:  chapter.CourseId,
					ChapterId: chapterId,
					Status:    models.EVALUATION_APPLY_STATUS_CREATED,
					Content:   evaluationContent,
					RecordId:  recordId,
				}
				_, err := models.InsertEvaluationApply(&evaluationApply)
				if err != nil {
					return nil, err
				}
			} else {
				if apply.Status == models.EVALUATION_APPLY_STATUS_IDLE {
					eveluationInfo := map[string]interface{}{
						"SessionId": sessionId,
						"Status":    models.EVALUATION_APPLY_STATUS_CREATED,
						"Content":   evaluationContent,
					}
					models.UpdateEvaluationApply(apply.Id, eveluationInfo)
				}
			}
		case "qa":
			content, err := evaluateSessionUpgrade(sessionId, userId, chapterId, recordId, evaluationContent)
			return content, err
		}
	} else if user.AccessRight == models.USER_ACCESSRIGHT_STUDENT { //学生插入评价
		content, err := evaluateSessionUpgrade(sessionId, userId, chapterId, recordId, evaluationContent)
		return content, err
	}
	return nil, nil
}

func evaluateSessionUpgrade(sessionId, userId, chapterId, recordId int64, evaluationContent string) (*models.Evaluation, error) {
	session, _ := models.ReadSession(sessionId)
	var targetId int64
	if userId == session.Creator {
		targetId = session.Tutor
	} else {
		targetId = session.Creator
	}
	oldEvaluation, _ := evaluationService.QueryEvaluation(userId, sessionId)
	if oldEvaluation.Id == 0 {
		evaluation := models.Evaluation{
			UserId:    userId,
			TargetId:  targetId,
			SessionId: sessionId,
			ChapterId: chapterId,
			Content:   evaluationContent,
			RecordId:  recordId}
		content, err := models.InsertEvaluation(&evaluation)
		return content, err
	} else {
		evaluationInfo := map[string]interface{}{
			"UserId":    userId,
			"TargetId":  targetId,
			"SessionId": sessionId,
			"ChapterId": chapterId,
			"Content":   evaluationContent,
			"RecordId":  recordId,
		}
		err := models.UpdateEvaluation(oldEvaluation.Id, evaluationInfo)
		if err != nil {
			return nil, err
		} else {
			content, err := models.ReadEvaluation(oldEvaluation.Id)
			return content, err
		}
	}
	return nil, nil
}

func QueryEvaluationInfoUpgrade(userId, sessionId, chapterId, recordId int64) ([]*evaluationInfo, error) {
	evalutionInfos := make([]*evaluationInfo, 0)
	selfEvaluation := evaluationInfo{}
	otherEvaluation := evaluationInfo{}
	var isStudent bool
	var studentEvaluation, teacherEvaluation *models.Evaluation
	if sessionId != 0 {
		session, _ := models.ReadSession(sessionId)
		studentEvaluation, _ = evaluationService.QueryEvaluation(session.Creator, sessionId)
		teacherEvaluation, _ = evaluationService.QueryEvaluation(session.Tutor, sessionId)
		if userId == session.Creator {
			isStudent = true
		} else if userId == session.Tutor {
			isStudent = false
		}
	} else {
		chapter, _ := models.ReadCourseCustomChapter(chapterId)
		studentEvaluation, _ = evaluationService.QueryEvaluationByChapter(chapter.UserId, chapterId, recordId)
		teacherEvaluation, _ = evaluationService.QueryEvaluationByChapter(chapter.TeacherId, chapterId, recordId)
		if userId == chapter.UserId {
			isStudent = true
		} else if userId == chapter.TeacherId {
			isStudent = false
		}
	}

	//旧版评价表里targetId为0，新版不为0，故根据该字段来判断获取的是旧版评论还是新版评论
	if isStudent {
		//导师和学生都有评论
		if teacherEvaluation.Id != 0 && studentEvaluation.Id != 0 {
			if (teacherEvaluation.TargetId != 0 && studentEvaluation.TargetId == 0) || (teacherEvaluation.TargetId == 0 && studentEvaluation.TargetId != 0) {
				if !(strings.HasPrefix(studentEvaluation.Content, "[") && strings.HasSuffix(studentEvaluation.Content, "]")) {
					selfEvaluation.Type = "student"
					selfEvaluation.Evalution = studentEvaluation
					evalutionInfos = append(evalutionInfos, &selfEvaluation)
				}
			} else {
				if !(strings.HasPrefix(studentEvaluation.Content, "[") && strings.HasSuffix(studentEvaluation.Content, "]")) {
					selfEvaluation.Type = "student"
					selfEvaluation.Evalution = studentEvaluation
					evalutionInfos = append(evalutionInfos, &selfEvaluation)
				}
				if !(strings.HasPrefix(teacherEvaluation.Content, "[") && strings.HasSuffix(teacherEvaluation.Content, "]")) {
					otherEvaluation.Type = "teacher"
					otherEvaluation.Evalution = teacherEvaluation
					evalutionInfos = append(evalutionInfos, &otherEvaluation)
				}
			}
		} else if teacherEvaluation.Id == 0 && studentEvaluation.Id != 0 { //导师未评论，学生有评论
			if !(strings.HasPrefix(studentEvaluation.Content, "[") && strings.HasSuffix(studentEvaluation.Content, "]")) {
				selfEvaluation.Type = "student"
				selfEvaluation.Evalution = studentEvaluation
				evalutionInfos = append(evalutionInfos, &selfEvaluation)
			}
		} else if teacherEvaluation.Id != 0 && studentEvaluation.Id == 0 { //导师有评论，学生未评论
			if teacherEvaluation.TargetId != 0 {
				otherEvaluation.Type = "teacher"
				otherEvaluation.Evalution = teacherEvaluation
				evalutionInfos = append(evalutionInfos, &otherEvaluation)
			}
		}
	} else {
		if teacherEvaluation.Id != 0 { //导师只能看到自己的评论，不能看到学生的评论
			if !(strings.HasPrefix(teacherEvaluation.Content, "[") && strings.HasSuffix(teacherEvaluation.Content, "]")) {
				selfEvaluation.Type = "teacher"
				selfEvaluation.Evalution = teacherEvaluation
				evalutionInfos = append(evalutionInfos, &selfEvaluation)
			}
		}
	}
	return evalutionInfos, nil
}
