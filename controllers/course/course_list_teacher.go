package course

import (
	"time"

	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
)

func GetCourseListTeacher(userId, page, count int64) (int64, []*courseTeacherListItem) {
	items := make([]*courseTeacherListItem, 0)

	records, err := courseService.QueryTeacherCoursePurchaseRecords(userId, page, count, false)
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

		chapterCompletePeriod, _ := courseService.GetLatestCompleteChapterPeriod(record.CourseId, record.UserId, record.Id)
		student, err := models.ReadUser(record.UserId)
		if err != nil {
			continue
		}

		item := courseTeacherListItem{
			Course:                 *course,
			StudentCount:           studentCount,
			ChapterCount:           chapterCount - 1,
			AuditionStatus:         record.AuditionStatus,
			PurchaseStatus:         record.PurchaseStatus,
			ChapterCompletedPeriod: chapterCompletePeriod,
			LastUpdateTime:         record.LastUpdateTime.Format(time.RFC3339),
			StudentInfo:            student,
		}

		items = append(items, &item)
	}

	return 0, items
}
