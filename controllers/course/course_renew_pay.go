// course_renew_pay
package course

import (
	"WolaiWebservice/models"
	"errors"

	"github.com/astaxie/beego/orm"
)

func HandleCourseRenewPay(userId, courseId, amount int64) (int64, error) {
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
	//	chapterCount := amount / currentRecord.PriceHourly

	return 0, nil
}
