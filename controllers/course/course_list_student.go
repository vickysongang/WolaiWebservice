package course

import (
	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
)

func GetCourseListStudent(userId, page, count int64) (int64, []*courseStudentListItem) {
	items := make([]*courseStudentListItem, 0)

	records, err := courseService.QueryStudentCoursePurchaseRecords(userId, page, count, false)
	if err != nil {
		return 0, items
	}

	for _, record := range records {
		course, err := models.ReadCourse(record.CourseId)
		if err != nil {
			continue
		}

		studentCount := courseService.GetCourseStudentCount(record.CourseId)
		chapterCount := courseService.GetCourseChapterCount(record.CourseId)

		chapterCompletePeriod, _ := courseService.GetLatestCompleteChapterPeriod(record.CourseId, userId, record.Id)

		item := courseStudentListItem{
			Course:                 *course,
			StudentCount:           studentCount,
			ChapterCount:           chapterCount - 1,
			AuditionStatus:         record.AuditionStatus,
			PurchaseStatus:         record.PurchaseStatus,
			ChapterCompletedPeriod: chapterCompletePeriod,
		}

		items = append(items, &item)
	}

	return 0, items
}
