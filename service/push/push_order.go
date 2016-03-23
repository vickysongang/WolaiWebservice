package push

import (
	"WolaiWebservice/models"
	"WolaiWebservice/redis"
	"WolaiWebservice/utils/apnsprovider"
	"WolaiWebservice/utils/leancloud/lcpush"
)

func PushNewOrderDispatch(userId, orderId int64) error {
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
			go apnsprovider.PushNewOrderDispatch(userDevice.DeviceToken, userDevice.DeviceProfile, orderId)
		} else if userDevice.DeviceType == models.DEVICE_TYPE_ANDROID {
			go lcpush.PushNewOrderDispatch(userDevice.ObjectId, orderId)
		}
	}
	return nil
}

func PushNewOrderAssign(userId, orderId int64) error {
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
			go apnsprovider.PushNewOrderAssign(userDevice.DeviceToken, userDevice.DeviceProfile, orderId)
		} else if userDevice.DeviceType == models.DEVICE_TYPE_ANDROID {
			go lcpush.PushNewOrderAssign(userDevice.ObjectId, orderId)
		}
	}
	return nil
}

func PushOrderAccept(userId, orderId, teacherId int64) error {
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
			go apnsprovider.PushOrderAccept(userDevice.DeviceToken, userDevice.DeviceProfile, orderId, teacherId)
		} else if userDevice.DeviceType == models.DEVICE_TYPE_ANDROID {
			go lcpush.PushOrderAccept(userDevice.ObjectId, orderId, teacherId)
		}
	}
	return nil
}

func PushOrderPersonalAccept(userId, orderId, teacherId int64) error {
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
			go apnsprovider.PushOrderPersonalAccept(userDevice.DeviceToken, userDevice.DeviceProfile, orderId, teacherId)
		} else if userDevice.DeviceType == models.DEVICE_TYPE_ANDROID {
			go lcpush.PushOrderAccept(userDevice.ObjectId, orderId, teacherId)
		}
	}
	return nil
}
