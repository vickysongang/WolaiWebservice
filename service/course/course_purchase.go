// course_purchase
package course

import (
	"WolaiWebservice/models"

	"github.com/astaxie/beego/orm"
)

func GetCoursePurchaseRecordByUserId(courseId, userId int64) (models.CoursePurchaseRecord, error) {
	o := orm.NewOrm()
	var purchase models.CoursePurchaseRecord
	err := o.QueryTable("course_purchase_record").
		Filter("course_id", courseId).
		Filter("user_id", userId).
		One(&purchase)
	return purchase, err
}

func QueryCoursePurchaseRecords(userId, teacherId, page, count int64) ([]*models.CoursePurchaseRecord, error) {
	o := orm.NewOrm()
	var records []*models.CoursePurchaseRecord
	_, err := o.QueryTable("course_purchase_record").
		Filter("user_id", userId).
		Filter("teacher_id", teacherId).
		OrderBy("-last_update_time").
		Offset(page * count).Limit(count).
		All(&records)
	return records, err
}

func QueryTeacherCoursePurchaseRecords(teacherId, page, count int64, excludeIdle bool) ([]*models.CoursePurchaseRecord, error) {
	o := orm.NewOrm()
	var records []*models.CoursePurchaseRecord
	cond := orm.NewCondition()
	cond = cond.And("teacher_id", teacherId)
	if excludeIdle {
		cond.AndNot("purchase_status", models.PURCHASE_RECORD_STATUS_IDLE)
	}
	_, err := o.QueryTable("course_purchase_record").
		SetCond(cond).
		OrderBy("-last_update_time").
		Offset(page * count).Limit(count).
		All(&records)
	return records, err
}

func QueryStudentCoursePurchaseRecords(userId, page, count int64, excludeIdle bool) ([]*models.CoursePurchaseRecord, error) {
	o := orm.NewOrm()
	var records []*models.CoursePurchaseRecord
	cond := orm.NewCondition()
	cond = cond.And("user_id", userId)
	if excludeIdle {
		cond.AndNot("purchase_status", models.PURCHASE_RECORD_STATUS_IDLE)
	}
	_, err := o.QueryTable("course_purchase_record").
		SetCond(cond).
		OrderBy("-last_update_time").
		Offset(page * count).Limit(count).
		All(&records)
	return records, err
}

func GetStudentCoursePurchaseRecordCount(userId int64) int64 {
	o := orm.NewOrm()
	count, _ := o.QueryTable("course_purchase_record").
		Filter("user_id", userId).
		Exclude("purchase_status", models.PURCHASE_RECORD_STATUS_IDLE).
		Count()
	return count
}

func GetTeacherCoursePurchaseRecordCount(teacherId int64) int64 {
	o := orm.NewOrm()
	count, _ := o.QueryTable("course_purchase_record").
		Filter("teacher_id", teacherId).
		Exclude("purchase_status", models.PURCHASE_RECORD_STATUS_IDLE).
		Count()
	return count
}
