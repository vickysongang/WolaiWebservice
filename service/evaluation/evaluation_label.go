// evaluation_label
package evaluation

import (
	"WolaiWebservice/config"

	"WolaiWebservice/models"

	"github.com/astaxie/beego/orm"
)

func QueryEvaluationLabels(genderType int64, attributeType, objectType string) ([]*models.EvaluationLabel, error) {
	var labels []*models.EvaluationLabel

	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(config.Env.Database.Type)
	qb.Select("name").From("evaluation_label").
		Where("gender_type in (?,2) and attribute_type = ? and object_type in (?,'both')")
	sql := qb.String()
	_, err := o.Raw(sql, genderType, attributeType, objectType).QueryRows(&labels)
	if err != nil {
		return nil, err
	}
	return labels, nil
}

func QueryEvaluationLabelsBySubject(subjectId int64) ([]*models.EvaluationLabel, error) {
	var labels []*models.EvaluationLabel

	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(config.Env.Database.Type)
	qb.Select("evaluation_label.name").From("evaluation_label").
		InnerJoin("evaluation_to_subject").On("evaluation_label.id = evaluation_to_subject.label_id").
		Where("evaluation_label.attribute_type = 'subject' and evaluation_to_subject.subject_id = ?")
	sql := qb.String()
	_, err := o.Raw(sql, subjectId).QueryRows(&labels)
	if err != nil {
		return nil, err
	}
	return labels, nil
}
