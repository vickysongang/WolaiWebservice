package course

import (
	"time"

	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
)

func GetCourseListTeacherUpgrade(teacherId, page, count int64) (int64, []*courseTeacherListItem) {
	var err error

	items := make([]*courseTeacherListItem, 0)

	if page == 0 {
		auditionUncompleteRecords, _ := courseService.QueryTeacherUncompletedAuditionRecords(teacherId)
		for _, record := range auditionUncompleteRecords {
			item := assignTeacherAuditionCourseInfo(record.CourseId,
				record.UserId,
				record.Status,
				record.LastUpdateTime,
				record.Id)
			items = append(items, item)
		}
	}

	records, err := courseService.QueryTeacherCoursePurchaseRecords(teacherId, page, count, true)
	if err != nil {
		return 0, items
	}
	totalCount := courseService.GetTeacherCoursePurchaseRecordCount(teacherId)
	for _, record := range records {
		course, err := models.ReadCourse(record.CourseId)
		if err != nil {
			continue
		}

		studentCount := courseService.GetCourseStudentCount(record.CourseId)
		chapterCount := record.ChapterCount

		chapterCompletePeriod, _ := courseService.GetLatestCompleteChapterPeriod(record.CourseId, record.UserId, record.Id)
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
			RecordId:               record.Id,
		}

		items = append(items, &item)
	}
	if page == totalCount/count {
		auditionCompleteRecords, _ := courseService.QueryTeacherCompletedAuditionRecords(teacherId)
		for _, record := range auditionCompleteRecords {
			item := assignTeacherAuditionCourseInfo(record.CourseId,
				record.UserId,
				record.Status,
				record.LastUpdateTime,
				record.Id)
			items = append(items, item)
		}
	}
	return 0, items
}

func assignTeacherAuditionCourseInfo(courseId, userId int64, status string, lastUpdateTime time.Time, recordId int64) *courseTeacherListItem {
	course, err := models.ReadCourse(courseId)
	if err != nil {
		return nil
	}

	studentCount := courseService.GetCourseStudentCount(courseId)
	chapterCount := int64(1)

	chapterCompletePeriod, _ := courseService.GetLatestCompleteChapterPeriod(courseId, userId, recordId)

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
		RecordId:               recordId,
	}
	return &item
}
