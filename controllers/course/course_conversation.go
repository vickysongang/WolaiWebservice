// course_session
package course

import (
	"WolaiWebservice/models"

	"github.com/astaxie/beego/orm"
)

func QueryCourseCountOfConversation(studentId, teacherId int64) int64 {
	o := orm.NewOrm()
	count, _ := o.QueryTable("course_purchase_record").
		Filter("user_id", studentId).
		Filter("teacher_id", teacherId).Count()
	return count
}

func GetCourseListStudentOfConversation(userId, teacherId, page, count int64) (int64, []*courseStudentListItem) {
	o := orm.NewOrm()

	items := make([]*courseStudentListItem, 0)

	var records []*models.CoursePurchaseRecord
	_, err := o.QueryTable("course_purchase_record").
		Filter("user_id", userId).
		Filter("teacher_id", teacherId).
		OrderBy("-last_update_time").Offset(page * count).Limit(count).All(&records)
	if err != nil {
		return 0, items
	}

	for _, record := range records {
		course, err := models.ReadCourse(record.CourseId)
		if err != nil {
			continue
		}

		studentCount := queryCourseStudentCount(record.CourseId)
		chapterCount := queryCourseChapterCount(record.CourseId)

		chapterCompletePeriod, _ := queryLatestCourseChapterPeriod(record.CourseId, userId)

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
