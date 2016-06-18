// evaluation_common
package evaluation

import (
	"WolaiWebservice/models"

	"github.com/astaxie/beego/orm"
)

func GetEvaluationDetailUrlPrefix() string {
	o := orm.NewOrm()
	var dictionary models.Dictionary
	o.QueryTable(new(models.Dictionary).TableName()).
		Filter("code", "detail_url").Filter("type", models.DICT_TYPE_EVALUATION).One(&dictionary)
	if dictionary.Id != 0 {
		return dictionary.Meaning
	}
	return ""
}

func QueryEvaluation(userId, sessionId int64) (*models.Evaluation, error) {
	o := orm.NewOrm()
	evalution := models.Evaluation{}
	err := o.QueryTable("evaluation").
		Filter("user_id", userId).
		Filter("session_id", sessionId).
		One(&evalution)
	return &evalution, err
}

func QueryEvaluationByChapter(userId, chapterId, recordId int64) (*models.Evaluation, error) {
	o := orm.NewOrm()
	evalution := models.Evaluation{}
	cond := orm.NewCondition()
	if userId != 0 {
		cond = cond.And("user_id", userId)
	}
	if chapterId != 0 {
		cond = cond.And("chapter_id", chapterId)
	}
	if recordId != 0 {
		cond = cond.And("record_id", recordId)
	}
	err := o.QueryTable("evaluation").SetCond(cond).One(&evalution)
	return &evalution, err
}
