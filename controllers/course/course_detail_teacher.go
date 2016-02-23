package course

import (
	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
)

type courseDetailTeacher struct {
	models.Course
	StudentCount           int64                       `json:"studentCount"`
	ChapterCount           int64                       `json:"chapterCount"`
	ChapterCompletedPeriod int64                       `json:"chapterCompletePeriod"`
	CharacteristicList     []models.CourseContentIntro `json:"characteristicList"`
	ChapterList            []*courseChapterStatus      `json:"chapterList"`
	StudentList            []*models.User              `json:"studentList"`
}

func GetCourseDetailTeacher(courseId, studentId int64) (int64, *courseDetailTeacher) {
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

	characteristicList, _ := queryCourseContentIntros(courseId)
	detail.CharacteristicList = characteristicList

	detail.StudentCount = queryCourseStudentCount(courseId)
	chapterCount, _ := o.QueryTable("course_chapter").Filter("course_id", courseId).Count()
	detail.ChapterCount = chapterCount - 1
	detail.ChapterCompletedPeriod, err = queryLatestCourseChapterPeriod(courseId, studentId)
	if err != nil {
		detail.ChapterList, _ = queryCourseChapterStatus(courseId, detail.ChapterCompletedPeriod)
	} else {
		detail.ChapterList, _ = queryCourseChapterStatus(courseId, detail.ChapterCompletedPeriod+1)
	}

	studentList := make([]*models.User, 0)
	student, _ := models.ReadUser(studentId)
	studentList = append(studentList, student)

	detail.StudentList = studentList

	return 0, &detail
}
