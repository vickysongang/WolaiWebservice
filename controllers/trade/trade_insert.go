package trade

import (
	//"errors"
	//"time"

	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
)

func createTradeRecord(userId, amount int64, tradeType, result, comment string,
	sessionId, recordId, pingppId int64) (*models.TradeRecord, error) {
	var err error

	err = addUserBalance(userId, amount)
	if err != nil {
		return nil, err
	}

	user, err := models.ReadUser(userId)
	if err != nil {
		return nil, err
	}

	record := models.TradeRecord{
		UserId:      userId,
		TradeType:   tradeType,
		TradeAmount: amount,
		Result:      result,
		Balance:     user.Balance,
		Comment:     comment,
		SessionId:   sessionId,
		RecordId:    recordId,
		PingppId:    pingppId,
	}

	tradeRecord, err := models.InsertTradeRecord(&record)

	if err != nil {
		return nil, err
	}

	return tradeRecord, nil
}

/*
* 增加用户的余额
 */
func addUserBalance(userId int64, amount int64) error {
	o := orm.NewOrm()

	_, err := o.QueryTable("users").Filter("id", userId).Update(orm.Params{
		"balance": orm.ColValue(orm.Col_Add, amount),
	})

	return err
}

/*
* 减少用户的余额
 */
func minusUserBalance(userId int64, amount int64) error {
	o := orm.NewOrm()

	_, err := o.QueryTable("users").Filter("id", userId).Update(orm.Params{
		"balance": orm.ColValue(orm.Col_Minus, amount),
	})

	return err
}
