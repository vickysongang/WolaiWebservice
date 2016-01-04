package course

import (
	"errors"

	"github.com/astaxie/beego/orm"

	"WolaiWebservice/controllers/trade"
	"WolaiWebservice/models"
	"WolaiWebservice/utils/leancloud"
)

func HandleCourseActionNextChapter(userId, studentId, courseId, chapterId int64) (int64, error) {
	var err error
	o := orm.NewOrm()

	_, err = models.ReadUser(userId)
	if err != nil {
		return 2, errors.New("用户信息异常")
	}

	_, err = models.ReadUser(studentId)
	if err != nil {
		return 2, errors.New("用户信息异常")
	}

	_, err = models.ReadCourse(courseId)
	if err != nil {
		return 2, errors.New("课程信息异常")
	}

	chapter, err := models.ReadCourseChapter(chapterId)
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

	latestPeriod, err := queryLatestCourseChapterPeriod(courseId, studentId)
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

	go leancloud.SendCourseChapterCompleteMsg(purchase.Id, chapter.Id)

	chapterCount, _ := o.QueryTable("course_chapter").Filter("course_id", courseId).Count()

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

	err = trade.HandleCourseEarning(purchase.Id, chapter.Period)
	if err != nil {
		return 2, errors.New("支付信息异常")
	}

	return 0, nil
}
