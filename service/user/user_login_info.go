package user

import (
	"WolaiWebservice/models"
	"WolaiWebservice/utils/leancloud"
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

func SaveDeviceInfo(userId, versionCode int64, objectId string) (*models.UserDevice, error) {
	var err error

	inst, err := leancloud.LCGetIntallation(objectId)
	if err != nil {
		return nil, err
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
