// evaluation_apply
package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

const (
	EVALUATION_APPLY_STATUS_IDLE     = "idle"
	EVALUATION_APPLY_STATUS_CREATED  = "created"
	EVALUATION_APPLY_STATUS_APPROVED = "approved"
)

type EvaluationApply struct {
	Id          int64     `json:"id" orm:"pk"`
	UserId      int64     `json:"userId"`
	SessionId   int64     `json:"sessionId"`
	CourseId    int64     `json:"courseId"`
	ChapterId   int64     `json:"chapterId"`
	Status      string    `json:"status"`
	Content     string    `json:"content" orm:"type(longtext)"`
	CreateTime  time.Time `json:"createTime" orm:"auto_now_add;type(datetime)"`
	ApproveTime time.Time `json:"approveTime" orm:"type(datetime)"`
	Approver    string    `json:"approver"`
	Comment     string    `json:"comment"`
	RecordId    int64     `json:"recordId"`
}

func (apply *EvaluationApply) TableName() string {
	return "evaluation_apply"
}

func init() {
	orm.RegisterModel(new(EvaluationApply))
}

func InsertEvaluationApply(apply *EvaluationApply) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(apply)
	return id, err
}

func ReadEvaluationApply(applyId int64) (*EvaluationApply, error) {
	o := orm.NewOrm()
	apply := EvaluationApply{Id: applyId}
	err := o.Read(&apply)
	if err != nil {
		return nil, err
	}
	return &apply, nil
}

func UpdateEvaluationApply(applyId int64, applyInfo map[string]interface{}) error {
	o := orm.NewOrm()
	var params orm.Params = make(orm.Params)
	for k, v := range applyInfo {
		params[k] = v
	}
	_, err := o.QueryTable("evaluation_apply").Filter("id", applyId).Update(params)
	return err
}
