// course_quota
package course

import (
	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
	"WolaiWebservice/service/trade"
)

type QuotaChargeRule struct {
	Price     int64                        `json:"price"`
	Discounts []models.CourseQuotaDiscount `json:"discounts"`
}

func GetQuotaChargeRule(gradeId int64) (*QuotaChargeRule, error) {
	chargeRule := QuotaChargeRule{}
	quotaPrice, err := courseService.GetCourseQuotaPrice(gradeId)
	if err != nil {
		return nil, err
	}
	chargeRule.Price = quotaPrice.Price
	discounts, err := courseService.QueryCourseQuotaDiscounts()
	if err == nil {
		chargeRule.Discounts = discounts
	}
	return &chargeRule, nil
}

func HandleCourseQuotaActionPayByBalance(userId, gradeId, quantity, amount int64) (int64, error) {
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
	status, err := handleCourseQuotaActionPay(userId, gradeId, quantity, amount, 0)
	return status, err
}

func HandleCourseQuotaActionPayByThird(userId, gradeId, quantity, pingppAmount, totalAmount, pingppId int64) (int64, error) {
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
	status, err := handleCourseQuotaActionPay(userId, gradeId, quantity, totalAmount, pingppId)
	return status, err
}

func handleCourseQuotaActionPay(userId, gradeId, quantity, amount, pingppId int64) (int64, error) {
	var err error
	_, err = models.ReadUser(userId)
	if err != nil {
		return 2, ErrUserAbnormal
	}
	profile, err := models.ReadStudentProfile(userId)
	if err != nil {
		return 2, ErrUserAbnormal
	}
	if profile.QuotaGradeId == 0 {
		profile.QuotaGradeId = gradeId
	}
	profile.QuotaQuantity += quantity
	quotaPrice, _ := courseService.GetCourseQuotaPrice(gradeId)
	quotaDiscount, _ := courseService.QueryQuotaDiscountByQuantity(quantity)
	_, err = models.UpdateStudentProfile(profile)
	if err != nil {
		return 2, err
	}
	courseQuotaTradeRecord := models.CourseQuotaTradeRecord{
		UserId:       userId,
		GradeId:      gradeId,
		Price:        quotaPrice.Price,
		TotalPrice:   amount,
		Discount:     quotaDiscount.Discount,
		Quantity:     quantity,
		LeftQuantity: quantity,
		Type:         models.COURSE_QUOTA_TYPE_ONLINE_PURCHASE,
	}
	recordId, err := models.InsertCourseQuotaTradeRecord(&courseQuotaTradeRecord)
	if err != nil {
		return 2, err
	}
	err = trade.HandleCourseQuotaPurchaseTradeRecord(recordId, amount, pingppId)
	if err != nil {
		return 2, err
	}
	return 0, nil
}
