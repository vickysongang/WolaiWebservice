package auth

import (
	"errors"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"

	"WolaiWebservice/models"
)

func OauthBind(userId int64, openId string) (*models.UserOauth, error) {
	var err error

	newUserOauth := models.UserOauth{
		UserId:   userId,
		OpenIdQQ: openId,
	}

	userOauth, err := models.CreateUserOauth(&newUserOauth)
	if err != nil {
		return nil, err
	}

	return userOauth, nil
}

func HasOauthBound(userId int64) (bool, error) {
	var err error

	_, err = models.ReadUserOauth(userId)
	if err != nil {
		return false, nil
	}

	return true, errors.New("用户已与其他账号绑定")
}

func QueryUserOauthByOpenId(openId string) (*models.UserOauth, error) {
	var err error

	o := orm.NewOrm()

	var userOauth models.UserOauth
	err = o.QueryTable("user_oauth").Filter("open_id_qq", openId).One(&userOauth)
	if err != nil {
		seelog.Error("%s | OpenIdQQ: %s", err.Error(), openId)
		return nil, errors.New("未找到该账号的绑定信息")
	}

	return &userOauth, nil
}
