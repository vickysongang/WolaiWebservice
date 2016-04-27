// evaluation_common
package evaluation

import (
	"WolaiWebservice/models"

	"github.com/astaxie/beego/orm"
)

func GetEvaluationApply(userId, chapterId, recordId int64) (*models.EvaluationApply, error) {
	o := orm.NewOrm()
	var apply models.EvaluationApply
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
	err := o.QueryTable(new(models.EvaluationApply).TableName()).SetCond(cond).One(&apply)
	return &apply, err
}

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
