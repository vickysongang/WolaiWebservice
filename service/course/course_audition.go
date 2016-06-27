// course_audition
package course

import (
	"WolaiWebservice/models"

	"github.com/astaxie/beego/orm"
)

func QueryAuditionCourse() *models.Course {
	o := orm.NewOrm()
	var course models.Course
	o.QueryTable(new(models.Course).TableName()).Filter("type", models.COURSE_TYPE_AUDITION).One(&course)
	if course.Id == 0 {
		return nil
	}
	return &course
}

func GetUncompletedAuditionRecord(userId int64) *models.CourseAuditionRecord {
	o := orm.NewOrm()
	var record models.CourseAuditionRecord
	o.QueryTable(new(models.CourseAuditionRecord).TableName()).
		Filter("user_id", userId).Exclude("status", models.AUDITION_RECORD_STATUS_COMPLETE).One(&record)
	if record.Id == 0 {
		return nil
	}
	return &record
}

func QueryUncompletedAuditionRecords(userId, teacherId int64) ([]*models.CourseAuditionRecord, error) {
	o := orm.NewOrm()
	var records []*models.CourseAuditionRecord
	_, err := o.QueryTable("course_audition_record").
		Filter("user_id", userId).
		Filter("teacher_id", teacherId).
		Exclude("status", models.AUDITION_RECORD_STATUS_COMPLETE).
		OrderBy("-last_update_time").All(&records)
	return records, err
}

func QueryStudentUncompletedAuditionRecords(userId int64) ([]*models.CourseAuditionRecord, error) {
	o := orm.NewOrm()
	var records []*models.CourseAuditionRecord
	_, err := o.QueryTable("course_audition_record").
		Filter("user_id", userId).
		Exclude("status", models.AUDITION_RECORD_STATUS_COMPLETE).
		OrderBy("-last_update_time").All(&records)
	return records, err
}

func QueryTeacherUncompletedAuditionRecords(teacherId int64) ([]*models.CourseAuditionRecord, error) {
	o := orm.NewOrm()
	var records []*models.CourseAuditionRecord
	_, err := o.QueryTable("course_audition_record").
		Filter("teacher_id", teacherId).
		Exclude("status", models.AUDITION_RECORD_STATUS_COMPLETE).
		OrderBy("-last_update_time").All(&records)
	return records, err
}

func QueryStudentCompletedAuditionRecords(userId int64) ([]*models.CourseAuditionRecord, error) {
	o := orm.NewOrm()
	var records []*models.CourseAuditionRecord
	_, err := o.QueryTable("course_audition_record").
		Filter("user_id", userId).
		Filter("status", models.AUDITION_RECORD_STATUS_COMPLETE).
		OrderBy("-last_update_time").
		All(&records)
	return records, err
}

func QueryTeacherCompletedAuditionRecords(teacherId int64) ([]*models.CourseAuditionRecord, error) {
	o := orm.NewOrm()
	var records []*models.CourseAuditionRecord
	_, err := o.QueryTable("course_audition_record").
		Filter("teacher_id", teacherId).
		Filter("status", models.AUDITION_RECORD_STATUS_COMPLETE).
		OrderBy("-last_update_time").All(&records)
	return records, err
}

func GetCourseAuditionRecordByUserId(courseId, userId int64) (models.CourseAuditionRecord, error) {
	o := orm.NewOrm()
	var audition models.CourseAuditionRecord
	err := o.QueryTable(new(models.CourseAuditionRecord).TableName()).
		Filter("course_id", courseId).Filter("user_id", userId).
		Exclude("status", models.AUDITION_RECORD_STATUS_COMPLETE).
		One(&audition)
	return audition, err
}
