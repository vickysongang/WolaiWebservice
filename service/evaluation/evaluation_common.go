// evaluation_common
package evaluation

import (
	"WolaiWebservice/models"

	"github.com/astaxie/beego/orm"
)

func GetEvaluationApply(userId, sessionId, courseId int64) (*models.EvaluationApply, error) {
	o := orm.NewOrm()
	var apply models.EvaluationApply
	err := o.QueryTable(new(models.EvaluationApply).TableName()).
		Filter("user_id", userId).
		Filter("session_id", sessionId).
		Filter("course_id", courseId).One(&apply)
	return &apply, err
}
