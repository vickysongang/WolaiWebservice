package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"
)

type UserDataUsage struct {
	UserId           int64     `json:"userId" orm:"column(user_id);pk"`
	Data             int64     `json:"data" orm:"column(total_data)"`
	DataClass        int64     `json:"dataClass" orm:"column(total_data_class)"`
	LastUpdateTime   time.Time `json:"lastUpdateTime" orm:"column(last_update_time);type(datetime);auto_now"`
	DataToClaim      int64     `json:"dataToClaim" orm:"column(total_data_to_claim)"`
	DataClassToClaim int64     `json:"dataClassToClaim" orm:"column(total_data_class_to_claim)"`
}

func init() {
	orm.RegisterModel(new(UserDataUsage))
}

func (tp *UserDataUsage) TableName() string {
	return "user_data_usage"
}

func CreateUserDataUsage(userDataUsage *UserDataUsage) (*UserDataUsage, error) {
	var err error

	o := orm.NewOrm()

	_, err = o.Insert(userDataUsage)
	if err != nil {
		seelog.Error("%s", err.Error())
		return nil, errors.New("创建用户流量失败")
	}
	return userDataUsage, nil
}

func ReadUserDataUsage(userId int64) (*UserDataUsage, error) {
	var err error

	o := orm.NewOrm()

	userDataUsage := UserDataUsage{UserId: userId}
	err = o.Read(&userDataUsage)
	if err != nil {
		seelog.Errorf("%s | UserId: %d", err.Error(), userId)
		return nil, errors.New("未找到用户流量信息")
	}

	return &userDataUsage, nil
}

func UpdateUserDataUsage(userDataUsage *UserDataUsage) (*UserDataUsage, error) {
	var err error

	o := orm.NewOrm()

	_, err = o.Update(userDataUsage)
	if err != nil {
		seelog.Errorf("%s | UserId: %d", err.Error(), userDataUsage.UserId)
		return nil, errors.New("更新用户失败")
	}

	return userDataUsage, nil
}

func HandleDataClaimUpdate(userId int64, data, dataClass, totalClaimAdd, totalClassClaimAdd int64) error {
	o := orm.NewOrm()
	_, err := o.QueryTable("user_data_usage").Filter("user_id", userId).Update(orm.Params{
		"total_data_to_claim":       orm.ColValue(orm.ColAdd, totalClaimAdd),
		"total_data_class_to_claim": orm.ColValue(orm.ColAdd, totalClassClaimAdd),
		"total_data":                data,
		"total_data_class":          dataClass,
		"last_update_time":          time.Now(),
	})
	return err
}
