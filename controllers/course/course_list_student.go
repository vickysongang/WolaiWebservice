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

	items := make([]*courseStudentListItem, 0)

	var records []*models.CoursePurchaseRecord
	_, err := o.QueryTable("course_purchase_record").Filter("user_id", userId).
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
		chapterCount := courseService.GetCourseChapterCount(record.CourseId)

		chapterCompletePeriod, _ := courseService.QueryLatestCourseChapterPeriod(record.CourseId, userId)

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
