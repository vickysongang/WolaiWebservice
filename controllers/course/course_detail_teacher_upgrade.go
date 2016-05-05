package course

import (
	"WolaiWebservice/models"

	courseService "WolaiWebservice/service/course"

	"github.com/astaxie/beego/orm"
)

func GetCourseDetailTeacherUpgrade(courseId, studentId, teacherId, recordId int64) (int64, *courseDetailTeacher) {
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
		status, course := GetAuditionCourseDetailTeacher(studentId, teacherId, course, recordId)
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
		Course:         *course,
		RecordId:       purchaseRecord.Id,
		PurchaseStatus: purchaseRecord.PurchaseStatus,
	}

	characteristicList, _ := courseService.QueryCourseContentIntros(courseId)
	detail.CharacteristicList = characteristicList

	detail.StudentCount = courseService.GetCourseStudentCount(courseId)

	detail.ChapterCount = purchaseRecord.ChapterCount

	detail.ChapterCompletedPeriod, err = courseService.GetLatestCompleteChapterPeriod(courseId, studentId, purchaseRecord.Id)
	if err != nil {
		detail.ChapterList, _ = queryCourseCustomChapterStatus(courseId,
			1,
			studentId,
			purchaseRecord.TeacherId,
			purchaseRecord.Id,
			models.COURSE_TYPE_DELUXE,
			true)
	} else {
		detail.ChapterList, _ = queryCourseCustomChapterStatus(courseId,
			detail.ChapterCompletedPeriod+1,
			studentId,
			purchaseRecord.TeacherId,
			purchaseRecord.Id,
			models.COURSE_TYPE_DELUXE,
			true)
	}

	studentList := make([]*models.User, 0)
	student, _ := models.ReadUser(studentId)
	studentList = append(studentList, student)

	detail.StudentList = studentList

	return 0, &detail
}

func GetAuditionCourseDetailTeacher(studentId, teacherId int64, course *models.Course, recordId int64) (int64, *courseDetailTeacher) {
	courseId := course.Id
	auditionoRecord, err := models.ReadCourseAuditionRecord(recordId)
	if err != nil {
		return 2, nil
	}

	detail := courseDetailTeacher{
		Course: *course,
	}

	characteristicList, _ := courseService.QueryCourseContentIntros(courseId)
	detail.CharacteristicList = characteristicList

	detail.StudentCount = courseService.GetCourseStudentCount(courseId)

	detail.ChapterCount = 1
	detail.RecordId = recordId
	detail.PurchaseStatus = auditionoRecord.Status
	detail.ChapterCompletedPeriod, err = courseService.GetLatestCompleteChapterPeriod(courseId, studentId, auditionoRecord.Id)
	if err != nil {
		detail.ChapterList, _ = queryCourseCustomChapterStatus(courseId,
			1,
			studentId,
			auditionoRecord.TeacherId,
			auditionoRecord.Id,
			models.COURSE_TYPE_AUDITION,
			true)
	} else {
		detail.ChapterList, _ = queryCourseCustomChapterStatus(courseId,
			detail.ChapterCompletedPeriod+1,
			studentId,
			auditionoRecord.TeacherId,
			auditionoRecord.Id,
			models.COURSE_TYPE_AUDITION,
			true)
	}

	studentList := make([]*models.User, 0)
	student, _ := models.ReadUser(studentId)
	studentList = append(studentList, student)

	detail.StudentList = studentList

	return 0, &detail
}
