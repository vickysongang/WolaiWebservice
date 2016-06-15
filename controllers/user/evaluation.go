// evaluation
package user

import (
	"WolaiWebservice/models"
	userService "WolaiWebservice/service/user"
	"time"
)

type EvaluationListItem struct {
	EvaluationId      int64     `json:"evaluationId"`
	UserId            int64     `json:"userId"`
	Avatar            string    `json:"avatar"`
	Nickname          string    `json:"nickname"`
	Phone             string    `json:"phone"`
	EvaluationContent string    `json:"evaluationContent"`
	PubTime           time.Time `json:"pubTime"`
}

func AssembleTeacherEvaluationList(teacherId, page, count int64) ([]*EvaluationListItem, error) {
	var err error

	result := make([]*EvaluationListItem, 0)

	evaluations, err := userService.GetTeacherEvaluations(teacherId, page, count)
	if err != nil {
		return result, nil
	}

	for _, evaluation := range evaluations {
		student, err := models.ReadUser(evaluation.UserId)
		if err != nil {
			continue
		}
		item := EvaluationListItem{
			EvaluationId:      evaluation.EvaluationId,
			UserId:            student.Id,
			Avatar:            student.Avatar,
			Nickname:          student.Nickname,
			Phone:             *student.Phone,
			EvaluationContent: evaluation.Content,
			PubTime:           evaluation.PubTime,
		}
		result = append(result, &item)
	}

	return result, nil
}
