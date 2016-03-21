package pingxx

import (
	"errors"
	"sync"
	"time"

	"github.com/astaxie/beego/orm"

	courseController "WolaiWebservice/controllers/course"
	qaPkgController "WolaiWebservice/controllers/qapkg"
	"WolaiWebservice/models"
	"WolaiWebservice/service/trade"

	seelog "github.com/cihub/seelog"
)

var ErrChargeNotFound = errors.New("Charge is not found")

type PingxxWebhookManager struct {
	chargeMap map[string]int64
	lock      sync.RWMutex
}

var WebhookManager *PingxxWebhookManager

func init() {
	WebhookManager = NewPingxxWebhookManager()
}

func NewPingxxWebhookManager() *PingxxWebhookManager {
	manager := PingxxWebhookManager{
		chargeMap: make(map[string]int64),
	}
	return &manager
}

func (pwm *PingxxWebhookManager) SetChargeOnline(chargeId string) int64 {
	_, ok := pwm.chargeMap[chargeId]
	if ok {
		return -1
	}
	pwm.lock.Lock()
	defer pwm.lock.Unlock()
	pwm.chargeMap[chargeId] = time.Now().Unix()
	return 0
}

func (pwm *PingxxWebhookManager) IsChargeOnline(chargeId string) bool {
	_, ok := pwm.chargeMap[chargeId]
	seelog.Debug("Pingxx webhook | IsChargeOnline:", chargeId, " ", ok)
	return ok
}

func (pwm *PingxxWebhookManager) ChargeSuccessEvent(chargeId string) {
	if !pwm.IsChargeOnline(chargeId) {
		state := pwm.SetChargeOnline(chargeId)
		if state == -1 {
			return
		}
	} else {
		return
	}

	recordInfo := map[string]interface{}{
		"Result": "success",
	}
	models.UpdatePingppRecord(chargeId, recordInfo)

	record, _ := models.QueryPingppRecordByChargeId(chargeId)
	if record.Id == 0 {
		return
	}

	if pwm.checkChargeSuccessExist(record, record.Type) {
		return
	}

	switch record.Type {
	case models.TRADE_CHARGE:
		premium, _ := trade.GetChargePremuim(record.UserId, int64(record.Amount))

		trade.HandleTradeChargePingpp(record.Id)
		if premium > 0 {
			trade.HandleTradeChargePremium(record.UserId, premium, "", record.Id, "")
		}
	case models.TRADE_COURSE_AUDITION:
		courseController.HandleCourseActionPayByThird(record.UserId, record.RefId, record.Type, int64(record.Amount), record.Id)

	case models.TRADE_COURSE_PURCHASE:
		courseController.HandleCourseActionPayByThird(record.UserId, record.RefId, record.Type, int64(record.Amount), record.Id)

	case models.TRADE_QA_PKG_PURCHASE:
		qaPkgController.HandleQaPkgActionPayByThird(record.UserId, record.RefId, int64(record.Amount), record.Id)
	}
}

func (pwm *PingxxWebhookManager) RefundSuccessEvent(chargeId string, refundId string) {
	recordInfo := map[string]interface{}{
		"Result":   "success",
		"RefundId": refundId,
	}
	models.UpdatePingppRecord(chargeId, recordInfo)
}

func (pwm *PingxxWebhookManager) checkChargeSuccessExist(record *models.PingppRecord, tradeType string) bool {
	o := orm.NewOrm()

	exist := o.QueryTable(new(models.TradeRecord).TableName()).
		Filter("user_id", record.UserId).
		Filter("trade_type", tradeType).
		Filter("pingpp_id", record.Id).Exist()

	return exist
}
