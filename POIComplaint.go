// POIComplains.go
package main

import (
	"time"

	"github.com/astaxie/beego/orm"
	seelog "github.com/cihub/seelog"
)

type POIComplaint struct {
	Id         int64     `json:"id" orm:"pk"`
	UserId     int64     `json:"userId"`
	SessionId  int64     `json:"sessionId"`
	CreateTime time.Time `json:"-" orm:"auto_now_add;type(datetime)"`
	Reasons    string    `json:"reasons"`
	Comment    string    `json:"comment"`
	Status     string    `json:"status"`
	Suggestion string    `json:"suggestion"`
}

func (c *POIComplaint) TableName() string {
	return "complaint"
}

func init() {
	orm.RegisterModel(new(POIComplaint))
}

func InsertPOIComplaint(complaint *POIComplaint) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(complaint)
	if err != nil {
		seelog.Error("complaint:", complaint, " ", err.Error())
		return 0, err
	}
	return id, nil
}

func UpdateComplaintInfo(complaintId int64, complaintInfo map[string]interface{}) error {
	o := orm.NewOrm()
	var params orm.Params = make(orm.Params)
	for k, v := range complaintInfo {
		params[k] = v
	}
	_, err := o.QueryTable("complaint").Filter("id", complaintId).Update(params)
	if err != nil {
		seelog.Error("complaintId:", complaintId, " complaintInfo:", complaintInfo, " ", err.Error())
		return err
	}
	return nil
}
