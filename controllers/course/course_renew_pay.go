// course_renew_pay
package course

import (
	"WolaiWebservice/models"
	"errors"

	courseService "WolaiWebservice/service/course"
	"WolaiWebservice/service/trade"
)

func HandleCourseRenewPayByBalance(userId, courseId, amount, quantity int64) (int64, error) {
	user, err := models.ReadUser(userId)
	if err != nil {
		return 2, ErrUserAbnormal
	}
	if user.Balance < amount {
		return 2, trade.ErrInsufficientFund
	}
	err = trade.HandleUserBalance(userId, 0-amount)
	if err != nil {
		return 2, err
	}
	status, err := handleCourseRenewPay(userId, courseId, amount, quantity, 0)
	return status, err
}

func HandleCourseRenewPayByThird(userId, courseId, pingppAmount, totalAmount, quantity, pingppId int64) (int64, error) {
	user, err := models.ReadUser(userId)
	if err != nil {
		return 2, ErrUserAbnormal
	}
	if pingppAmount < totalAmount {
		err = trade.HandleUserBalance(userId, 0-user.Balance)
		if err != nil {
			return 2, err
		}
	}
	status, err := handleCourseRenewPay(userId, courseId, totalAmount, quantity, pingppId)
	return status, err
}

func handleCourseRenewPay(userId, courseId, amount, quantity int64, pingppId int64) (int64, error) {
	_, err := models.ReadCourse(courseId)
	if err != nil {
		return 2, errors.New("课程包资料异常")
	}
	_, err = models.ReadUser(userId)
	if err != nil {
		return 2, ErrUserAbnormal
	}
	currentRecord, err := courseService.GetCoursePurchaseRecordByUserId(courseId, userId)
	if err != nil {
		return 2, ErrPurchaseAbnormal
	}
	var renewCount int64
	var renewRecordId int64
	oldRenewRecord := courseService.GetCourseRenewWaitingRecord(userId, courseId)
	if oldRenewRecord != nil && oldRenewRecord.PriceTotal == amount {
		renewRecordInfo := map[string]interface{}{
			"Status": models.COURSE_RENEW_STATUS_COMPLETE,
		}
		err = models.UpdateCourseRenewRecord(oldRenewRecord.Id, renewRecordInfo)
		if err != nil {
			return 2, err
		}
		renewCount = oldRenewRecord.RenewCount
		renewRecordId = oldRenewRecord.Id
	} else {
		chapterCount := amount / currentRecord.PriceHourly
		if quantity != 0 {
			chapterCount = quantity
		}
		newRenewRecord := models.CourseRenewRecord{
			CourseId:      courseId,
			UserId:        userId,
			TeacherId:     currentRecord.TeacherId,
			PriceHourly:   currentRecord.PriceHourly,
			PriceTotal:    amount,
			RenewCount:    chapterCount,
			Status:        models.COURSE_RENEW_STATUS_COMPLETE,
			PaymentMethod: models.PAYMENT_METHOD_ONLINE_WALLET,
		}
		id, err := models.CreateCourseRenewRecord(&newRenewRecord)
		if err != nil {
			return 2, err
		}
		renewCount = chapterCount
		renewRecordId = id
	}
	err = trade.HandleCourseRenewTradeRecord(renewRecordId, amount, pingppId)
	if err != nil {
		return 2, err
	}
	purchaseRecordInfo := map[string]interface{}{
		"ChapterCount":   currentRecord.ChapterCount + renewCount,
		"PurchaseStatus": models.AUDITION_RECORD_STATUS_PAID,
	}
	_, err = models.UpdateCoursePurchaseRecord(currentRecord.Id, purchaseRecordInfo)
	if err != nil {
		return 2, ErrPurchaseAbnormal
	}
	return 0, nil
}

func HandleCourseRenewPayByQuota(userId, courseId, quantity int64) (int64, error) {
	course, err := models.ReadCourse(courseId)
	if err != nil {
		return 2, errors.New("课程包资料异常")
	}
	_, err = models.ReadUser(userId)
	if err != nil {
		return 2, ErrUserAbnormal
	}
	record, err := courseService.GetCoursePurchaseRecordByUserId(courseId, userId)
	if err != nil {
		return 2, ErrPurchaseAbnormal
	}
	var renewRecordId int64

	oldRenewRecord := courseService.GetCourseRenewWaitingRecord(userId, courseId)
	if oldRenewRecord != nil && oldRenewRecord.RenewCount == quantity {
		renewRecordInfo := map[string]interface{}{
			"Status":        models.COURSE_RENEW_STATUS_COMPLETE,
			"PaymentMethod": models.PAYMENT_METHOD_ONLINE_QUOTA,
		}
		err = models.UpdateCourseRenewRecord(oldRenewRecord.Id, renewRecordInfo)
		if err != nil {
			return 2, err
		}
		renewRecordId = oldRenewRecord.Id
	} else {
		newRenewRecord := models.CourseRenewRecord{
			CourseId:      courseId,
			UserId:        userId,
			TeacherId:     record.TeacherId,
			PriceHourly:   record.PriceHourly,
			RenewCount:    quantity,
			Status:        models.COURSE_RENEW_STATUS_COMPLETE,
			PaymentMethod: models.PAYMENT_METHOD_ONLINE_QUOTA,
		}
		id, err := models.CreateCourseRenewRecord(&newRenewRecord)
		if err != nil {
			return 2, err
		}
		renewRecordId = id
	}
	totalPrice, err := courseService.HandleCourseQuotaPay(userId, renewRecordId, course.GradeId, quantity, "renew")
	if err != nil {
		return 2, err
	}

	renewInfo := map[string]interface{}{
		"PriceTotal": totalPrice,
	}
	err = models.UpdateCourseRenewRecord(renewRecordId, renewInfo)
	if err != nil {
		return 2, err
	}

	err = trade.HandleCourseRenewTradeRecord(renewRecordId, totalPrice, 0)
	if err != nil {
		return 2, err
	}

	purchaseRecordInfo := map[string]interface{}{
		"ChapterCount":   record.ChapterCount + quantity,
		"PurchaseStatus": models.AUDITION_RECORD_STATUS_PAID,
	}
	_, err = models.UpdateCoursePurchaseRecord(record.Id, purchaseRecordInfo)
	if err != nil {
		return 2, ErrPurchaseAbnormal
	}

	return 0, nil
}
