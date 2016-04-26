package course

import (
	"time"

	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
)

func GetCourseListTeacherUpgrade(teacherId, page, count int64) (int64, []*courseTeacherListItem) {
	var err error
	o := orm.NewOrm()

	items := make([]*courseTeacherListItem, 0)

	if page == 0 {
		var auditionUncompleteRecords []*models.CourseAuditionRecord
		_, err = o.QueryTable("course_audition_record").Filter("teacher_id", teacherId).
			Exclude("status", models.AUDITION_RECORD_STATUS_COMPLETE).
			OrderBy("-last_update_time").All(&auditionUncompleteRecords)

		for _, auditionRecord := range auditionUncompleteRecords {
			item := assignTeacherAuditionCourseInfo(auditionRecord.CourseId,
				auditionRecord.UserId,
				auditionRecord.Status,
				auditionRecord.LastUpdateTime,
				auditionRecord.AuditionNum)
			items = append(items, item)
		}
	}

	var records []*models.CoursePurchaseRecord
	_, err = o.QueryTable("course_purchase_record").Filter("teacher_id", teacherId).
		OrderBy("-last_update_time").Offset(page * count).Limit(count).All(&records)
	if err != nil {
		return 0, items
	}

	for _, record := range records {
		course, err := models.ReadCourse(record.CourseId)
		if err != nil {
			continue
		}

		studentCount := courseService.GetCourseStudentCount(record.CourseId)
		chapterCount := record.ChapterCount

		chapterCompletePeriod, _ := courseService.QueryLatestCourseChapterPeriod(record.CourseId, record.UserId)
		student, err := models.ReadUser(record.UserId)
		if err != nil {
			continue
		}

		item := courseTeacherListItem{
			Course:                 *course,
			StudentCount:           studentCount,
			ChapterCount:           chapterCount,
			PurchaseStatus:         record.PurchaseStatus,
			ChapterCompletedPeriod: chapterCompletePeriod,
			LastUpdateTime:         record.LastUpdateTime.Format(time.RFC3339),
			StudentInfo:            student,
			AuditionNum:            0,
		}

		items = append(items, &item)
	}

	recordsLen := int64(len(records))
	if recordsLen < count {
		var auditionCompleteRecords []*models.CourseAuditionRecord
		o.QueryTable("course_audition_record").Filter("teacher_id", teacherId).
			Filter("status", models.AUDITION_RECORD_STATUS_COMPLETE).
			OrderBy("-last_update_time").All(&auditionCompleteRecords)

		for _, auditionRecord := range auditionCompleteRecords {
			item := assignTeacherAuditionCourseInfo(auditionRecord.CourseId,
				auditionRecord.UserId,
				auditionRecord.Status,
				auditionRecord.LastUpdateTime,
				auditionRecord.AuditionNum)
			items = append(items, item)
		}
	}
	return 0, items
}

func assignTeacherAuditionCourseInfo(courseId, userId int64, status string, lastUpdateTime time.Time, auditionNum int64) *courseTeacherListItem {
	course, err := models.ReadCourse(courseId)
	if err != nil {
		return nil
	}

	studentCount := courseService.GetAuditionCourseStudentCount(courseId)
	chapterCount := int64(1)

	chapterCompletePeriod, _ := courseService.QueryLatestCourseChapterPeriod(courseId, userId)

	student, err := models.ReadUser(userId)
	if err != nil {
		return nil
	}

	item := courseTeacherListItem{
		Course:                 *course,
		StudentCount:           studentCount,
		ChapterCount:           chapterCount,
		PurchaseStatus:         status,
		ChapterCompletedPeriod: chapterCompletePeriod,
		LastUpdateTime:         lastUpdateTime.Format(time.RFC3339),
		StudentInfo:            student,
		AuditionNum:            auditionNum,
	}
	return &item
}
