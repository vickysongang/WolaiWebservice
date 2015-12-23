package order

import (
	"errors"

	"WolaiWebservice/models"
)

func GetOrderInfo(orderId int64) (int64, error, *models.Order) {
	order, err := models.ReadOrder(orderId)
	if err != nil {
		return 2, errors.New("订单资料异常"), nil
	}

	return 0, nil, order
}
