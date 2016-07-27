// course_session
package course

import (
	"WolaiWebservice/models"
	"errors"

	courseService "WolaiWebservice/service/course"
)

func GetCourseListStudentOfConversation(userId, teacherId, page, count int64) (int64, []*courseStudentListItem, error) {
	items := make([]*courseStudentListItem, 0)

	teacher, err := models.ReadUser(teacherId)
	if err != nil {
		return 2, items, errors.New("对方用户不存在")
	}
	if teacher.AccessRight != models.USER_ACCESSRIGHT_TEACHER {
		return 2, items, errors.New("对方不是导师，没有可以选择的课程哦")
	}

	courseCount := courseService.GetConversationCourseCount(userId, teacherId)
	auditionCourseCount := courseService.GetConversationAuditonCourseCount(userId, teacherId)

	if courseCount == 0 && auditionCourseCount == 0 {
		return 2, items, errors.New("还没有选择该导师的课程，可以先去导师的个人主页看看哦")
	}

	if page == 0 {
		auditionUncompleteRecords, _ := courseService.QueryUncompletedAuditionRecords(userId, teacherId)
		for _, record := range auditionUncompleteRecords {
			item := assignStudentAuditionCourseInfo(record.CourseId,
				record.UserId,
				record.Status,
				record.Id,
				record.TeacherId)
			items = append(items, item)
		}
	}

	records, err := courseService.QueryCoursePurchaseRecords(userId, teacherId, page, count)
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
