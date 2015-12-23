package course

import (
	"errors"

	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
)

func HandleCourseActionNextChapter(userId, studentId, courseId, chapterId int64) (int64, error) {
	o := orm.NewOrm()

	_, err := models.ReadUser(userId)
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
		if latestPeriod != chapter.Period+1 {
			return 2, errors.New("课程信息异常")
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
		TeacherId: userId,
		Period:    chapter.Period,
	}

	_, err = models.CreateCourseChapterToUser(&record)
	if err != nil {
		return 2, errors.New("服务器操作异常")
	}

	chapterCount, _ := o.QueryTable("course_chapter").Filter("course_id", courseId).Count()

	if chapter.Period == 0 {
		recordInfo := map[string]interface{}{
			"audition_status": models.PURCHASE_RECORD_STATUS_COMPLETE,
		}
		models.UpdateCoursePurchaseRecord(record.Id, recordInfo)
	} else if chapter.Period == chapterCount-1 {
		recordInfo := map[string]interface{}{
			"purchase_status": models.PURCHASE_RECORD_STATUS_COMPLETE,
		}
		models.UpdateCoursePurchaseRecord(record.Id, recordInfo)
	}

	return 0, nil
}
