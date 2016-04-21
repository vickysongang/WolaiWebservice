package course

import (
	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
)

func GetCourseDetailTeacherUpgrade(courseId, studentId int64) (int64, *courseDetailTeacher) {
	o := orm.NewOrm()

	var purchaseRecord models.CoursePurchaseRecord
	err := o.QueryTable("course_purchase_record").Filter("user_id", studentId).Filter("course_id", courseId).
		One(&purchaseRecord)
	if err != nil {
		return 2, nil
	}

	course, err := models.ReadCourse(courseId)
	if err != nil {
		return 2, nil
	}

	detail := courseDetailTeacher{
		Course: *course,
	}

	characteristicList, _ := courseService.QueryCourseContentIntros(courseId)
	detail.CharacteristicList = characteristicList

	detail.StudentCount = courseService.GetCourseStudentCount(courseId)

	detail.ChapterCount = purchaseRecord.ChapterCount

	detail.ChapterCompletedPeriod, err = courseService.QueryLatestCourseChapterPeriod(courseId, studentId)
	if err != nil {
		detail.ChapterList, _ = queryCourseCustomChapterStatus(courseId, detail.ChapterCompletedPeriod, studentId, purchaseRecord.TeacherId, true)
	} else {
		detail.ChapterList, _ = queryCourseCustomChapterStatus(courseId, detail.ChapterCompletedPeriod+1, studentId, purchaseRecord.TeacherId, true)
	}

	studentList := make([]*models.User, 0)
	student, _ := models.ReadUser(studentId)
	studentList = append(studentList, student)

	detail.StudentList = studentList

	return 0, &detail
}
