// POILocation
package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type POILocation struct {
	Id        int64 `orm:"pk"`
	UserId    int64
	ObjectId  string
	Address   string
	Ip        string
	UserAgent string
	Time      time.Time `orm:"auto_now_add;type(datetime)"`
}

func (location *POILocation) TableName() string {
	return "location"
}

func init() {
	orm.RegisterModel(new(POILocation))
}

func InsertLocation(location *POILocation) (*POILocation, error) {
	o := orm.NewOrm()
	id, err := o.Insert(location)
	if err != nil {
		return nil, err
	}
	location.Id = id
	return location, nil
}
