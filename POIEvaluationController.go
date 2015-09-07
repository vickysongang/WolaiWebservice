// POIEvaluationController.go
package main

import (
	"math/rand"
	"time"
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
	if sliceSize <= 0 {
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

func QuerySystemEvaluationLabels(userId, sessionId, count int64) (POIEvaluationLabels, error) {
	labels := POIEvaluationLabels{}
	session := QuerySessionById(sessionId)
	user := QueryUserById(userId)
	//学生
	if userId == session.Created {
		//个人标签
		studentPersonalLabels, err := QueryEvaluationLabels(user.Gender, PERSONAL_EVALUATION_LABEL, STUDENT_EVALUATION_LABEL)
		if err != nil {
			return nil, err
		}
		personalCount := count / 2
		for _, v := range GetRandNumSlice(personalCount, int64(len(studentPersonalLabels))) {
			labels = append(labels, studentPersonalLabels[v])
		}
		//能力程度
		studentAbilityLabels, err := QueryEvaluationLabels(user.Gender, ABILITY_EVALUATION_LABEL, STUDENT_EVALUATION_LABEL)
		if err != nil {
			return nil, err
		}
		for _, v := range GetRandNumSlice(count-personalCount, int64(len(studentAbilityLabels))) {
			labels = append(labels, studentAbilityLabels[v])
		}
	} else if userId == session.Tutor { //老师
		//个人标签
		teacherPersonalLabels, err := QueryEvaluationLabels(user.Gender, PERSONAL_EVALUATION_LABEL, TEACHER_EVALUATION_LABEL)
		if err != nil {
			return nil, err
		}
		personalCount := (count - 2) / 2
		for _, v := range GetRandNumSlice(personalCount, int64(len(teacherPersonalLabels))) {
			labels = append(labels, teacherPersonalLabels[v])
		}
		//讲课风格
		teacherStyleLabels, err := QueryEvaluationLabels(user.Gender, STYLE_EVALUATION_LABEL, TEACHER_EVALUATION_LABEL)
		if err != nil {
			return nil, err
		}
		for _, v := range GetRandNumSlice(count-personalCount-2, int64(len(teacherStyleLabels))) {
			labels = append(labels, teacherStyleLabels[v])
		}
		//科目标签
		order := QueryOrderById(session.OrderId)
		teacherSubjectLabels, err := QueryEvaluationLabelsBySubject(order.SubjectId)
		if err != nil {
			return nil, err
		}
		for _, v := range GetRandNumSlice(2, int64(len(teacherSubjectLabels))) {
			labels = append(labels, teacherSubjectLabels[v])
		}
	}
	return labels, nil
}
