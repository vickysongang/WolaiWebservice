package push

import (
	"WolaiWebservice/models"
	"WolaiWebservice/utils/apnsprovider"
	"WolaiWebservice/utils/leancloud/lcpush"
)

func PushSessionInstantStart(userId, sessionId int64) error {
	var err error

	_, err = models.ReadUser(userId)
	if err != nil {
		return err
	}

	userDevice, err := models.ReadUserDevice(userId)
	if err != nil {
		return err
	}

	if userDevice.DeviceType == models.DEVICE_TYPE_IOS {
		apnsprovider.PushSessionInstantStart(userDevice.DeviceToken, sessionId)
	} else if userDevice.DeviceType == models.DEVICE_TYPE_ANDROID {
		lcpush.PushSessionInstantStart(userDevice.ObjectId, sessionId)
	}

	return nil
}

func PushSessionResume(userId, sessionId int64) error {
	var err error

	_, err = models.ReadUser(userId)
	if err != nil {
		return err
	}

	userDevice, err := models.ReadUserDevice(userId)
	if err != nil {
		return err
	}

	if userDevice.DeviceType == models.DEVICE_TYPE_IOS {
		apnsprovider.PushSessionResume(userDevice.DeviceToken, sessionId)
	} else if userDevice.DeviceType == models.DEVICE_TYPE_ANDROID {
		lcpush.PushSessionResume(userDevice.ObjectId, sessionId)
	}

	return nil
}
