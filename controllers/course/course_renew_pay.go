// course_renew_pay
package course

import (
	"WolaiWebservice/models"
	"errors"

	courseService "WolaiWebservice/service/course"
	"WolaiWebservice/service/trade"

	"github.com/astaxie/beego/orm"
)

func HandleCourseRenewPayByBalance(userId, courseId, amount int64) (int64, error) {
	user, err := models.ReadUser(userId)
	if err != nil {
		return 2, errors.New("用户信息异常")
	}
	if user.Balance < amount {
		return 2, trade.ErrInsufficientFund
	}
	err = trade.HandleUserBalance(userId, 0-amount)
	if err != nil {
		return 2, err
	}
	status, err := HandleCourseRenewPay(userId, courseId, amount, 0)
	return status, err
}

func HandleCourseRenewPayByThird(userId, courseId, pingppAmount, totalAmount, pingppId int64) (int64, error) {
	user, err := models.ReadUser(userId)
	if err != nil {
		return 2, errors.New("用户信息异常")
	}
	if pingppAmount < totalAmount {
		err = trade.HandleUserBalance(userId, 0-user.Balance)
		if err != nil {
			return 2, err
		}
	}
	status, err := HandleCourseRenewPay(userId, courseId, totalAmount, pingppId)
	return status, err
}

func HandleCourseRenewPay(userId, courseId, amount int64, pingppId int64) (int64, error) {
	o := orm.NewOrm()
	_, err := models.ReadCourse(courseId)
	if err != nil {
		return 2, errors.New("课程包资料异常")
	}
	_, err = models.ReadUser(userId)
	if err != nil {
		return 2, errors.New("用户资料异常")
	}
	var currentRecord models.CoursePurchaseRecord
	err = o.QueryTable("course_purchase_record").Filter("course_id", courseId).Filter("user_id", userId).
		One(&currentRecord)
	if err != nil {
		return 2, errors.New("购买记录异常")
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
		newRenewRecord := models.CourseRenewRecord{
			CourseId:    courseId,
			UserId:      userId,
			TeacherId:   currentRecord.TeacherId,
			PriceHourly: currentRecord.PriceHourly,
			PriceTotal:  amount,
			RenewCount:  chapterCount,
			Status:      models.COURSE_RENEW_STATUS_COMPLETE,
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
		return 2, errors.New("购买记录异常")
	}
	return 0, nil
}
