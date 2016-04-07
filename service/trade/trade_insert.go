package trade

import (
	"math"
	//"errors"
	//"time"

	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
	"WolaiWebservice/utils/leancloud/lcmessage"

	seelog "github.com/cihub/seelog"
)

func createTradeRecord(userId, amount int64, tradeType, result, comment string,
	sessionId, recordId, pingppId int64, chargeCode string, qapkgTimeLength int64, chapterId int64) (*models.TradeRecord, error) {
	var err error

	//	err = HandleUserBalance(userId, amount)
	//	if err != nil {
	//		return nil, err
	//	}

	user, err := models.ReadUser(userId)
	if err != nil {
		return nil, err
	}

	record := models.TradeRecord{
		UserId:          userId,
		TradeType:       tradeType,
		TradeAmount:     amount,
		Result:          result,
		Balance:         user.Balance,
		Comment:         comment,
		SessionId:       sessionId,
		RecordId:        recordId,
		PingppId:        pingppId,
		ChargeCode:      chargeCode,
		QapkgTimeLength: qapkgTimeLength,
		ChapterId:       chapterId,
	}

	tradeRecord, err := models.InsertTradeRecord(&record)

	if err != nil {
		return nil, err
	}

	go lcmessage.SendTradeNotification(tradeRecord.Id)

	return tradeRecord, nil
}

/*
 * 操作用户的余额
 */
func HandleUserBalance(userId int64, amount int64) error {
	if amount < 0 {
		user, _ := models.ReadUser(userId)
		if int64(math.Abs(float64(amount))) > user.Balance {
			amount = 0 - user.Balance
		}
	}
	o := orm.NewOrm()
	_, err := o.QueryTable("users").Filter("id", userId).Update(orm.Params{
		"balance": orm.ColValue(orm.ColAdd, amount),
	})
	seelog.Debug("HandleUserBalance | userId:", userId, " amount:", amount)
	return err
}
