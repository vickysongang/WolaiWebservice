package trade

import (
	"errors"

	"WolaiWebservice/models"
	tradeService "WolaiWebservice/service/trade"
)

func GetChargeBanner(userId int64) (int64, error, []*models.ChargeBanner) {
	_, err := models.ReadUser(userId)
	if err != nil {
		return 2, errors.New("用户信息异常"), nil
	}

	banners, err := tradeService.QueryChargeBanners()
	if err != nil {
		return 2, errors.New("数据异常"), nil
	}

	return 0, nil, banners
}
