// course_action_nextchapter_upgrade
package course

import (
	"errors"

	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
	"WolaiWebservice/utils/leancloud/lcmessage"
)

func HandleCourseActionNextChapterUpgrade(userId, chapterId, recordId int64) (int64, error) {
	var err error
	_, err = models.ReadUser(userId)
	if err != nil {
		return 2, ErrUserAbnormal
	}
	chapter, err := models.ReadCourseCustomChapter(chapterId)
	if err != nil {
		return 2, ErrChapterAbnormal
	}

	course, err := models.ReadCourse(chapter.CourseId)
	if err != nil {
		return 2, nil
	}
	if course.Type == models.COURSE_TYPE_DELUXE {
		status, err := HandleDeluxeCourseNextChapterUpgrade(userId, chapterId, recordId)
		return status, err
	} else if course.Type == models.COURSE_TYPE_AUDITION {
		status, err := HandleAuditionCourseNextChapterUpgrade(userId, chapterId, recordId)
		return status, err
	}
	return 0, nil
}

func HandleDeluxeCourseNextChapterUpgrade(userId, chapterId, recordId int64) (int64, error) {
	var err error

	chapter, err := models.ReadCourseCustomChapter(chapterId)
	if err != nil {
		return 2, ErrChapterAbnormal
	}

	purchase, err := models.ReadCoursePurchaseRecord(recordId)
	if err != nil {
		return 2, ErrCourseAbnormal
	}
	if purchase.TeacherId != userId {
		return 2, errors.New("课程导师信息异常")
	}
	courseId := purchase.CourseId
	studentId := purchase.UserId
	latestPeriod, _ := courseService.GetLatestCompleteChapterPeriod(courseId, studentId, purchase.Id)

	if chapter.Period != latestPeriod+1 {
		return 2, errors.New("课程课时号信息异常")
	}
	if purchase.PurchaseStatus == models.PURCHASE_RECORD_STATUS_COMPLETE {
		return 2, errors.New("学生还未购买该课时")
	} else if purchase.PurchaseStatus != models.PURCHASE_RECORD_STATUS_PAID {
		return 2, errors.New("学生尚未完成课程支付")
	}

	record := models.CourseChapterToUser{
		CourseId:  courseId,
		ChapterId: chapterId,
		UserId:    studentId,
		TeacherId: userId,
		Period:    chapter.Period,
		RecordId:  purchase.Id,
	}

	_, err = models.CreateCourseChapterToUser(&record)
	if err != nil {
		return 2, ErrServerAbnormal
	}

	go lcmessage.SendCourseChapterCompleteMsg(purchase.Id, chapter.Id)

	if chapter.Period == purchase.ChapterCount {
		recordInfo := map[string]interface{}{
			"AuditionStatus": models.PURCHASE_RECORD_STATUS_COMPLETE,
			"PurchaseStatus": models.PURCHASE_RECORD_STATUS_COMPLETE,
		}
		models.UpdateCoursePurchaseRecord(purchase.Id, recordInfo)
	}

	return 0, nil
}

func HandleAuditionCourseNextChapterUpgrade(teacherId, chapterId, recordId int64) (int64, error) {
	var err error
	chapter, err := models.ReadCourseCustomChapter(chapterId)
	if err != nil {
		return 2, ErrCourseAbnormal
	}

	auditionRecord, err := models.ReadCourseAuditionRecord(recordId)
	if err != nil {
		return 2, ErrCourseAbnormal
	}
	if auditionRecord.TeacherId != teacherId {
		return 2, ErrCourseAbnormal
	}
	courseId := auditionRecord.CourseId
	studentId := auditionRecord.UserId

	if auditionRecord.Status != models.AUDITION_RECORD_STATUS_PAID {
		return 2, errors.New("学生尚未完成课程支付")
	}

	latestPeriod, err := courseService.GetLatestCompleteChapterPeriod(courseId, studentId, auditionRecord.Id)

	if chapter.Period != latestPeriod+1 {
		return 2, errors.New("课程课时信息异常")
	}

	record := models.CourseChapterToUser{
		CourseId:  courseId,
		ChapterId: chapterId,
		UserId:    studentId,
		TeacherId: teacherId,
		Period:    chapter.Period,
		RecordId:  auditionRecord.Id,
	}

	_, err = models.CreateCourseChapterToUser(&record)
	if err != nil {
		return 2, ErrServerAbnormal
	}

	go lcmessage.SendAuditionCourseChapterCompleteMsg(auditionRecord.Id, chapter.Id)

	recordInfo := map[string]interface{}{
		"Status": models.AUDITION_RECORD_STATUS_COMPLETE,
	}
	_, err = models.UpdateCourseAuditionRecord(auditionRecord.Id, recordInfo)
	if err != nil {
		return 2, ErrServerAbnormal
	}
	return 0, nil
}
