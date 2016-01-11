package trade

import (
	"errors"
	"time"

	"WolaiWebservice/models"
)

var (
	ErrChargeCodeInvalid error
	ErrChargeCodeError   error
)

func init() {
	ErrChargeCodeInvalid = errors.New("充值码无效")
	ErrChargeCodeError = errors.New("充值码兑换异常")
}

func ValidateChargeCode(code string) (*models.ChargeCode, error) {
	var err error

	chargeCode, err := models.ReadChargeCode(code)
	if err != nil {
		return nil, ErrChargeCodeInvalid
	}

	if chargeCode.UseFlag == models.CODE_USE_FLAG_YES {
		return nil, ErrChargeCodeInvalid
	}

	return chargeCode, nil
}

func ApplyChargeCode(userId int64, code string) error {
	var err error

	chargeCode, err := models.ReadChargeCode(code)
	if err != nil {
		return ErrChargeCodeInvalid
	}

	if chargeCode.UseFlag == models.CODE_USE_FLAG_YES {
		return ErrChargeCodeInvalid
	}

	chargeCode.UseFlag = models.CODE_USE_FLAG_NO
	chargeCode.UseTime = time.Now()

	chargeCode, err = models.UpdateChargeCode(chargeCode)
	if err != nil {
		return ErrChargeCodeError
	}

	return nil
}
