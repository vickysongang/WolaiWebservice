// qa_pkg_purchase
package qapkg

import (
	"WolaiWebservice/models"
	"errors"
	"time"
)

func HandleQaPkgPurchaseRecord(userId, qaPkgId int64) (int64, error) {
	qaPkg, err := models.ReadQaPkg(qaPkgId)
	if err != nil {
		return 2, errors.New("答疑包资料异常")
	}
	if qaPkg.Type == models.QA_PKG_TYPE_PERMANENT {
		purchaseRecord := models.QaPkgPurchaseRecord{
			QaPkgId:    qaPkgId,
			TimeLength: qaPkg.TimeLength,
			Price:      qaPkg.DiscountPrice,
			UserId:     userId,
			Type:       models.QA_PKG_TYPE_PERMANENT,
			Status:     models.QA_PKG_PURCHASE_STATUS_SERVING,
		}
		models.InsertQaPkgPurchaseRecord(&purchaseRecord)
	} else if qaPkg.Type == models.QA_PKG_TYPE_MONTHLY {
		latestRecord, _ := GetLatestMonthlyQaPkg(userId)
		var startTime time.Time
		if latestRecord.Id == 0 {
			startTime = time.Now()
		} else {
			startTime = latestRecord.TimeTo
		}
		for i := 0; i < int(qaPkg.Month); i++ {
			purchaseRecord := models.QaPkgPurchaseRecord{
				QaPkgId:    qaPkgId,
				TimeLength: qaPkg.TimeLength,
				Price:      qaPkg.DiscountPrice,
				UserId:     userId,
				TimeFrom:   startTime.AddDate(0, 0, i*30),
				TimeTo:     startTime.AddDate(0, 0, (i+1)*30),
				Type:       models.QA_PKG_TYPE_MONTHLY,
				Status:     models.QA_PKG_PURCHASE_STATUS_SERVING,
			}
			models.InsertQaPkgPurchaseRecord(&purchaseRecord)
		}
	}
	return 0, nil
}
