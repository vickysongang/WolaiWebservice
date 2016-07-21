package course

import (
	"errors"

	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
	"WolaiWebservice/service/trade"
)

func HandleCourseActionPayByBalance(userId int64, courseId int64, payType string) (int64, error) {
	var err error

	course, err := models.ReadCourse(courseId)
	if err != nil {
		return 2, ErrCourseAbnormal
	}
	user, err := models.ReadUser(userId)
	if err != nil {
		return 2, ErrUserAbnormal
	}
	if course.Type == models.COURSE_TYPE_DELUXE {
		// 先查询该用户是否有购买（或试图购买）过这个课程
		var record *models.CoursePurchaseRecord
		currentRecord, err := courseService.GetCoursePurchaseRecordByUserId(courseId, userId)
		if err != nil {
			return 2, ErrPurchaseAbnormal
		}
		record = &currentRecord

		switch payType {
		case PAYMENT_TYPE_AUDITION:
			if record.AuditionStatus != models.PURCHASE_RECORD_STATUS_WAITING {
				return 2, ErrPurchaseAbnormal
			}
			if user.Balance < PAYMENT_PRICE_AUDITION {
				return 2, trade.ErrInsufficientFund
			}

			err = trade.HandleUserBalance(record.UserId, 0-PAYMENT_PRICE_AUDITION)
			if err != nil {
				return 2, err
			}

			err = trade.HandleCourseAuditionTradeRecord(record.Id, PAYMENT_PRICE_AUDITION, 0)
			if err != nil {
				return 2, err
			}

			recordInfo := map[string]interface{}{
				"AuditionStatus": models.PURCHASE_RECORD_STATUS_PAID,
				"PaymentMethod":  models.PAYMENT_METHOD_ONLINE_WALLET,
			}

			record, err = models.UpdateCoursePurchaseRecord(record.Id, recordInfo)
			if err != nil {
				return 2, ErrPurchaseAbnormal
			}

		case PAYMENT_TYPE_PURCHASE:
			if record.PurchaseStatus != models.PURCHASE_RECORD_STATUS_WAITING {
				return 2, ErrPurchaseAbnormal
			}
			if user.Balance < record.PriceTotal {
				return 2, trade.ErrInsufficientFund
			}

			err = trade.HandleUserBalance(record.UserId, 0-record.PriceTotal)
			if err != nil {
				return 2, err
			}

			err = trade.HandleCoursePurchaseTradeRecord(record.Id, 0, "")
			if err != nil {
				return 2, err
			}

			recordInfo := map[string]interface{}{
				"PurchaseStatus": models.PURCHASE_RECORD_STATUS_PAID,
				"PaymentMethod":  models.PAYMENT_METHOD_ONLINE_WALLET,
			}

			if record.AuditionStatus == models.PURCHASE_RECORD_STATUS_IDLE ||
				record.AuditionStatus == models.PURCHASE_RECORD_STATUS_APPLY ||
				record.AuditionStatus == models.PURCHASE_RECORD_STATUS_WAITING {
				recordInfo["AuditionStatus"] = models.PURCHASE_RECORD_STATUS_PAID
			}

			record, err = models.UpdateCoursePurchaseRecord(record.Id, recordInfo)
			if err != nil {
				return 2, ErrPurchaseAbnormal
			}
		}
	} else if course.Type == models.COURSE_TYPE_AUDITION {
		currentRecord, _ := courseService.GetCourseAuditionRecordByUserId(courseId, userId)
		if currentRecord.Id == 0 {
			return 2, errors.New("试听记录不存在")
		} else {
			if currentRecord.Status != models.AUDITION_RECORD_STATUS_APPLY &&
				currentRecord.Status != models.AUDITION_RECORD_STATUS_WAITING {
				return 2, ErrAuditionAbnormal
			}
			if user.Balance < PAYMENT_PRICE_AUDITION {
				return 2, trade.ErrInsufficientFund
			}

			err = trade.HandleUserBalance(currentRecord.UserId, 0-PAYMENT_PRICE_AUDITION)
			if err != nil {
				return 2, err
			}

			err = trade.HandleAuditionCoursePurchaseTradeRecord(currentRecord.Id, PAYMENT_PRICE_AUDITION, 0)
			if err != nil {
				return 2, err
			}
			recordInfo := map[string]interface{}{
				"Status": models.AUDITION_RECORD_STATUS_PAID,
			}
			_, err := models.UpdateCourseAuditionRecord(currentRecord.Id, recordInfo)
			if err != nil {
				return 2, errors.New("更新试听记录异常")
			}
		}
	}
	return 0, nil
}

func CheckCourseActionPayByThird(userId int64, courseId int64, tradeType string) (int64, error) {
	var err error

	course, err := models.ReadCourse(courseId)
	if err != nil {
		return 2, ErrCourseAbnormal
	}
	if course.Type == models.COURSE_TYPE_DELUXE {
		// 先查询该用户是否有购买（或试图购买）过这个课程
		currentRecord, err := courseService.GetCoursePurchaseRecordByUserId(courseId, userId)
		if err != nil {
			return 2, ErrPurchaseAbnormal
		}
		switch tradeType {
		case models.TRADE_COURSE_AUDITION:
			if currentRecord.AuditionStatus != models.PURCHASE_RECORD_STATUS_WAITING {
				return 2, ErrPurchaseAbnormal
			}
		case models.TRADE_COURSE_PURCHASE:
			if currentRecord.PurchaseStatus != models.PURCHASE_RECORD_STATUS_WAITING {
				return 2, ErrPurchaseAbnormal
			}
		}
	} else if course.Type == models.COURSE_TYPE_AUDITION {
		currentRecord, _ := courseService.GetCourseAuditionRecordByUserId(courseId, userId)
		if currentRecord.Id == 0 {
			return 2, ErrAuditionAbnormal
		} else {
			if currentRecord.Status != models.AUDITION_RECORD_STATUS_APPLY &&
				currentRecord.Status != models.AUDITION_RECORD_STATUS_WAITING {
				return 2, ErrAuditionAbnormal
			}
		}
	}
	return 0, nil
}

