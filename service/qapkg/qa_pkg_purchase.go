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
			LeftTime:   qaPkg.TimeLength,
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
				LeftTime:   qaPkg.TimeLength,
				Type:       models.QA_PKG_TYPE_MONTHLY,
				Status:     models.QA_PKG_PURCHASE_STATUS_SERVING,
			}
			models.InsertQaPkgPurchaseRecord(&purchaseRecord)
		}
	}
	return 0, nil
}

func HandleGivenQaPkgPurchaseRecord(userId, qaPkgId int64) (int64, error) {
	qaPkg, err := models.ReadQaPkg(qaPkgId)
	if err != nil {
		return 2, errors.New("答疑包资料异常")
	}
	endTime, _ := time.Parse(time.RFC3339, "2100-01-01T00:00:00.00+08:00")
	purchaseRecord := models.QaPkgPurchaseRecord{
		QaPkgId:    qaPkg.Id,
		TimeLength: qaPkg.TimeLength,
		Price:      qaPkg.DiscountPrice,
		UserId:     userId,
		LeftTime:   qaPkg.TimeLength,
		TimeFrom:   time.Now(),
		TimeTo:     endTime,
		Type:       models.QA_PKG_TYPE_GIVEN,
		Status:     models.QA_PKG_PURCHASE_STATUS_SERVING,
	}
	_, err = models.InsertQaPkgPurchaseRecord(&purchaseRecord)
	if err != nil {
		return 2, err
	}
	return 0, nil
}
