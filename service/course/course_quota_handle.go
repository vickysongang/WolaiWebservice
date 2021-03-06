// course_quota_handle
package course

import (
	"WolaiWebservice/models"
	"errors"
)

func HandleCourseQuotaPay(userId, recordId, gradeId, chapterCount int64, recordType string) (int64, error) {
	profile, err := models.ReadStudentProfile(userId)
	if err != nil {
		return 0, err
	}
	if profile.QuotaQuantity < chapterCount {
		return 0, errors.New("可用课时不够")
	}
	profile.QuotaQuantity -= chapterCount
	_, err = models.UpdateStudentProfile(profile)
	if err != nil {
		return 0, err
	}
	records, _ := QueryCourseQuotaPurchaseRecords(userId)
	leftChapterCount := chapterCount
	var totalPrice int64
	breakFlag := false
	for _, record := range records {
		var quantity int64
		if record.LeftQuantity >= leftChapterCount {
			quantity = leftChapterCount
			record.LeftQuantity -= leftChapterCount
			models.UpdateCourseQuotaTradeRecord(record)
			breakFlag = true
		} else {
			quantity = record.LeftQuantity
			record.LeftQuantity = 0
			models.UpdateCourseQuotaTradeRecord(record)
			leftChapterCount = leftChapterCount - record.LeftQuantity
		}
		if quantity > 0 {
			totalPriceItem := quantity * record.Price * record.Discount / 100
			totalPrice += totalPriceItem
			paymentDetail := models.CourseQuotaPaymentDetail{
				RecordId:         record.Id,
				CourseRecordId:   recordId,
				CourseRecordType: recordType,
				Quantity:         quantity,
				TotalPrice:       totalPriceItem,
			}
			models.InsertCourseQuotaPaymentDetail(&paymentDetail)
		}
		if breakFlag {
			break
		}
	}
	quotaPrice, _ := GetCourseQuotaPrice(gradeId)
	var quotaPayType string
	switch recordType {
	case "purchase":
		quotaPayType = models.COURSE_QUOTA_TYPE_QUOTA_PAY_PURCHASE
	case "renew":
		quotaPayType = models.COURSE_QUOTA_TYPE_QUOTA_PAY_RENEW
	}
	tradeRecord := models.CourseQuotaTradeRecord{
		UserId:         userId,
		GradeId:        gradeId,
		Price:          quotaPrice.Price,
		TotalPrice:     totalPrice,
		Discount:       0,
		Quantity:       chapterCount,
		LeftQuantity:   0,
		CourseRecordId: recordId,
		Type:           quotaPayType,
	}
	_, err = models.InsertCourseQuotaTradeRecord(&tradeRecord)
	if err != nil {
		return 0, err
	}
	return totalPrice, nil
}
