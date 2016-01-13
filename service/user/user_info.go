package user

import (
	"time"

	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
)

func IsTeacherFirstLogin(user *models.User) (bool, error) {
	var err error

	if user.AccessRight != models.USER_ACCESSRIGHT_TEACHER {
		return false, nil
	}

	if user.Status != models.USER_STATUS_INACTIVE {
		return false, nil
	}

	user.Status = models.USER_STATUS_ACTIVE
	user, err = models.UpdateUser(user)
	if err != nil {
		return false, err
	}

	return true, nil
}

func IsUserExisting(userId int64) bool {
	o := orm.NewOrm()

	exist := o.QueryTable(new(models.User).TableName()).
		Filter("id", userId).Exist()

	return exist
}

func UpdateUserLastLoginTime(user *models.User) error {
	var err error

	user.LastLoginTime = time.Now()
	user, err = models.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}
