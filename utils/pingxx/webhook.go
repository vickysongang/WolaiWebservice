package pingxx

import (
	"errors"
	"sync"
	"time"

	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
	"WolaiWebservice/service/trade"

	seelog "github.com/cihub/seelog"
)

type PingxxWebhookManager struct {
	chargeMap map[string]int64
	lock      *sync.RWMutex
}

var ErrChargeNotFound = errors.New("Charge is not found")

var WebhookManager *PingxxWebhookManager

func NewPingxxWebhookManager() *PingxxWebhookManager {
	manager := PingxxWebhookManager{
		chargeMap: make(map[string]int64),
	}
	return &manager
}

func init() {
	WebhookManager = NewPingxxWebhookManager()
}

func (pwm *PingxxWebhookManager) SetChargeOnline(chargeId string) {
	_, ok := pwm.chargeMap[chargeId]
	if ok {
		return
	}
	pwm.lock.Lock()
	defer pwm.lock.Unlock()
	pwm.chargeMap[chargeId] = time.Now().Unix()
	seelog.Debug("Pingxx webhook | SetChargeOnline:", chargeId)
}

func (pwm *PingxxWebhookManager) IsChargeOnline(chargeId string) bool {
	pwm.lock.RLock()
	_, ok := pwm.chargeMap[chargeId]
	defer pwm.lock.RUnlock()
	seelog.Debug("Pingxx webhook | IsChargeOnline:", chargeId, " ", ok)
	return ok
}

func (pwm *PingxxWebhookManager) ChargeSuccessEvent(chargeId string) {
	if !pwm.IsChargeOnline(chargeId) {
		pwm.SetChargeOnline(chargeId)
	} else {
		return
	}
	var err error

	recordInfo := map[string]interface{}{
		"Result": "success",
	}
	models.UpdatePingppRecord(chargeId, recordInfo)
	record, err := models.QueryPingppRecordByChargeId(chargeId)

	if err != nil {
		return
	}

	if pwm.checkChargeSuccessExist(record) {
		return
	}

	premium, err := trade.GetChargePremuim(record.UserId, int64(record.Amount))
	trade.HandleTradeChargePingpp(record.Id)
	if premium > 0 {
		trade.HandleTradeChargePremium(record.UserId, premium, "", record.Id, "")
	}
}

func (pwm *PingxxWebhookManager) RefundSuccessEvent(chargeId string, refundId string) {
	recordInfo := map[string]interface{}{
		"Result":   "success",
		"RefundId": refundId,
	}
	models.UpdatePingppRecord(chargeId, recordInfo)
}

func (pwm *PingxxWebhookManager) checkChargeSuccessExist(record *models.PingppRecord) bool {
	o := orm.NewOrm()

	exist := o.QueryTable(new(models.TradeRecord).TableName()).
		Filter("user_id", record.UserId).
		Filter("trade_type", models.TRADE_CHARGE).
		Filter("pingpp_id", record.Id).Exist()

	return exist
}
