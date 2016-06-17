package course

import (
	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
)

func GetCourseDetailTeacher(courseId, studentId int64) (int64, *courseDetailTeacher) {
	purchaseRecord, err := courseService.GetCoursePurchaseRecordByUserId(courseId, studentId)
	if err != nil {
		return 2, nil
	}

	course, err := models.ReadCourse(courseId)
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

	detail.ChapterCompletedPeriod, err = courseService.GetLatestCompleteChapterPeriod(courseId, studentId, purchaseRecord.Id)
	if err != nil {
		detail.ChapterList, _ = queryCourseCustomChapterStatus(
			courseId,
			detail.ChapterCompletedPeriod,
			studentId,
			purchaseRecord.TeacherId,
			purchaseRecord.Id,
			models.COURSE_TYPE_DELUXE,
			false)
	} else {
		detail.ChapterList, _ = queryCourseCustomChapterStatus(
			courseId,
			detail.ChapterCompletedPeriod+1,
			studentId,
			purchaseRecord.TeacherId,
			purchaseRecord.Id,
			models.COURSE_TYPE_DELUXE,
			false)
	}

	studentList := make([]*models.User, 0)
	student, _ := models.ReadUser(studentId)
	studentList = append(studentList, student)

	detail.StudentList = studentList

	return 0, &detail
}
