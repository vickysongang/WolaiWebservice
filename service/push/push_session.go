package push

import (
	"WolaiWebservice/models"
	"WolaiWebservice/redis"
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
	if redis.HasUserObjectId(userId) {
		if userDevice.DeviceType == models.DEVICE_TYPE_IOS {
			go apnsprovider.PushSessionInstantStart(userDevice.DeviceToken, userDevice.DeviceProfile, sessionId)
		} else if userDevice.DeviceType == models.DEVICE_TYPE_ANDROID {
			go lcpush.PushSessionInstantStart(userDevice.ObjectId, sessionId)
		}
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
	if redis.HasUserObjectId(userId) {
		if userDevice.DeviceType == models.DEVICE_TYPE_IOS {
			go apnsprovider.PushSessionResume(userDevice.DeviceToken, userDevice.DeviceProfile, sessionId)
		} else if userDevice.DeviceType == models.DEVICE_TYPE_ANDROID {
			go lcpush.PushSessionResume(userDevice.ObjectId, sessionId)
		}
	}
	return nil
}
