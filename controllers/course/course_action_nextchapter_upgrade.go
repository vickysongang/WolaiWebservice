// course_action_nextchapter_upgrade
package course

import (
	"errors"

	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
	"WolaiWebservice/utils/leancloud/lcmessage"
)

func HandleCourseActionNextChapterUpgrade(userId, studentId, courseId, chapterId int64) (int64, error) {
	var err error
	_, err = models.ReadUser(userId)
	if err != nil {
		return 2, errors.New("用户信息异常")
	}

	_, err = models.ReadUser(studentId)
	if err != nil {
		return 2, errors.New("用户信息异常")
	}

	var course *models.Course
	if courseId == 0 { //代表试听课，从H5页面跳转过来的
		course = courseService.QueryAuditionCourse()
		if course == nil {
			return 2, errors.New("课程信息异常")
		}
	} else {
		course, err = models.ReadCourse(courseId)
		if err != nil {
			return 2, nil
		}
	}
	if course.Type == models.COURSE_TYPE_DELUXE {
		status, err := HandleDeluxeCourseNextChapterUpgrade(userId, studentId, courseId, chapterId)
		return status, err
	} else if course.Type == models.COURSE_TYPE_AUDITION {
		status, err := HandleAuditionCourseNextChapterUpgrade(userId, studentId, courseId, chapterId)
		return status, err
	}
	return 0, nil
}

func HandleDeluxeCourseNextChapterUpgrade(userId, studentId, courseId, chapterId int64) (int64, error) {
	var err error
	o := orm.NewOrm()

	chapter, err := models.ReadCourseCustomChapter(chapterId)
	if err != nil {
		return 2, errors.New("课程信息异常")
	}

	var purchase models.CoursePurchaseRecord
	err = o.QueryTable("course_purchase_record").
		Filter("course_id", courseId).Filter("user_id", studentId).One(&purchase)
	if err != nil {
		return 2, errors.New("课程信息异常")
	}
	if purchase.TeacherId != userId {
		return 2, errors.New("课程信息异常")
	}

	latestPeriod, err := courseService.QueryLatestCourseChapterPeriod(courseId, studentId)
	if err == nil {
		if chapter.Period != latestPeriod+1 {
			return 2, errors.New("课程信息异常")
		}

		if purchase.PurchaseStatus != models.PURCHASE_RECORD_STATUS_PAID {
			return 2, errors.New("学生尚未完成课程支付")
		}

	} else {
		if latestPeriod != 0 {
			return 2, errors.New("课程信息异常")
		}

		if purchase.AuditionStatus != models.PURCHASE_RECORD_STATUS_PAID {
			return 2, errors.New("学生尚未完成试听支付")
		}
	}

	record := models.CourseChapterToUser{
		CourseId:  courseId,
		ChapterId: chapterId,
		UserId:    studentId,
		TeacherId: userId,
		Period:    chapter.Period,
	}

	_, err = models.CreateCourseChapterToUser(&record)
	if err != nil {
		return 2, errors.New("服务器操作异常")
	}

	go lcmessage.SendCourseChapterCompleteMsg(purchase.Id, chapter.Id)

	chapterCount := purchase.ChapterCount

	recordInfo := map[string]interface{}{
		"audition_status": purchase.AuditionStatus,
	}
	models.UpdateCoursePurchaseRecord(purchase.Id, recordInfo)

	if chapter.Period == 0 {
		recordInfo := map[string]interface{}{
			"audition_status": models.PURCHASE_RECORD_STATUS_COMPLETE,
		}
		models.UpdateCoursePurchaseRecord(purchase.Id, recordInfo)
	} else if chapter.Period == chapterCount-1 {
		recordInfo := map[string]interface{}{
			"purchase_status": models.PURCHASE_RECORD_STATUS_COMPLETE,
		}
		models.UpdateCoursePurchaseRecord(purchase.Id, recordInfo)
	}

	return 0, nil
}

func HandleAuditionCourseNextChapterUpgrade(teacherId, studentId, courseId, chapterId int64) (int64, error) {
	var err error
	o := orm.NewOrm()

	chapter, err := models.ReadCourseCustomChapter(chapterId)
	if err != nil {
		return 2, errors.New("课程信息异常")
	}

	var auditionRecord models.CourseAuditionRecord
	err = o.QueryTable("course_audition_record").
		Filter("course_id", courseId).
		Filter("user_id", studentId).
		Filter("teacher_id", teacherId).
		One(&auditionRecord)
	if err != nil {
		return 2, errors.New("课程信息异常")
	}
	if auditionRecord.TeacherId != teacherId {
		return 2, errors.New("课程信息异常")
	}

	latestPeriod, err := courseService.QueryLatestCourseChapterPeriod(courseId, studentId)
	if err == nil {
		if chapter.Period != latestPeriod+1 {
			return 2, errors.New("课程信息异常")
		}

		if auditionRecord.Status != models.AUDITION_RECORD_STATUS_PAID {
			return 2, errors.New("学生尚未完成课程支付")
		}

	} else {
		if latestPeriod != 0 {
			return 2, errors.New("课程信息异常")
		}
	}

	record := models.CourseChapterToUser{
		CourseId:  courseId,
		ChapterId: chapterId,
		UserId:    studentId,
		TeacherId: teacherId,
		Period:    chapter.Period,
	}

	_, err = models.CreateCourseChapterToUser(&record)
	if err != nil {
		return 2, errors.New("服务器操作异常")
	}

	go lcmessage.SendAuditionCourseChapterCompleteMsg(auditionRecord.Id, chapter.Id)

	recordInfo := map[string]interface{}{
		"Status": models.AUDITION_RECORD_STATUS_COMPLETE,
	}
	_, err = models.UpdateCourseAuditionRecord(auditionRecord.Id, recordInfo)
	if err != nil {
		return 2, errors.New("服务器操作异常")
	}
	return 0, nil
}
