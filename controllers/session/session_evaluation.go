package session

import (
	"math/rand"
	"time"

	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
)

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
		teacher := models.QueryUserById(session.Tutor)
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
		student := models.QueryUserById(session.Creator)
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

type evaluationInfo struct {
	Type      string             `json:"type"`
	Evalution *models.Evaluation `json:"evaluationInfo"`
}

func QueryEvaluationInfo(userId, sessionId int64) ([]*evaluationInfo, error) {
	session, _ := models.ReadSession(sessionId)
	self, err1 := models.QueryEvaluation4Self(userId, sessionId)
	other, err2 := models.QueryEvaluation4Other(userId, sessionId)

	selfEvaluation := evaluationInfo{}
	otherEvaluation := evaluationInfo{}

	evalutionInfos := make([]*evaluationInfo, 0)
	if userId == session.Tutor {
		if err1 == nil {
			selfEvaluation.Type = "teacher"
			selfEvaluation.Evalution = self

			evalutionInfos = append(evalutionInfos, &selfEvaluation)
		}
		if err2 == nil {
			otherEvaluation.Type = "student"
			otherEvaluation.Evalution = other

			evalutionInfos = append(evalutionInfos, &otherEvaluation)
		}
	} else if userId == session.Creator {
		if err1 == nil {
			selfEvaluation.Type = "student"
			selfEvaluation.Evalution = self

			evalutionInfos = append(evalutionInfos, &selfEvaluation)
		}
		if err2 == nil {
			otherEvaluation.Type = "teacher"
			otherEvaluation.Evalution = other

			evalutionInfos = append(evalutionInfos, &otherEvaluation)
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
