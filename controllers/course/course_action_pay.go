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

	course, err := models.ReadCourse(courseId)
	if err != nil {
		return 2, errors.New("课程信息异常")
	}
	user, err := models.ReadUser(userId)
	if err != nil {
		return 2, errors.New("用户信息异常")
	}
	if course.Type == models.COURSE_TYPE_DELUXE {
		// 先查询该用户是否有购买（或试图购买）过这个课程
		var currentRecord models.CoursePurchaseRecord
		var record *models.CoursePurchaseRecord
		err = o.QueryTable(new(models.CoursePurchaseRecord).TableName()).Filter("course_id", courseId).Filter("user_id", userId).
			One(&currentRecord)
		if err != nil {
			return 2, errors.New("购买记录异常")
		}
		record = &currentRecord

		switch payType {
		case PAYMENT_TYPE_AUDITION:
			if record.AuditionStatus != models.PURCHASE_RECORD_STATUS_WAITING {
				return 2, errors.New("购买记录异常")
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

			err = trade.HandleUserBalance(record.UserId, 0-record.PriceTotal)
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
	} else if course.Type == models.COURSE_TYPE_AUDITION {
		var currentRecord models.CourseAuditionRecord
		o.QueryTable(new(models.CourseAuditionRecord).TableName()).
			Filter("course_id", courseId).Filter("user_id", userId).
			Exclude("status", models.AUDITION_RECORD_STATUS_COMPLETE).
			One(&currentRecord)
		if currentRecord.Id == 0 {
			return 2, errors.New("试听记录异常")
		} else {
			if currentRecord.Status != models.AUDITION_RECORD_STATUS_APPLY && currentRecord.Status != models.AUDITION_RECORD_STATUS_WAITING {
				return 2, errors.New("试听记录异常")
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
				return 2, errors.New("购买记录异常")
			}
		}
	}
	return 0, nil
}

func CheckCourseActionPayByThird(userId int64, courseId int64, tradeType string) (int64, error) {
	var err error
	o := orm.NewOrm()

	course, err := models.ReadCourse(courseId)
	if err != nil {
		return 2, errors.New("课程信息异常")
	}
	if course.Type == models.COURSE_TYPE_DELUXE {
		// 先查询该用户是否有购买（或试图购买）过这个课程
		var currentRecord models.CoursePurchaseRecord
		err = o.QueryTable("course_purchase_record").Filter("course_id", courseId).Filter("user_id", userId).
			One(&currentRecord)
		if err != nil {
			return 2, errors.New("购买记录异常")
		}
		switch tradeType {
		case models.TRADE_COURSE_AUDITION:
			if currentRecord.AuditionStatus != models.PURCHASE_RECORD_STATUS_WAITING {
				return 2, errors.New("购买记录异常")
			}
		case models.TRADE_COURSE_PURCHASE:
			if currentRecord.PurchaseStatus != models.PURCHASE_RECORD_STATUS_WAITING {
				return 2, errors.New("购买记录异常")
			}
		}
	} else if course.Type == models.COURSE_TYPE_AUDITION {
		var currentRecord models.CourseAuditionRecord
		o.QueryTable(new(models.CourseAuditionRecord).TableName()).
			Filter("course_id", courseId).Filter("user_id", userId).
			Exclude("status", models.AUDITION_RECORD_STATUS_COMPLETE).
			One(&currentRecord)
		if currentRecord.Id == 0 {
			return 2, errors.New("试听记录异常")
		} else {
			if currentRecord.Status != models.AUDITION_RECORD_STATUS_APPLY && currentRecord.Status != models.AUDITION_RECORD_STATUS_WAITING {
				return 2, errors.New("试听记录异常")
			}
		}
	}
	return 0, nil
}

func HandleCourseActionPayByThird(userId int64, courseId int64, tradeType string, pingppAmount, pingppId int64) (int64, error) {
	var err error
	o := orm.NewOrm()
	course, err := models.ReadCourse(courseId)
	if err != nil {
		return 2, errors.New("课程信息异常")
	}
	user, err := models.ReadUser(userId)
	if err != nil {
		return 2, errors.New("用户信息异常")
	}
	if course.Type == models.COURSE_TYPE_DELUXE {
		// 先查询该用户是否有购买（或试图购买）过这个课程
		var currentRecord models.CoursePurchaseRecord
		var record *models.CoursePurchaseRecord
		err = o.QueryTable("course_purchase_record").Filter("course_id", courseId).Filter("user_id", userId).
			One(&currentRecord)
		if err != nil {
			return 2, errors.New("购买记录异常")
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
				"audition_status": models.PURCHASE_RECORD_STATUS_PAID,
			}

			record, err = models.UpdateCoursePurchaseRecord(record.Id, recordInfo)
			if err != nil {
				return 2, errors.New("购买记录异常")
			}

		case models.TRADE_COURSE_PURCHASE:
			if pingppAmount < record.PriceTotal {
				err = trade.HandleUserBalance(userId, 0-user.Balance)
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
	} else if course.Type == models.COURSE_TYPE_AUDITION {
		if pingppAmount < PAYMENT_PRICE_AUDITION {
			err = trade.HandleUserBalance(userId, 0-user.Balance)
			if err != nil {
				return 2, err
			}
		}
		var currentRecord models.CourseAuditionRecord
		o.QueryTable(new(models.CourseAuditionRecord).TableName()).
			Filter("course_id", courseId).Filter("user_id", userId).
			Exclude("status", models.AUDITION_RECORD_STATUS_COMPLETE).
			One(&currentRecord)
		if currentRecord.Id == 0 {
			return 2, errors.New("试听记录异常")
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
			return 2, errors.New("购买记录异常")
		}
	}

	return 0, nil
}
