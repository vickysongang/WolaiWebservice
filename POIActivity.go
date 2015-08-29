// POIActivity.go
package main

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type POIActivity struct {
	Id       int64     `json:"id" orm:"pk"`
	Title    string    `json:"title"`
	Subtitle string    `json:"subtitle"`
	Amount   int64     `json:"amount"`
	TimeFrom time.Time `json:"timeFrom" orm:"type(datetime)"`
	TimeTo   time.Time `json:"timeTo" orm:"type(datetime)"`
	Extra    string    `json:"extra"`
	MediaId  string    `json:"mediaId"`
	Status   string    `json:"status"`
}

type POIUserToActivity struct {
	Id         int64 `json:"-" orm:"pk"`
	UserId     int64 `json:"userId"`
	ActivityId int64 `json:"activityId"`
}

type POIActivities []POIActivity

func (activity *POIActivity) TableName() string {
	return "activities"
}

func init() {
	orm.RegisterModel(new(POIActivity))
}

func InsertActivity(activity *POIActivity) (*POIActivity, error) {
	o := orm.NewOrm()
	id, err := o.Insert(&activity)
	if err != nil {
		return nil, err
	}
	activity.Id = id
	return activity, nil
}
