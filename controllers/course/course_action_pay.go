package course

import (
	"errors"

	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
	"WolaiWebservice/service/trade"
)

func HandleCourseActionPayByBalance(userId int64, courseId int64, payType string) (int64, error) {
	var err error
	o := orm.NewOrm()

	_, err = models.ReadCourse(courseId)
	if err != nil {
		return 2, errors.New("课程信息异常")
	}

	// 先查询该用户是否有购买（或试图购买）过这个课程
	var currentRecord models.CoursePurchaseRecord
	var record *models.CoursePurchaseRecord
	err = o.QueryTable("course_purchase_record").Filter("course_id", courseId).Filter("user_id", userId).
		One(&currentRecord)
	if err != nil {
		return 2, errors.New("购买记录异常")
	}
	record = &currentRecord

	user, err := models.ReadUser(record.UserId)
	if err != nil {
		return 2, errors.New("用户信息异常")
	}

	switch payType {
	case PAYMENT_TYPE_AUDITION:
		if record.AuditionStatus != models.PURCHASE_RECORD_STATUS_WAITING {
			return 2, errors.New("购买记录异常")
		}

		if user.Balance < PAYMENT_PRICE_AUDITION {
			return 2, trade.ErrInsufficientFund
		}

		err = trade.AddUserBalance(record.UserId, 0-PAYMENT_PRICE_AUDITION)
		if err != nil {
			return 2, err
		}

		err = trade.HandleCourseAuditionTradeRecord(record.Id, PAYMENT_PRICE_AUDITION, 0)
		if err != nil {
			return 2, err
		}

		recordInfo := map[string]interface{}{
			"audition_status": models.PURCHASE_RECORD_STATUS_PAID,
		}

		record, err = models.UpdateCoursePurchaseRecord(record.Id, recordInfo)
		if err != nil {
			return 2, errors.New("购买记录异常")
		}

	case PAYMENT_TYPE_PURCHASE:
		if record.PurchaseStatus != models.PURCHASE_RECORD_STATUS_WAITING {
			return 2, errors.New("购买记录异常")
		}
		if user.Balance < record.PriceTotal {
			return 2, trade.ErrInsufficientFund
		}

		err = trade.AddUserBalance(record.UserId, 0-record.PriceTotal)
		if err != nil {
			return 2, err
		}

		err = trade.HandleCoursePurchaseTradeRecord(record.Id, 0)
		if err != nil {
			return 2, err
		}

		recordInfo := map[string]interface{}{
			"purchase_status": models.PURCHASE_RECORD_STATUS_PAID,
		}

		if record.AuditionStatus == models.PURCHASE_RECORD_STATUS_IDLE ||
			record.AuditionStatus == models.PURCHASE_RECORD_STATUS_APPLY ||
			record.AuditionStatus == models.PURCHASE_RECORD_STATUS_WAITING {
			recordInfo["audition_status"] = models.PURCHASE_RECORD_STATUS_PAID
		}

		record, err = models.UpdateCoursePurchaseRecord(record.Id, recordInfo)
		if err != nil {
			return 2, errors.New("购买记录异常")
		}
	}

	return 0, nil
}

func CheckCourseActionPayByThird(userId int64, courseId int64, tradeType string) (int64, error) {
	var err error
	o := orm.NewOrm()

	_, err = models.ReadCourse(courseId)
	if err != nil {
		return 2, errors.New("课程信息异常")
	}
	// 先查询该用户是否有购买（或试图购买）过这个课程
	var currentRecord models.CoursePurchaseRecord
	var record *models.CoursePurchaseRecord
	err = o.QueryTable("course_purchase_record").Filter("course_id", courseId).Filter("user_id", userId).
		One(&currentRecord)
	if err != nil {
		return 2, errors.New("购买记录异常")
	}
	switch tradeType {
	case models.TRADE_COURSE_AUDITION:
		if record.AuditionStatus != models.PURCHASE_RECORD_STATUS_WAITING {
			return 2, errors.New("购买记录异常")
		}
	case models.TRADE_COURSE_PURCHASE:
		if record.PurchaseStatus != models.PURCHASE_RECORD_STATUS_WAITING {
			return 2, errors.New("购买记录异常")
		}
	}
	return 0, nil
}

func HandleCourseActionPayByThird(userId int64, courseId int64, tradeType string, pingppAmount, pingppId int64) (int64, error) {
	var err error
	o := orm.NewOrm()

	// 先查询该用户是否有购买（或试图购买）过这个课程
	var currentRecord models.CoursePurchaseRecord
	var record *models.CoursePurchaseRecord
	err = o.QueryTable("course_purchase_record").Filter("course_id", courseId).Filter("user_id", userId).
		One(&currentRecord)
	if err != nil {
		return 2, errors.New("购买记录异常")
	}
	user, err := models.ReadUser(record.UserId)
	if err != nil {
		return 2, errors.New("用户信息异常")
	}
	record = &currentRecord
	switch tradeType {
	case models.TRADE_COURSE_AUDITION:
		if pingppAmount < PAYMENT_PRICE_AUDITION {
			err = trade.AddUserBalance(userId, 0-user.Balance)
			if err != nil {
				return 2, err
			}
		}
		err = trade.HandleCourseAuditionTradeRecord(record.Id, PAYMENT_PRICE_AUDITION, pingppId)
		if err != nil {
			return 2, err
		}

		recordInfo := map[string]interface{}{
			"audition_status": models.PURCHASE_RECORD_STATUS_PAID,
		}

		record, err = models.UpdateCoursePurchaseRecord(record.Id, recordInfo)
		if err != nil {
			return 2, errors.New("购买记录异常")
		}

	case models.TRADE_COURSE_PURCHASE:
		if pingppAmount < record.PriceTotal {
			err = trade.AddUserBalance(userId, 0-user.Balance)
			if err != nil {
				return 2, err
			}
		}
		err = trade.HandleCoursePurchaseTradeRecord(record.Id, pingppId)
		if err != nil {
			return 2, err
		}

		recordInfo := map[string]interface{}{
			"purchase_status": models.PURCHASE_RECORD_STATUS_PAID,
		}

		if record.AuditionStatus == models.PURCHASE_RECORD_STATUS_IDLE ||
			record.AuditionStatus == models.PURCHASE_RECORD_STATUS_APPLY ||
			record.AuditionStatus == models.PURCHASE_RECORD_STATUS_WAITING {
			recordInfo["audition_status"] = models.PURCHASE_RECORD_STATUS_PAID
		}

		record, err = models.UpdateCoursePurchaseRecord(record.Id, recordInfo)
		if err != nil {
			return 2, errors.New("购买记录异常")
		}
	}
	return 0, nil
}
