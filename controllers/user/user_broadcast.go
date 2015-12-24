package user

import (
	"errors"

	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
)

func GetUserBroadcast(userId int64) (int64, error, []*models.Broadcast) {
	o := orm.NewOrm()

	user, err := models.ReadUser(userId)
	if err != nil {
		return 2, errors.New("用户信息异常"), nil
	}

	var broadcasts []*models.Broadcast
	_, err = o.QueryTable("broadcast").Filter("access_right", user.AccessRight).Filter("active", "Y").
		OrderBy("rank").All(&broadcasts)

	if err != nil {
		return 2, errors.New("数据异常"), nil
	}

	return 0, nil, broadcasts
}
