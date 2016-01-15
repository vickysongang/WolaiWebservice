package push

import (
	"WolaiWebservice/models"
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

	if userDevice.DeviceType == models.DEVICE_TYPE_IOS {
		apnsprovider.PushNewOrderDispatch(userDevice.DeviceToken, orderId)
	} else if userDevice.DeviceType == models.DEVICE_TYPE_ANDROID {
		lcpush.PushNewOrderDispatch(userDevice.ObjectId, orderId)
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

	if userDevice.DeviceType == models.DEVICE_TYPE_IOS {
		apnsprovider.PushNewOrderAssign(userDevice.DeviceToken, orderId)
	} else if userDevice.DeviceType == models.DEVICE_TYPE_ANDROID {
		lcpush.PushNewOrderAssign(userDevice.ObjectId, orderId)
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

	if userDevice.DeviceType == models.DEVICE_TYPE_IOS {
		apnsprovider.PushOrderAccept(userDevice.DeviceToken, orderId, teacherId)
	} else if userDevice.DeviceType == models.DEVICE_TYPE_ANDROID {
		lcpush.PushOrderAccept(userDevice.ObjectId, orderId, teacherId)
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

	if userDevice.DeviceType == models.DEVICE_TYPE_IOS {
		apnsprovider.PushOrderPersonalAccept(userDevice.DeviceToken, orderId, teacherId)
	} else if userDevice.DeviceType == models.DEVICE_TYPE_ANDROID {
		lcpush.PushOrderAccept(userDevice.ObjectId, orderId, teacherId)
	}

	return nil
}
