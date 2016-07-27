package course

import (
	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
)

func GetCourseListStudentUpgrade(userId, page, count int64) (int64, []*courseStudentListItem) {
	var err error

	items := make([]*courseStudentListItem, 0)
	if page == 0 {
		auditionUncompleteRecords, _ := courseService.QueryStudentUncompletedAuditionRecords(userId)
		for _, record := range auditionUncompleteRecords {
			item := assignStudentAuditionCourseInfo(record.CourseId, userId, record.Status, record.Id, record.TeacherId)
			items = append(items, item)
		}
	}

	records, err := courseService.QueryStudentCoursePurchaseRecords(userId, page, count, true)
	if err != nil {
		return 0, items
	}
	totalCount := courseService.GetStudentCoursePurchaseRecordCount(userId)

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
			PurchaseStatus:         record.PurchaseStatus,
			ChapterCompletedPeriod: chapterCompletePeriod,
			TeacherId:              record.TeacherId,
			RecordId:               record.Id,
		}

		items = append(items, &item)
	}

	if page == totalCount/count {
		auditionCompleteRecords, _ := courseService.QueryStudentCompletedAuditionRecords(userId)
		for _, record := range auditionCompleteRecords {
			item := assignStudentAuditionCourseInfo(record.CourseId, userId, record.Status, record.Id, record.TeacherId)
			items = append(items, item)
		}
	}

	return 0, items
}

func assignStudentAuditionCourseInfo(courseId, userId int64, status string, recordId, teacherId int64) *courseStudentListItem {
	course, err := models.ReadCourse(courseId)
	if err != nil {
		return nil
	}

	studentCount := courseService.GetCourseStudentCount(courseId)
	chapterCount := int64(1)

	chapterCompletePeriod, _ := courseService.GetLatestCompleteChapterPeriod(courseId, userId, recordId)

	item := courseStudentListItem{
		Course:                 *course,
		StudentCount:           studentCount,
		ChapterCount:           chapterCount,
		PurchaseStatus:         status,
		ChapterCompletedPeriod: chapterCompletePeriod,
		TeacherId:              teacherId,
		RecordId:               recordId,
	}
	return &item
}
