// POIExperience
package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type POIExperience struct {
	Id         int64     `json:"-" orm:"pk"`
	Nickname   string    `json:"nickname"`
	Phone      string    `json:"phone"`
	CreateTime time.Time `json:"-" orm:"auto_now_add;type(datetime)"`
}

func (experience *POIExperience) TableName() string {
	return "experience"
}

func init() {
	orm.RegisterModel(new(POIExperience))
}

func InsertExperience(experience *POIExperience) (*POIExperience, error) {
	o := orm.NewOrm()
	id, err := o.Insert(experience)
	if err != nil {
		return nil, err
	}
	experience.Id = id
	return experience, nil
}

func CheckExperienceExsits(phone string) bool {
	o := orm.NewOrm()
	return o.QueryTable("experience").Filter("phone", phone).Exist()
}
