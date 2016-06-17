package course

import (
	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
)

func GetCourseDetailStudentUpgrade(userId int64, courseId int64, recordId int64) (int64, *courseDetailStudent) {
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
		status, course := GetDeluxeCourseDetail(userId, course)
		return status, course
	} else if course.Type == models.COURSE_TYPE_AUDITION {
		status, course := GetAuditionCourseDetail(userId, course, recordId)
		return status, course
	}
	return 0, nil
}

func GetDeluxeCourseDetail(userId int64, course *models.Course) (int64, *courseDetailStudent) {
	var err error
	courseId := course.Id
	studentCount := courseService.GetCourseStudentCount(courseId)

	detail := courseDetailStudent{
		Course:       *course,
		StudentCount: studentCount,
	}

	characteristicList, _ := courseService.QueryCourseContentIntros(courseId)
	detail.CharacteristicList = characteristicList

	purchaseRecord, err := courseService.GetCoursePurchaseRecordByUserId(courseId, userId)
	if err != nil && err != orm.ErrNoRows {
		return 2, nil
	}
	detail.RecordId = purchaseRecord.Id
	purchaseFlag := (err != orm.ErrNoRows) //判断是否购买或者试听
	if !purchaseFlag {
		detail.PurchaseStatus = models.PURCHASE_RECORD_STATUS_IDLE
		detail.TeacherList, _ = queryCourseTeacherList(courseId)
		detail.ChapterList, _ = queryCourseChapterStatus(courseId, 1, true)
		chapterCount := courseService.GetCourseChapterCount(courseId)
		detail.ChapterCount = chapterCount
	} else {
		detail.PurchaseStatus = purchaseRecord.PurchaseStatus
		detail.TeacherList, _ = queryCourseCurrentTeacher(purchaseRecord.TeacherId)
		detail.ChapterCount = purchaseRecord.ChapterCount
		if purchaseRecord.TeacherId == 0 {
			detail.ChapterList, _ = queryCourseChapterStatus(courseId, 1, true)
		} else {
			detail.ChapterCompletedPeriod, err = courseService.GetLatestCompleteChapterPeriod(courseId, userId, purchaseRecord.Id)
			if err != nil {
				detail.ChapterList, _ = queryCourseCustomChapterStatus(courseId,
					1,
					userId,
					purchaseRecord.TeacherId,
					purchaseRecord.Id,
					models.COURSE_TYPE_DELUXE,
					true)
			} else {
				detail.ChapterList, _ = queryCourseCustomChapterStatus(courseId,
					detail.ChapterCompletedPeriod+1,
					userId,
					purchaseRecord.TeacherId,
					purchaseRecord.Id,
					models.COURSE_TYPE_DELUXE,
					true)
			}
		}
	}
	auditionCourse := courseService.QueryAuditionCourse()
	if auditionCourse != nil {
		detail.AuditionCourseId = auditionCourse.Id
	}
	return 0, &detail
}

func GetAuditionCourseDetail(userId int64, course *models.Course, recordId int64) (int64, *courseDetailStudent) {
	courseId := course.Id
	studentCount := courseService.GetAuditionCourseStudentCount(courseId)
	detail := courseDetailStudent{
		Course:       *course,
		ChapterCount: 1,
		StudentCount: studentCount,
	}
	characteristicList, _ := courseService.QueryCourseContentIntros(courseId)
	detail.CharacteristicList = characteristicList

	if recordId == 0 {
		auditionRecord, _ := courseService.GetCourseAuditionRecordByUserId(courseId, userId)
		if auditionRecord.Id != 0 {
			detail.PurchaseStatus = auditionRecord.Status
			detail.RecordId = auditionRecord.Id
		} else {
			detail.PurchaseStatus = models.PURCHASE_RECORD_STATUS_IDLE
		}
		detail.TeacherList = make([]*teacherItem, 0)
		detail.ChapterList, _ = queryCourseChapterStatus(courseId, 1, true)
	} else {
		auditionRecord, err := models.ReadCourseAuditionRecord(recordId)
		if err != nil {
			return 2, nil
		}
		detail.RecordId = auditionRecord.Id
		detail.PurchaseStatus = auditionRecord.Status
		detail.TeacherList, _ = queryCourseCurrentTeacher(auditionRecord.TeacherId)
		if auditionRecord.TeacherId == 0 {
			detail.ChapterList, _ = queryCourseChapterStatus(courseId, 1, true)
		} else {
			detail.ChapterCompletedPeriod, err = courseService.GetLatestCompleteChapterPeriod(courseId, userId, auditionRecord.Id)
			if err != nil {
				detail.ChapterList, _ = queryCourseCustomChapterStatus(courseId,
					1,
					userId,
					auditionRecord.TeacherId,
					auditionRecord.Id,
					models.COURSE_TYPE_AUDITION,
					true)
			} else {
				detail.ChapterList, _ = queryCourseCustomChapterStatus(courseId,
					detail.ChapterCompletedPeriod+1,
					userId,
					auditionRecord.TeacherId,
					auditionRecord.Id,
					models.COURSE_TYPE_AUDITION,
					true)
			}
		}
	}
	return 0, &detail
}
