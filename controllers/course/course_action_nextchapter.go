package course

import (
	"errors"

	"WolaiWebservice/models"
)

func HandleCourseActionNextChapter(userId, studentId, courseId, chapterId int64) (int64, error) {
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

	latestPeriod := queryLatestCourseChapterPeriod(courseId, studentId)
	if chapter.Period != 0 {
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

	return 0, nil
}
