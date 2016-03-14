package user

import (
	"WolaiWebservice/models"
	"WolaiWebservice/utils/leancloud"
	"errors"
	"strings"
)

func SaveLoginInfo(userId int64, objectId, address, ip, userAgent string) (*models.UserLoginInfo, error) {
	var err error

	_, err = models.ReadUser(userId)
	if err != nil {
		return nil, err
	}

	info := models.UserLoginInfo{
		UserId:    userId,
		ObjectId:  objectId,
		Address:   address,
		IP:        ip,
		UserAgent: userAgent,
	}

	loginInfo, err := models.CreateUserLoginInfo(&info)
	if err != nil {
		return nil, err
	}

	return loginInfo, nil
}

func GetDeviceTypeFromUserAgent(userAgent string) (string, error) {
	if strings.HasPrefix(userAgent, "WoLaiAndroid") || strings.Contains(userAgent, "Android") {
		return "android", nil
	} else if strings.HasPrefix(userAgent, "wolaiSocial") || strings.Contains(userAgent, "iOS") {
		return "ios", nil
	} else {
		return "other", errors.New("User Agent does not contain Android/iOS type")
	}
}

func SaveDeviceInfo(userId, versionCode int64, objectId string, userAgent string, voipToken string) (*models.UserDevice, error) {
	var err error

	inst, err := leancloud.LCGetIntallation(objectId)
	if err != nil {
		return nil, err
	}

	if inst.DeviceType == "" {
		deviceType, err := GetDeviceTypeFromUserAgent(userAgent)
		if err != nil {
			return nil, err
		}
		inst.DeviceType = deviceType
	}

	userDevice, err := models.ReadUserDevice(userId)
	if err != nil {
		info := models.UserDevice{
			UserId:         userId,
			VersionCode:    versionCode,
			DeviceType:     inst.DeviceType,
			ObjectId:       inst.ObjectId,
			DeviceToken:    inst.DeviceToken,
			DeviceProfile:  inst.DeviceProfile,
			InstallationId: inst.InstallationId,
			TimeZone:       inst.TimeZone,
			VoipToken:      voipToken,
		}

		userDevice, err = models.CreateUserDevice(&info)
		if err != nil {
			return nil, err
		}
	} else {
		userDevice.VersionCode = versionCode
		userDevice.DeviceType = inst.DeviceType
		userDevice.ObjectId = inst.ObjectId
		userDevice.DeviceToken = inst.DeviceToken
		userDevice.DeviceProfile = inst.DeviceProfile
		userDevice.InstallationId = inst.InstallationId
		userDevice.TimeZone = inst.TimeZone

		userDevice, err = models.UpdateUserDevice(userDevice)
		if err != nil {
			return nil, err
		}
	}

	return userDevice, nil
}
