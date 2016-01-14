package push

import (
	"WolaiWebservice/models"
	"WolaiWebservice/utils/apnsprovider"
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
		apnsprovider.PushNewOrderDispatch(orderId, userDevice.DeviceToken)
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
		apnsprovider.PushNewOrderAssign(orderId, userDevice.DeviceToken)
	}

	return nil
}
