package session

import (
	"math/rand"
	"time"

	evaluationService "WolaiWebservice/service/evaluation"

	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
)

type evaluationInfo struct {
	Type      string             `json:"type"`
	Evalution *models.Evaluation `json:"evaluationInfo"`
}

func CheckRandNumInSlice(slice []int64, randNum int64) bool {
	for _, v := range slice {
		if v == randNum {
			return true
		}
	}
	return false
}

func GetRandNumSlice(sliceSize int64, length int64) []int64 {
	var result []int64
	if sliceSize <= 0 || length <= 0 {
		return result
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		randNum := r.Int63n(length)
		if !CheckRandNumInSlice(result, randNum) {
			result = append(result, randNum)
		}
		if sliceSize > length {
			if int64(len(result)) == length {
				break
			}
		} else {
			if int64(len(result)) == sliceSize {
				break
			}
		}
	}
	return result
}

func QuerySystemEvaluationLabels(userId, sessionId, count int64) ([]*models.EvaluationLabel, error) {
	labels := []*models.EvaluationLabel{}
	session, _ := models.ReadSession(sessionId)

	//如果当前用户是学生，则要返回老师的标签信息，如果当前用户是老师，则要返回学生的标签信息
	//学生
	if userId == session.Creator {
		teacher, _ := models.ReadUser(session.Tutor)
		//个人标签
		teacherPersonalLabels, err := models.QueryEvaluationLabels(teacher.Gender, models.PERSONAL_EVALUATION_LABEL, models.TEACHER_EVALUATION_LABEL)
		if err != nil {
			return nil, err
		}
		personalCount := (count - 2) / 2
		for _, v := range GetRandNumSlice(personalCount, int64(len(teacherPersonalLabels))) {
			labels = append(labels, teacherPersonalLabels[v])
		}
		//讲课风格
		teacherStyleLabels, err := models.QueryEvaluationLabels(teacher.Gender, models.STYLE_EVALUATION_LABEL, models.TEACHER_EVALUATION_LABEL)
		if err != nil {
			return nil, err
		}
		for _, v := range GetRandNumSlice(count-personalCount-2, int64(len(teacherStyleLabels))) {
			labels = append(labels, teacherStyleLabels[v])
		}
		//科目标签
		order, _ := models.ReadOrder(session.OrderId)
		teacherSubjectLabels, err := models.QueryEvaluationLabelsBySubject(order.SubjectId)
		if err != nil {
			return nil, err
		}
		for _, v := range GetRandNumSlice(2, int64(len(teacherSubjectLabels))) {
			labels = append(labels, teacherSubjectLabels[v])
		}
	} else if userId == session.Tutor { //老师
		student, _ := models.ReadUser(session.Creator)
		//个人标签
		studentPersonalLabels, err := models.QueryEvaluationLabels(student.Gender, models.PERSONAL_EVALUATION_LABEL, models.STUDENT_EVALUATION_LABEL)
		if err != nil {
			return nil, err
		}
		personalCount := count / 2
		for _, v := range GetRandNumSlice(personalCount, int64(len(studentPersonalLabels))) {
			labels = append(labels, studentPersonalLabels[v])
		}
		//能力程度
		studentAbilityLabels, err := models.QueryEvaluationLabels(student.Gender, models.ABILITY_EVALUATION_LABEL, models.STUDENT_EVALUATION_LABEL)
		if err != nil {
			return nil, err
		}
		for _, v := range GetRandNumSlice(count-personalCount, int64(len(studentAbilityLabels))) {
			labels = append(labels, studentAbilityLabels[v])
		}
	}
	return labels, nil
}

func CreateEvaluation(userId, targetId, sessionId int64, evaluationContent string) (*models.Evaluation, error) {
	session, _ := models.ReadSession(sessionId)
	order, _ := models.ReadOrder(session.OrderId)
	if order.CourseId != 0 && userId == session.Tutor {
		apply, _ := evaluationService.GetEvaluationApply(userId, order.ChapterId)
		if apply.Id == 0 {
			evaluationApply := models.EvaluationApply{
				UserId:    userId,
				SessionId: sessionId,
				CourseId:  order.CourseId,
				ChapterId: order.ChapterId,
				Status:    models.EVALUATION_APPLY_STATUS_CREATED,
				Content:   evaluationContent,
			}
			models.InsertEvaluationApply(&evaluationApply)
		}
		return nil, nil
	} else {
		if targetId == 1 {
			if userId == session.Creator {
				targetId = session.Tutor
			} else {
				targetId = session.Creator
			}
		}
		evaluation := models.Evaluation{
			UserId:    userId,
			TargetId:  targetId,
			SessionId: sessionId,
			Content:   evaluationContent}
		content, err := models.InsertEvaluation(&evaluation)
		return content, err
	}
	return nil, nil
}

func QueryEvaluationInfo(userId, sessionId int64) ([]*evaluationInfo, error) {
	session, _ := models.ReadSession(sessionId)
	studentEvaluation, _ := models.QueryEvaluation(session.Creator, sessionId)
	teacherEvaluation, _ := models.QueryEvaluation(session.Tutor, sessionId)

	selfEvaluation := evaluationInfo{}
	otherEvaluation := evaluationInfo{}

	evalutionInfos := make([]*evaluationInfo, 0)
	//旧版评价表里targetId为0，新版不为0，故根据该字段来判断获取的是旧版评论还是新版评论
	if userId == session.Creator {
		if teacherEvaluation.Id != 0 && studentEvaluation.Id != 0 {
			if (teacherEvaluation.TargetId != 0 && studentEvaluation.TargetId == 0) || (teacherEvaluation.TargetId == 0 && studentEvaluation.TargetId != 0) {
				selfEvaluation.Type = "student"
				selfEvaluation.Evalution = studentEvaluation
				evalutionInfos = append(evalutionInfos, &selfEvaluation)
			} else {
				selfEvaluation.Type = "student"
				selfEvaluation.Evalution = studentEvaluation
				evalutionInfos = append(evalutionInfos, &selfEvaluation)

				otherEvaluation.Type = "teacher"
				otherEvaluation.Evalution = teacherEvaluation
				evalutionInfos = append(evalutionInfos, &otherEvaluation)
			}
		}
	} else if userId == session.Tutor {
		if teacherEvaluation.Id != 0 && studentEvaluation.Id != 0 {
			if (teacherEvaluation.TargetId != 0 && studentEvaluation.TargetId == 0) || (teacherEvaluation.TargetId == 0 && studentEvaluation.TargetId != 0) {
				selfEvaluation.Type = "teacher"
				selfEvaluation.Evalution = teacherEvaluation
				evalutionInfos = append(evalutionInfos, &selfEvaluation)
			} else {
				selfEvaluation.Type = "teacher"
				selfEvaluation.Evalution = teacherEvaluation
				evalutionInfos = append(evalutionInfos, &selfEvaluation)

				otherEvaluation.Type = "student"
				otherEvaluation.Evalution = studentEvaluation
				evalutionInfos = append(evalutionInfos, &otherEvaluation)
			}
		}
	}
	return evalutionInfos, nil
}

func HasOrderInSessionEvaluated(sessionId int64, userId int64) bool {
	o := orm.NewOrm()
	count, err := o.QueryTable("evaluation").Filter("session_id", sessionId).Filter("user_id", userId).Count()
	if err != nil {
		return false
	}
	if count > 0 {
		return true
	}
	return false
}
