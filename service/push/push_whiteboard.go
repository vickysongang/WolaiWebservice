package push

import (
	"WolaiWebservice/models"
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

	if userDevice.DeviceType == models.DEVICE_TYPE_IOS {
		apnsprovider.PushWhiteboardCall(userDevice.DeviceToken, userDevice.DeviceProfile, callerId)
	} else if userDevice.DeviceType == models.DEVICE_TYPE_ANDROID {
		lcpush.PushWhiteboardCall(userDevice.ObjectId, callerId)
	}

	return nil
}
