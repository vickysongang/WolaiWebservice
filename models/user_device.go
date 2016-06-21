package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"
)

type UserDevice struct {
	UserId         int64     `json:"userId" orm:"pk"`
	VersionCode    int64     `json:"versionCode"`
	DeviceType     string    `json:"deviceType"`
	ObjectId       string    `json:"objectId"`
	DeviceToken    string    `json:"deviceToken"`
	DeviceProfile  string    `json:"deviceProfile"`
	InstallationId string    `json:"installationId"`
	TimeZone       string    `json:"timeZone"`
	VoipToken      string    `json:"voipToken"`
	CreateTime     time.Time `json:"createTime" orm:"auto_now_add;type(datetime)"`
	LastUpdateTime time.Time `json:"lastUpdateTime" orm:"type(datetime)"`
}

const (
	DEVICE_TYPE_ANDROID = "android"
	DEVICE_TYPE_IOS     = "ios"

	DEVICE_PROFILE_APPSTORE = "appstore"
	DEVICE_PROFILE_INHOUSE  = "inhouse"
	DEVICE_PROFILE_VOIP     = "voip"
)

func init() {
	orm.RegisterModel(new(UserDevice))
}

func (d *UserDevice) TableName() string {
	return "user_device"
}

func CreateUserDevice(userDevice *UserDevice) (*UserDevice, error) {
	var err error

	o := orm.NewOrm()

	_, err = o.Insert(userDevice)
	if err != nil {
		seelog.Errorf("%s", err.Error())
		return nil, errors.New("创建用户设备信息失败")
	}

	return userDevice, nil
}

func ReadUserDevice(userId int64) (*UserDevice, error) {
	var err error

	o := orm.NewOrm()

	userDevice := UserDevice{UserId: userId}
	err = o.Read(&userDevice)
	if err != nil {
		seelog.Errorf("%s | UserId: %d", err.Error(), userId)
		return nil, errors.New("用户设备信息不存在")
	}

	return &userDevice, nil
}

func UpdateUserDevice(userDevice *UserDevice) (*UserDevice, error) {
	var err error

	o := orm.NewOrm()
	userDevice.LastUpdateTime = time.Now()
	_, err = o.Update(userDevice)
	if err != nil {
		seelog.Errorf("%s | UserId: %d", err.Error(), userDevice.UserId)
		return nil, errors.New("更新用户设备信息失败")
	}

	return userDevice, nil
}
