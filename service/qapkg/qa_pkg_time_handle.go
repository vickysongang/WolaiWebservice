package qapkg

import (
	"WolaiWebservice/models"
	"time"
)

func HandleUserQaPkgTime(userId int64, timeLength int64) error {
	now := time.Now()
	balanceTime := timeLength
	monthlyQaPkgRecords, _ := GetMonthlyQaPkgRecords(userId)
	for _, monthlyQaPkgRecord := range monthlyQaPkgRecords {
		if now.After(monthlyQaPkgRecord.TimeFrom) && monthlyQaPkgRecord.TimeTo.After(now) {
			if timeLength < monthlyQaPkgRecord.LeftTime {
				recordInfo := map[string]interface{}{
					"LeftTime": monthlyQaPkgRecord.LeftTime - timeLength,
				}
				models.UpdateQaPkgPurchaseRecord(monthlyQaPkgRecord.Id, recordInfo)
				return nil
			} else {
				recordInfo := map[string]interface{}{
					"LeftTime": 0,
					"Status":   models.QA_PKG_PURCHASE_STATUS_COMPLETE,
				}
				models.UpdateQaPkgPurchaseRecord(monthlyQaPkgRecord.Id, recordInfo)
				balanceTime = timeLength - monthlyQaPkgRecord.LeftTime
				break
			}
		}
	}
	permanentQaPkgRecords, _ := GetPermanentQaPkgRecords(userId)
	for _, permanentQaPkgRecord := range permanentQaPkgRecords {
		if balanceTime < permanentQaPkgRecord.LeftTime {
			recordInfo := map[string]interface{}{
				"LeftTime": permanentQaPkgRecord.LeftTime - balanceTime,
			}
			models.UpdateQaPkgPurchaseRecord(permanentQaPkgRecord.Id, recordInfo)
			return nil
		} else {
			recordInfo := map[string]interface{}{
				"LeftTime": 0,
				"Status":   models.QA_PKG_PURCHASE_STATUS_COMPLETE,
			}
			models.UpdateQaPkgPurchaseRecord(permanentQaPkgRecord.Id, recordInfo)
			balanceTime = balanceTime - permanentQaPkgRecord.LeftTime
			continue
		}
	}
	return nil
}
