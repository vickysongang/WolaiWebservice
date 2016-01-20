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

type VersionRequire struct {
	MinAndroidVersion int64
	MinIOSVersion     int64
}

func CheckUserVersion(userId int64, req *VersionRequire) bool {
	var err error

	userDevice, err := models.ReadUserDevice(userId)
	if err != nil {
		return false
	}

	if userDevice.DeviceType == models.DEVICE_TYPE_ANDROID {
		if userDevice.VersionCode >= req.MinAndroidVersion {
			return true
		} else {
			return false
		}
	}

	if userDevice.DeviceType == models.DEVICE_TYPE_IOS {
		if userDevice.VersionCode >= req.MinIOSVersion {
			return true
		} else {
			return false
		}
	}

	return false
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

func UpdateUserInfo(userId, gender int64, nickname, avatar string) (*models.User, error) {
	var err error

	user, err := models.ReadUser(userId)
	if err != nil {
		return nil, err
	}

	user.Nickname = nickname
	user.Avatar = avatar
	user.Gender = gender

	user, err = models.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
