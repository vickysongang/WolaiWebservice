package push

import (
	"WolaiWebservice/models"
	"WolaiWebservice/redis"
	"WolaiWebservice/utils/apnsprovider"
	"WolaiWebservice/utils/leancloud/lcpush"
)

func PushWhiteboardCall(userId, callerId int64) error {
	var err error

	_, err = models.ReadUser(userId)
	if err != nil {
		return err
	}

	_, err = models.ReadUser(callerId)
	if err != nil {
		return err
	}

	userDevice, err := models.ReadUserDevice(userId)
	if err != nil {
		return err
	}
	if redis.HasUserObjectId(userId) {
		if userDevice.DeviceType == models.DEVICE_TYPE_IOS {
			go apnsprovider.PushWhiteboardCall(userDevice.DeviceToken, userDevice.DeviceProfile, callerId)
		} else if userDevice.DeviceType == models.DEVICE_TYPE_ANDROID {
			go lcpush.PushWhiteboardCall(userDevice.ObjectId, callerId)
		}
	}
	return nil
}
