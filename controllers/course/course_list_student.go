package course

import (
	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
)

type courseStudentListItem struct {
	models.Course
	StudentCount           int64  `json:"studentCount"`
	ChapterCount           int64  `json:"chapterCount"`
	AuditionStatus         string `json:"auditionStatus"`
	PurchaseStatus         string `json:"purchaseStatus"`
	ChapterCompletedPeriod int64  `json:"chapterCompletePeriod"`
}

func GetCourseListStudent(userId, page, count int64) (int64, []*courseStudentListItem) {
	o := orm.NewOrm()
	var err error

	items := make([]*courseStudentListItem, 0)
	if page == 0 {
		var auditionUncompleteRecords []*models.CourseAuditionRecord

		_, err = o.QueryTable("course_audition_record").Filter("user_id", userId).
			Exclude("status", models.AUDITION_RECORD_STATUS_COMPLETE).
			OrderBy("-last_update_time").All(&auditionUncompleteRecords)

		for _, auditionRecord := range auditionUncompleteRecords {
			item := assignAuditionCourseInfo(auditionRecord.CourseId, userId, auditionRecord.Status)
			items = append(items, item)
		}
	}

	var records []*models.CoursePurchaseRecord
	_, err = o.QueryTable("course_purchase_record").Filter("user_id", userId).
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

		chapterCompletePeriod, _ := courseService.QueryLatestCourseChapterPeriod(record.CourseId, userId)

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

	recordsLen := int64(len(records))
	if recordsLen < count {
		var auditionCompleteRecords []*models.CourseAuditionRecord
		o.QueryTable("course_audition_record").Filter("user_id", userId).
			Filter("status", models.AUDITION_RECORD_STATUS_COMPLETE).
			OrderBy("-last_update_time").All(&auditionCompleteRecords)

		for _, auditionRecord := range auditionCompleteRecords {
			item := assignAuditionCourseInfo(auditionRecord.CourseId, userId, auditionRecord.Status)
			items = append(items, item)
		}
	}

	return 0, items
}

func assignAuditionCourseInfo(courseId, userId int64, status string) *courseStudentListItem {
	course, err := models.ReadCourse(courseId)
	if err != nil {
		return nil
	}

	studentCount := courseService.GetCourseStudentCount(courseId)
	chapterCount := int64(1)

	chapterCompletePeriod, _ := courseService.QueryLatestCourseChapterPeriod(courseId, userId)

	item := courseStudentListItem{
		Course:                 *course,
		StudentCount:           studentCount,
		ChapterCount:           chapterCount,
		AuditionStatus:         status,
		PurchaseStatus:         status,
		ChapterCompletedPeriod: chapterCompletePeriod,
	}
	return &item
}
