package course

import (
	"time"

	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
)

type courseTeacherListItem struct {
	models.Course
	StudentCount           int64        `json:"studentCount"`
	ChapterCount           int64        `json:"chapterCount"`
	AuditionStatus         string       `json:"auditionStatus"`
	PurchaseStatus         string       `json:"purchaseStatus"`
	ChapterCompletedPeriod int64        `json:"chapterCompletePeriod"`
	LastUpdateTime         string       `json:"lastUpdateTime"`
	StudentInfo            *models.User `json:"studentInfo"`
}

func GetCourseListTeacher(userId, page, count int64) (int64, []*courseTeacherListItem) {
	o := orm.NewOrm()

	items := make([]*courseTeacherListItem, 0)

	var records []*models.CoursePurchaseRecord
	_, err := o.QueryTable("course_purchase_record").Filter("teacher_id", userId).
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

		chapterCompletePeriod, _ := courseService.QueryLatestCourseChapterPeriod(record.CourseId, record.UserId)
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
