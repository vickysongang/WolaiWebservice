// course_detail_student
package course

import (
	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
	userService "WolaiWebservice/service/user"
)

func GetCourseDetailStudent(userId int64, courseId int64) (int64, *courseDetailStudent) {
	course, err := models.ReadCourse(courseId)
	if err != nil {
		return 2, nil
	}

	studentCount := courseService.GetCourseStudentCount(courseId)

	chapterCount := courseService.GetCourseChapterCount(courseId)

	detail := courseDetailStudent{
		Course:       *course,
		StudentCount: studentCount,
		ChapterCount: chapterCount,
	}

	characteristicList, _ := courseService.QueryCourseContentIntros(courseId)
	detail.CharacteristicList = characteristicList

	purchaseRecord, err := courseService.GetCoursePurchaseRecordByUserId(courseId, userId)
	if err != nil && err != orm.ErrNoRows {
		return 2, nil
	}

	purchaseFlag := (err != orm.ErrNoRows) //判断是否购买或者试听
	if !purchaseFlag {
		detail.AuditionStatus = models.PURCHASE_RECORD_STATUS_IDLE
		detail.PurchaseStatus = models.PURCHASE_RECORD_STATUS_IDLE
		detail.TeacherList, _ = queryCourseTeacherList(courseId)
		detail.ChapterList, _ = queryCourseChapterStatus(courseId, 0, false)
	} else {
		detail.AuditionStatus = purchaseRecord.AuditionStatus
		detail.PurchaseStatus = purchaseRecord.PurchaseStatus
		detail.TeacherList, _ = queryCourseCurrentTeacher(purchaseRecord.TeacherId)

		if purchaseRecord.TeacherId == 0 {
			detail.ChapterList, _ = queryCourseChapterStatus(courseId, 0, false)
		} else {
			detail.ChapterCompletedPeriod, err = courseService.GetLatestCompleteChapterPeriod(courseId, userId, purchaseRecord.Id)
			if err != nil {
				detail.ChapterList, _ = queryCourseCustomChapterStatus(courseId,
					detail.ChapterCompletedPeriod,
					userId,
					purchaseRecord.TeacherId,
					purchaseRecord.Id,
					models.COURSE_TYPE_DELUXE,
					false)
			} else {
				detail.ChapterList, _ = queryCourseCustomChapterStatus(courseId,
					detail.ChapterCompletedPeriod+1,
					userId,
					purchaseRecord.TeacherId,
					purchaseRecord.Id,
					models.COURSE_TYPE_DELUXE,
					false)
			}
		}
	}

	return 0, &detail
}

func queryCourseTeacherList(courseId int64) ([]*teacherItem, error) {
	result := make([]*teacherItem, 0)

	courseTeachers, err := courseService.QueryCourseTeachers(courseId)
	if err != nil {
		return result, nil
	}

	for _, courseTeacher := range courseTeachers {
		user, _ := models.ReadUser(courseTeacher.UserId)
		profile, err := models.ReadTeacherProfile(courseTeacher.UserId)
		if err != nil {
			continue
		}
		school, _ := models.ReadSchool(profile.SchoolId)

		subjectNames, err := userService.GetTeacherSubjectNameSlice(courseTeacher.UserId)
		if err != nil {
			return result, err
		}

		item := teacherItem{
			Id:           courseTeacher.UserId,
			Nickname:     user.Nickname,
			Avatar:       user.Avatar,
			Gender:       user.Gender,
			AccessRight:  user.AccessRight,
			School:       school.Name,
			Intro:        profile.Intro,
			SubjectList:  subjectNames,
			OnlineStatus: "online",
		}
		result = append(result, &item)
	}

	return result, nil
}

func queryCourseCurrentTeacher(teacherId int64) ([]*teacherItem, error) {
	result := make([]*teacherItem, 0)

	user, err := models.ReadUser(teacherId)
	if err != nil {
		return result, err
	}
	profile, _ := models.ReadTeacherProfile(teacherId)
	school, _ := models.ReadSchool(profile.SchoolId)

	subjectNames, err := userService.GetTeacherSubjectNameSlice(teacherId)
	if err != nil {
		return result, err
	}

	item := teacherItem{
		Id:           user.Id,
		Nickname:     user.Nickname,
		Avatar:       user.Avatar,
		Gender:       user.Gender,
		AccessRight:  user.AccessRight,
		School:       school.Name,
		Intro:        profile.Intro,
		SubjectList:  subjectNames,
		OnlineStatus: "online",
	}
	result = append(result, &item)

	return result, nil
}
