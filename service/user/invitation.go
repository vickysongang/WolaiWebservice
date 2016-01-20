package user

import (
	"errors"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"

	"WolaiWebservice/models"
	"WolaiWebservice/service/trade"
)

func CheckUserInvitation(userId int64) (bool, error) {
	var err error

	o := orm.NewOrm()

	user, err := models.ReadUser(userId)
	if err != nil {
		return false, err
	}

	var record models.RegisterInvitation
	err = o.QueryTable(new(models.RegisterInvitation).TableName()).
		Filter("phone", user.Phone).
		One(&record)
	if err != nil {
		return false, errors.New("没有邀请记录")
	}

	if record.ProcessFlag != models.REGISTER_INVITATION_FLAG_NO {
		return false, nil
	}

	_, err = o.QueryTable(new(models.RegisterInvitation).TableName()).
		Filter("phone", user.Phone).
		Update(orm.Params{
		"process_flag": models.REGISTER_INVITATION_FLAG_YES,
	})
	if err != nil {
		seelog.Errorf("%s | Phone: %s", err.Error(), user.Phone)
		return false, errors.New("没有邀请记录")
	}

	trade.HandleTradeRewardInvitation(record.Inviter, record.Amount)

	return true, nil
}
