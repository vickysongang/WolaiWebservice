package course

import (
	"WolaiWebservice/models"
	"fmt"

	courseService "WolaiWebservice/service/course"

	"github.com/astaxie/beego/orm"
)

func GetCourseDetailTeacherUpgrade(courseId, studentId, teacherId, auditionNum int64) (int64, *courseDetailTeacher) {
	var err error
	var course *models.Course
	if courseId == 0 { //代表试听课，从H5页面跳转过来的
		course = courseService.QueryAuditionCourse()
		if course == nil {
			return 2, nil
		}
	} else {
		course, err = models.ReadCourse(courseId)
		if err != nil {
			return 2, nil
		}
	}
	if course.Type == models.COURSE_TYPE_DELUXE {
		status, course := GetDeluxeCourseDetailTeacher(studentId, teacherId, course)
		return status, course
	} else if course.Type == models.COURSE_TYPE_AUDITION {
		status, course := GetAuditionCourseDetailTeacher(studentId, teacherId, course, auditionNum)
		return status, course
	}
	return 0, nil
}

func GetDeluxeCourseDetailTeacher(studentId, teacherId int64, course *models.Course) (int64, *courseDetailTeacher) {
	o := orm.NewOrm()
	courseId := course.Id
	var purchaseRecord models.CoursePurchaseRecord
	err := o.QueryTable("course_purchase_record").Filter("user_id", studentId).Filter("course_id", courseId).Filter("teacher_id", teacherId).
		One(&purchaseRecord)
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
		detail.ChapterList, _ = queryCourseCustomChapterStatus(courseId, 1, studentId, purchaseRecord.TeacherId, true)
	} else {
		detail.ChapterList, _ = queryCourseCustomChapterStatus(courseId, detail.ChapterCompletedPeriod+1, studentId, purchaseRecord.TeacherId, true)
	}

	studentList := make([]*models.User, 0)
	student, _ := models.ReadUser(studentId)
	studentList = append(studentList, student)

	detail.StudentList = studentList

	return 0, &detail
}

func GetAuditionCourseDetailTeacher(studentId, teacherId int64, course *models.Course, auditionNum int64) (int64, *courseDetailTeacher) {
	o := orm.NewOrm()
	courseId := course.Id
	var auditionoRecord models.CourseAuditionRecord
	err := o.QueryTable("course_audition_record").
		Filter("user_id", studentId).
		Filter("course_id", courseId).
		Filter("teacher_id", teacherId).
		Filter("audition_num", auditionNum).
		One(&auditionoRecord)
	if err != nil {
		fmt.Println("ss", err.Error())
		return 2, nil
	}

	detail := courseDetailTeacher{
		Course: *course,
	}

	characteristicList, _ := courseService.QueryCourseContentIntros(courseId)
	detail.CharacteristicList = characteristicList

	detail.StudentCount = courseService.GetCourseStudentCount(courseId)

	detail.ChapterCount = 1

	detail.ChapterCompletedPeriod, err = courseService.QueryLatestCourseChapterPeriod(courseId, studentId)
	if err != nil {
		detail.ChapterList, _ = queryCourseCustomChapterStatus(courseId, 1, studentId, auditionoRecord.TeacherId, true)
	} else {
		detail.ChapterList, _ = queryCourseCustomChapterStatus(courseId, detail.ChapterCompletedPeriod+1, studentId, auditionoRecord.TeacherId, true)
	}

	studentList := make([]*models.User, 0)
	student, _ := models.ReadUser(studentId)
	studentList = append(studentList, student)

	detail.StudentList = studentList

	return 0, &detail
}
