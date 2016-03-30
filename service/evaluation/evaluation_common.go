// evaluation_common
package evaluation

import (
	"WolaiWebservice/models"

	"github.com/astaxie/beego/orm"
)

func GetEvaluationApply(userId, chapterId int64) (*models.EvaluationApply, error) {
	o := orm.NewOrm()
	var apply models.EvaluationApply
	err := o.QueryTable(new(models.EvaluationApply).TableName()).
		Filter("user_id", userId).
		Filter("chapter_id", chapterId).One(&apply)
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
