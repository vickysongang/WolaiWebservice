// POIEvaluationController.go
package controllers

import (
	"math/rand"
	"time"

	"github.com/tmhenry/POIWolaiWebService/models"
)

const (
	PERSONAL_EVALUATION_LABEL = "personal"
	STYLE_EVALUATION_LABEL    = "style"
	SUBJECT_EVALUATION_LABEL  = "subject"
	ABILITY_EVALUATION_LABEL  = "ability"

	TEACHER_EVALUATION_LABEL = "teacher"
	STUDENT_EVALUATION_LABEL = "student"
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

func QuerySystemEvaluationLabels(userId, sessionId, count int64) (models.POIEvaluationLabels, error) {
	labels := models.POIEvaluationLabels{}
	session := models.QuerySessionById(sessionId)
	//如果当前用户是学生，则要返回老师的标签信息，如果当前用户是老师，则要返回学生的标签信息
	//学生
	if userId == session.Created {
		teacher := models.QueryUserById(session.Tutor)
		//个人标签
		teacherPersonalLabels, err := models.QueryEvaluationLabels(teacher.Gender, PERSONAL_EVALUATION_LABEL, TEACHER_EVALUATION_LABEL)
		if err != nil {
			return nil, err
		}
		personalCount := (count - 2) / 2
		for _, v := range GetRandNumSlice(personalCount, int64(len(teacherPersonalLabels))) {
			labels = append(labels, teacherPersonalLabels[v])
		}
		//讲课风格
		teacherStyleLabels, err := models.QueryEvaluationLabels(teacher.Gender, STYLE_EVALUATION_LABEL, TEACHER_EVALUATION_LABEL)
		if err != nil {
			return nil, err
		}
		for _, v := range GetRandNumSlice(count-personalCount-2, int64(len(teacherStyleLabels))) {
			labels = append(labels, teacherStyleLabels[v])
		}
		//科目标签
		order := models.QueryOrderById(session.OrderId)
		teacherSubjectLabels, err := models.QueryEvaluationLabelsBySubject(order.SubjectId)
		if err != nil {
			return nil, err
		}
		for _, v := range GetRandNumSlice(2, int64(len(teacherSubjectLabels))) {
			labels = append(labels, teacherSubjectLabels[v])
		}
	} else if userId == session.Tutor { //老师
		student := models.QueryUserById(session.Created)
		//个人标签
		studentPersonalLabels, err := models.QueryEvaluationLabels(student.Gender, PERSONAL_EVALUATION_LABEL, STUDENT_EVALUATION_LABEL)
		if err != nil {
			return nil, err
		}
		personalCount := count / 2
		for _, v := range GetRandNumSlice(personalCount, int64(len(studentPersonalLabels))) {
			labels = append(labels, studentPersonalLabels[v])
		}
		//能力程度
		studentAbilityLabels, err := models.QueryEvaluationLabels(student.Gender, ABILITY_EVALUATION_LABEL, STUDENT_EVALUATION_LABEL)
		if err != nil {
			return nil, err
		}
		for _, v := range GetRandNumSlice(count-personalCount, int64(len(studentAbilityLabels))) {
			labels = append(labels, studentAbilityLabels[v])
		}
	}
	return labels, nil
}