func HandleCourseActionPayByThird(userId int64, courseId int64, tradeType string, pingppAmount, pingppId int64) (int64, error) {
	var err error
	course, err := models.ReadCourse(courseId)
	if err != nil {
		return 2, ErrCourseAbnormal
	}
	user, err := models.ReadUser(userId)
	if err != nil {
		return 2, ErrUserAbnormal
	}
	if course.Type == models.COURSE_TYPE_DELUXE {
		// 先查询该用户是否有购买（或试图购买）过这个课程
		var record *models.CoursePurchaseRecord
		currentRecord, err := courseService.GetCoursePurchaseRecordByUserId(courseId, userId)
		if err != nil {
			return 2, ErrPurchaseAbnormal
		}
		record = &currentRecord
		switch tradeType {
		case models.TRADE_COURSE_AUDITION:
			if pingppAmount < PAYMENT_PRICE_AUDITION {
				err = trade.HandleUserBalance(userId, 0-user.Balance)
				if err != nil {
					return 2, err
				}
			}
			err = trade.HandleCourseAuditionTradeRecord(record.Id, PAYMENT_PRICE_AUDITION, pingppId)
			if err != nil {
				return 2, err
			}

			recordInfo := map[string]interface{}{
				"AuditionStatus": models.PURCHASE_RECORD_STATUS_PAID,
				"PaymentMethod":  models.PAYMENT_METHOD_ONLINE_WALLET,
			}

			record, err = models.UpdateCoursePurchaseRecord(record.Id, recordInfo)
			if err != nil {
				return 2, ErrPurchaseAbnormal
			}

		case models.TRADE_COURSE_PURCHASE:
			if pingppAmount < record.PriceTotal {
				err = trade.HandleUserBalance(userId, 0-user.Balance)
				if err != nil {
					return 2, err
				}
			}
			err = trade.HandleCoursePurchaseTradeRecord(record.Id, pingppId, "")
			if err != nil {
				return 2, err
			}

			recordInfo := map[string]interface{}{
				"PurchaseStatus": models.PURCHASE_RECORD_STATUS_PAID,
				"PaymentMethod":  models.PAYMENT_METHOD_ONLINE_WALLET,
			}

			if record.AuditionStatus == models.PURCHASE_RECORD_STATUS_IDLE ||
				record.AuditionStatus == models.PURCHASE_RECORD_STATUS_APPLY ||
				record.AuditionStatus == models.PURCHASE_RECORD_STATUS_WAITING {
				recordInfo["AuditionStatus"] = models.PURCHASE_RECORD_STATUS_PAID
			}

			record, err = models.UpdateCoursePurchaseRecord(record.Id, recordInfo)
			if err != nil {
				return 2, ErrPurchaseAbnormal
			}
		}
	} else if course.Type == models.COURSE_TYPE_AUDITION {
		if pingppAmount < PAYMENT_PRICE_AUDITION {
			err = trade.HandleUserBalance(userId, 0-user.Balance)
			if err != nil {
				return 2, err
			}
		}
		currentRecord, _ := courseService.GetCourseAuditionRecordByUserId(courseId, userId)
		if currentRecord.Id == 0 {
			return 2, ErrAuditionAbnormal
		}
		err = trade.HandleAuditionCoursePurchaseTradeRecord(currentRecord.Id, PAYMENT_PRICE_AUDITION, pingppId)
		if err != nil {
			return 2, err
		}
		recordInfo := map[string]interface{}{
			"Status": models.AUDITION_RECORD_STATUS_PAID,
		}
		_, err := models.UpdateCourseAuditionRecord(currentRecord.Id, recordInfo)
		if err != nil {
			return 2, ErrPurchaseAbnormal
		}
	}

	return 0, nil
}

func HandleDeluxeCoursePayByQuota(userId, courseId int64) (int64, error) {
	var err error
	course, err := models.ReadCourse(courseId)
	if err != nil {
		return 2, ErrCourseAbnormal
	}
	_, err = models.ReadUser(userId)
	if err != nil {
		return 2, ErrUserAbnormal
	}
	record, err := courseService.GetCoursePurchaseRecordByUserId(courseId, userId)
	if err != nil {
		return 2, ErrPurchaseAbnormal
	}
	totalPrice, err := courseService.HandleCourseQuotaPay(userId, record.Id, course.GradeId, record.ChapterCount, "purchase")
	if err != nil {
		return 2, err
	}
	recordInfo := map[string]interface{}{
		"PurchaseStatus": models.PURCHASE_RECORD_STATUS_PAID,
		"PaymentMethod":  models.PAYMENT_METHOD_ONLINE_QUOTA,
	}
	_, err = models.UpdateCoursePurchaseRecord(record.Id, recordInfo)
	if err != nil {
		return 2, ErrPurchaseAbnormal
	}
	err = trade.HandleCoursePurchaseByQuotaTradeRecord(record.Id, totalPrice)
	if err != nil {
		return 2, err
	}

	return 0, nil
}
