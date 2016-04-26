package course

import (
	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
)

func GetCourseDetailStudentUpgrade(userId int64, courseId int64, auditionNum int64) (int64, *courseDetailStudent) {
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
		status, course := GetAuditionCourseDetail(userId, course, auditionNum)
		return status, course
	}
	return 0, nil
}

func GetDeluxeCourseDetail(userId int64, course *models.Course) (int64, *courseDetailStudent) {
	var err error
	o := orm.NewOrm()
	courseId := course.Id
	studentCount := courseService.GetCourseStudentCount(courseId)

	detail := courseDetailStudent{
		Course:       *course,
		StudentCount: studentCount,
	}

	characteristicList, _ := courseService.QueryCourseContentIntros(courseId)
	detail.CharacteristicList = characteristicList

	var purchaseRecord models.CoursePurchaseRecord
	err = o.QueryTable(new(models.CoursePurchaseRecord).TableName()).Filter("user_id", userId).Filter("course_id", courseId).
		One(&purchaseRecord)
	if err != nil && err != orm.ErrNoRows {
		return 2, nil
	}

	purchaseFlag := (err != orm.ErrNoRows) //判断是否购买或者试听
	if !purchaseFlag {
		detail.AuditionStatus = models.PURCHASE_RECORD_STATUS_IDLE
		detail.PurchaseStatus = models.PURCHASE_RECORD_STATUS_IDLE
		detail.TeacherList, _ = queryCourseTeacherList(courseId)
		detail.ChapterList, _ = queryCourseChapterStatus(courseId, 1, true)
		chapterCount := courseService.GetCourseChapterCount(courseId)
		detail.ChapterCount = chapterCount
	} else {
		detail.AuditionStatus = purchaseRecord.AuditionStatus
		detail.PurchaseStatus = purchaseRecord.PurchaseStatus
		detail.TeacherList, _ = queryCourseCurrentTeacher(purchaseRecord.TeacherId)
		detail.ChapterCount = purchaseRecord.ChapterCount
		if purchaseRecord.TeacherId == 0 {
			detail.ChapterList, _ = queryCourseChapterStatus(courseId, 1, true)
		} else {
			detail.ChapterCompletedPeriod, err = courseService.QueryLatestCourseChapterPeriod(courseId, userId)
			if err != nil {
				detail.ChapterList, _ = queryCourseCustomChapterStatus(courseId, 1, userId, purchaseRecord.TeacherId, true)
			} else {
				detail.ChapterList, _ = queryCourseCustomChapterStatus(courseId, detail.ChapterCompletedPeriod+1, userId, purchaseRecord.TeacherId, true)
			}
		}
	}
	auditionCourse := courseService.QueryAuditionCourse()
	if auditionCourse != nil {
		detail.AuditionCourseId = auditionCourse.Id
	}
	return 0, &detail
}

func GetAuditionCourseDetail(userId int64, course *models.Course, auditionNum int64) (int64, *courseDetailStudent) {
	o := orm.NewOrm()
	courseId := course.Id
	studentCount := courseService.GetAuditionCourseStudentCount(courseId)
	detail := courseDetailStudent{
		Course:       *course,
		ChapterCount: 1,
		StudentCount: studentCount,
	}
	characteristicList, _ := courseService.QueryCourseContentIntros(courseId)
	detail.CharacteristicList = characteristicList

	if auditionNum == 0 {
		detail.AuditionStatus = models.PURCHASE_RECORD_STATUS_IDLE
		detail.PurchaseStatus = models.PURCHASE_RECORD_STATUS_IDLE
		detail.TeacherList = make([]*teacherItem, 0)
		detail.ChapterList, _ = queryCourseChapterStatus(courseId, 1, true)
	} else {
		var auditionRecord models.CourseAuditionRecord
		err := o.QueryTable(new(models.CourseAuditionRecord).TableName()).
			Filter("course_id", courseId).Filter("user_id", userId).Filter("audition_num", auditionNum).
			One(&auditionRecord)
		if err != nil && err != orm.ErrNoRows {
			return 2, nil
		}
		detail.AuditionStatus = auditionRecord.Status
		detail.PurchaseStatus = auditionRecord.Status
		detail.TeacherList, _ = queryCourseCurrentTeacher(auditionRecord.TeacherId)
		if auditionRecord.TeacherId == 0 {
			detail.ChapterList, _ = queryCourseChapterStatus(courseId, 1, true)
		} else {
			detail.ChapterCompletedPeriod, err = courseService.QueryLatestCourseChapterPeriod(courseId, userId)
			if err != nil {
				detail.ChapterList, _ = queryCourseCustomChapterStatus(courseId, 1, userId, auditionRecord.TeacherId, true)
			} else {
				detail.ChapterList, _ = queryCourseCustomChapterStatus(courseId, detail.ChapterCompletedPeriod+1, userId, auditionRecord.TeacherId, true)
			}
		}
	}
	return 0, &detail
}
