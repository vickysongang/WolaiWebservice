package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Complaint struct {
	Id         int64     `json:"id" orm:"pk"`
	UserId     int64     `json:"userId"`
	SessionId  int64     `json:"sessionId"`
	CreateTime time.Time `json:"-" orm:"auto_now_add;type(datetime)"`
	Reasons    string    `json:"reasons"`
	Comment    string    `json:"comment"`
	Status     string    `json:"status"`
	Suggestion string    `json:"suggestion"`
}

func (c *Complaint) TableName() string {
	return "complaint"
}

func init() {
	orm.RegisterModel(new(Complaint))
}

func InsertComplaint(complaint *Complaint) (*Complaint, error) {
	var err error

	o := orm.NewOrm()

	id, err := o.Insert(complaint)
	if err != nil {
		return nil, err
	}
	complaint.Id = id

	return complaint, nil
}

func UpdateComplaint(complaintId int64, complaintInfo map[string]interface{}) error {
	var err error

	o := orm.NewOrm()

	var params orm.Params = make(orm.Params)
	for k, v := range complaintInfo {
		params[k] = v
	}

	_, err = o.QueryTable("complaint").Filter("id", complaintId).Update(params)

	return err
}
