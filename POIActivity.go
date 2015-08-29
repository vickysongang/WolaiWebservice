// POIActivity.go
package main

import (
	"time"

	"github.com/astaxie/beego/orm"
)

const (
	REGISTER_ACTIVITY = "REGISTER"
)

type POIActivity struct {
	Id       int64     `json:"id" orm:"pk"`
	Theme    string    `json:"theme"`
	Title    string    `json:"title"`
	Subtitle string    `json:"subtitle"`
	Amount   int64     `json:"amount"`
	TimeFrom time.Time `json:"timeFrom" orm:"type(datetime)"`
	TimeTo   time.Time `json:"timeTo" orm:"type(datetime)"`
	Extra    string    `json:"extra"`
	MediaId  string    `json:"mediaId"`
	Status   string    `json:"status"`
	Type     string    `json:"activity"`
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

func InsertUserToActivity(userToActivity *POIUserToActivity) (*POIUserToActivity, error) {
	o := orm.NewOrm()
	id, err := o.Insert(&userToActivity)
	if err != nil {
		return nil, err
	}
	userToActivity.Id = id
	return userToActivity, nil
}

func QueryEffectiveActivities(activityType string) (POIActivities, error) {
	activities := make(POIActivities, 0)
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(DB_TYPE)
	qb.Select("id,theme,title,subtitle,amount,time_from,time_to,extra,media_id,status,type").From("activities").
		Where("type = ? and now() BETWEEN time_from and time_to and status = 'open'")
	sql := qb.String()
	_, err := o.Raw(sql, activityType).QueryRows(&activities)
	if err != nil {
		return nil, err
	}
	return activities, nil
}
