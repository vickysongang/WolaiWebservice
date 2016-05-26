package qapkg

import (
	"WolaiWebservice/models"
	"time"
)

func HandleUserQaPkgTime(userId int64, timeLength int64) error {
	now := time.Now()
	balanceTime := timeLength
	monthlyQaPkgRecords, _ := GetMonthlyQaPkgRecords(userId)
	for _, record := range monthlyQaPkgRecords {
		if now.After(record.TimeFrom) && record.TimeTo.After(now) {
			if timeLength < record.LeftTime {
				recordInfo := map[string]interface{}{
					"LeftTime": record.LeftTime - timeLength,
				}
				models.UpdateQaPkgPurchaseRecord(record.Id, recordInfo)
				return nil
			} else {
				recordInfo := map[string]interface{}{
					"LeftTime": 0,
					"Status":   models.QA_PKG_PURCHASE_STATUS_COMPLETE,
				}
				models.UpdateQaPkgPurchaseRecord(record.Id, recordInfo)
				balanceTime = timeLength - record.LeftTime
				break
			}
		}
	}
	permanentQaPkgRecords, _ := GetPermanentQaPkgRecords(userId)
	for _, record := range permanentQaPkgRecords {
		if record.Type == models.QA_PKG_TYPE_GIVEN && !(now.After(record.TimeFrom) && record.TimeTo.After(now)) {
			continue
		}
		if balanceTime < record.LeftTime {
			recordInfo := map[string]interface{}{
				"LeftTime": record.LeftTime - balanceTime,
			}
			models.UpdateQaPkgPurchaseRecord(record.Id, recordInfo)
			return nil
		} else {
			recordInfo := map[string]interface{}{
				"LeftTime": 0,
				"Status":   models.QA_PKG_PURCHASE_STATUS_COMPLETE,
			}
			models.UpdateQaPkgPurchaseRecord(record.Id, recordInfo)
			balanceTime = balanceTime - record.LeftTime
			continue
		}
	}
	return nil
}
