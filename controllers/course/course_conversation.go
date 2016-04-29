// course_session
package course

import (
	"WolaiWebservice/models"
	"errors"

	courseService "WolaiWebservice/service/course"

	"github.com/astaxie/beego/orm"
)

func QueryCourseCountOfConversation(studentId, teacherId int64) int64 {
	o := orm.NewOrm()
	count, _ := o.QueryTable("course_purchase_record").
		Filter("user_id", studentId).
		Filter("teacher_id", teacherId).Count()
	return count
}

func QueryAuditonCourseCountOfConversation(studentId, teacherId int64) int64 {
	o := orm.NewOrm()
	count, _ := o.QueryTable("course_audition_record").
		Exclude("status", models.AUDITION_RECORD_STATUS_COMPLETE).
		Filter("user_id", studentId).
		Filter("teacher_id", teacherId).Count()
	return count
}

func GetCourseListStudentOfConversation(userId, teacherId, page, count int64) (int64, []*courseStudentListItem, error) {
	o := orm.NewOrm()

	items := make([]*courseStudentListItem, 0)

	teacher, err := models.ReadUser(teacherId)
	if err != nil {
		return 2, items, errors.New("对方用户不存在")
	}
	if teacher.AccessRight != models.USER_ACCESSRIGHT_TEACHER {
		return 2, items, errors.New("对方不是导师，没有可以选择的课程哦")
	}

	courseCount := QueryCourseCountOfConversation(userId, teacherId)
	auditionCourseCount := QueryAuditonCourseCountOfConversation(userId, teacherId)
	if courseCount == 0 && auditionCourseCount == 0 {
		return 2, items, errors.New("还没有选择该导师的课程，可以先去导师的个人主页看看哦")
	}

	if page == 0 {
		var auditionUncompleteRecords []*models.CourseAuditionRecord
		_, err = o.QueryTable("course_audition_record").Filter("user_id", userId).Filter("teacher_id", teacherId).
			Exclude("status", models.AUDITION_RECORD_STATUS_COMPLETE).
			OrderBy("-last_update_time").All(&auditionUncompleteRecords)

		for _, auditionRecord := range auditionUncompleteRecords {
			item := assignStudentAuditionCourseInfo(auditionRecord.CourseId,
				auditionRecord.UserId,
				auditionRecord.Status,
				auditionRecord.Id,
				auditionRecord.TeacherId)
			items = append(items, item)
		}
	}

	var records []*models.CoursePurchaseRecord
	_, err = o.QueryTable("course_purchase_record").
		Filter("user_id", userId).
		Filter("teacher_id", teacherId).
		OrderBy("-last_update_time").Offset(page * count).Limit(count).All(&records)
	if err != nil {
		return 2, items, err
	}

	for _, record := range records {
		course, err := models.ReadCourse(record.CourseId)
		if err != nil {
			continue
		}

		studentCount := courseService.GetCourseStudentCount(record.CourseId)
		chapterCount := record.ChapterCount

		chapterCompletePeriod, _ := courseService.GetLatestCompleteChapterPeriod(record.CourseId, userId, record.Id)

		item := courseStudentListItem{
			Course:                 *course,
			StudentCount:           studentCount,
			ChapterCount:           chapterCount,
			AuditionStatus:         record.AuditionStatus,
			PurchaseStatus:         record.PurchaseStatus,
			ChapterCompletedPeriod: chapterCompletePeriod,
		}

		items = append(items, &item)
	}

	return 0, items, nil
}
