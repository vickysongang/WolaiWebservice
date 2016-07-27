package course

import (
	"errors"

	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
	"WolaiWebservice/utils/leancloud/lcmessage"
)

func HandleCourseActionNextChapter(userId, studentId, courseId, chapterId int64) (int64, error) {
	var err error

	_, err = models.ReadUser(userId)
	if err != nil {
		return 2, ErrUserAbnormal
	}

	_, err = models.ReadUser(studentId)
	if err != nil {
		return 2, ErrUserAbnormal
	}

	_, err = models.ReadCourse(courseId)
	if err != nil {
		return 2, ErrCourseAbnormal
	}

	chapter, err := models.ReadCourseCustomChapter(chapterId)
	if err != nil {
		return 2, ErrCourseAbnormal
	}

	purchase, err := courseService.GetCoursePurchaseRecordByUserId(courseId, studentId)
	if err != nil {
		return 2, ErrCourseAbnormal
	}
	if purchase.TeacherId != userId {
		return 2, ErrCourseAbnormal
	}

	latestPeriod, err := courseService.GetLatestCompleteChapterPeriod(courseId, studentId, purchase.Id)
	if err == nil {
		if chapter.Period != latestPeriod+1 {
			return 2, ErrCourseAbnormal
		}

		if purchase.PurchaseStatus != models.PURCHASE_RECORD_STATUS_PAID {
			return 2, errors.New("学生尚未完成课程支付")
		}

	} else {
		if latestPeriod != 0 {
			return 2, ErrCourseAbnormal
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
		RecordId:  purchase.Id,
	}

	_, err = models.CreateCourseChapterToUser(&record)
	if err != nil {
		return 2, ErrServerAbnormal
	}

	go lcmessage.SendCourseChapterCompleteMsg(purchase.Id, chapter.Id)

	chapterCount := purchase.ChapterCount

	recordInfo := map[string]interface{}{
		"AuditionStatus": purchase.AuditionStatus,
	}
	models.UpdateCoursePurchaseRecord(purchase.Id, recordInfo)

	if chapter.Period == 0 {
		recordInfo := map[string]interface{}{
			"AuditionStatus": models.PURCHASE_RECORD_STATUS_COMPLETE,
		}
		models.UpdateCoursePurchaseRecord(purchase.Id, recordInfo)
	} else if chapter.Period == chapterCount {
		recordInfo := map[string]interface{}{
			"PurchaseStatus": models.PURCHASE_RECORD_STATUS_COMPLETE,
		}
		models.UpdateCoursePurchaseRecord(purchase.Id, recordInfo)
	}

	return 0, nil
}
