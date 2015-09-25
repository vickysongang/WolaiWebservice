// POICourseController
package controllers

import (
	"POIWolaiWebService/models"
	"time"
)

const (
	TIME_FORMAT = "2006-01-02 00:00:00"
)

func QueryUserCourses(userId int64) (models.POICourseInfos, error) {
	courseInfos := make(models.POICourseInfos, 0)
	courses, _ := models.QueryCourses()
	for _, course := range courses {
		course4User := models.POICourse4User{Course: course, CurrentTime: time.Now()}
		userToCourse, _ := models.QueryUserToCourse(course.Id, userId)
		if userToCourse == nil {
			newUserToCourse := models.POIUserToCourse{}
			newUserToCourse.UserId = userId
			newUserToCourse.CourseId = course.Id
			newUserToCourse.Status = models.COURSE_UNOPEN
			course4User.UserToCourse = &newUserToCourse
		} else {
			course4User.UserToCourse = userToCourse
		}
		course4User.CurrentTime = time.Now()
		course4User.JoinCount = models.GetCourseJoinCount(course.Id) + 555
		courseInfos = append(courseInfos, &course4User)
	}
	return courseInfos, nil
}

func QueryUserCourse(userId, courseId int64) (models.POICourse4User, error) {
	course, _ := models.QueryCourseById(courseId)
	course4User := models.POICourse4User{Course: course}
	userToCourse, _ := models.QueryUserToCourse(courseId, userId)
	if userToCourse == nil {
		newUserToCourse := models.POIUserToCourse{}
		newUserToCourse.UserId = userId
		newUserToCourse.CourseId = course.Id
		newUserToCourse.Status = models.COURSE_UNOPEN
		course4User.UserToCourse = &newUserToCourse
	} else {
		course4User.UserToCourse = userToCourse
	}
	course4User.JoinCount = models.GetCourseJoinCount(course.Id) + 555
	course4User.CurrentTime = time.Now()
	return course4User, nil
}

func JoinCourse(userId, courseId int64) (models.POICourse4User, error) {
	userToCourse := models.POIUserToCourse{UserId: userId, CourseId: courseId, Status: models.COURSE_OPENING}
	models.InsertUserToCourse(&userToCourse)
	course4User, _ := QueryUserCourse(userId, courseId)

	purchaseRecord := models.POICoursePurchaseRecord{UserId: userId, CourseId: courseId, Type: models.COURSE_JOIN, Status: models.COURSE_PURCHASE_PENDING}
	models.InsertCoursePurchaseRecord(&purchaseRecord)

	return course4User, nil
}

func ActiveUserCourse(userId, courseId int64) (models.POICourse4User, error) {
	now := time.Now()
	giveCourse, _ := models.QueryGiveCourse()
	timeTo := now.AddDate(0, 0, int(giveCourse.Length)+1)
	updateInfo := map[string]interface{}{
		"Status":   models.COURSE_SERVING,
		"TimeFrom": now.Format(TIME_FORMAT),
		"TimeTo":   timeTo.Format(TIME_FORMAT),
	}
	models.UpdateUserCourseInfo(userId, courseId, updateInfo)

	purchaseRecordInfo := map[string]interface{}{
		"Status": models.COURSE_PURCHASE_COMPLETED,
		"TimeTo": timeTo.Format(TIME_FORMAT),
	}
	models.UpdatePurchaseRecord(userId, courseId, purchaseRecordInfo)
	course4User, _ := QueryUserCourse(userId, courseId)

	return course4User, nil
}

func UserRenewCourse(userId, courseId int64) (models.POICourse4User, error) {
	purchaseRecord, _ := models.QueryPendingPurchaseRecord(userId, courseId)
	if purchaseRecord == nil {
		newPurchaseRecord := models.POICoursePurchaseRecord{UserId: userId, CourseId: courseId, Type: models.COURSE_RENEW, Status: models.COURSE_PURCHASE_PENDING}
		models.InsertCoursePurchaseRecord(&newPurchaseRecord)
	}
	course4User, _ := QueryUserCourse(userId, courseId)
	return course4User, nil
}

func SupportRenewUserCourse(userId, courseId int64, renewCount int64) (models.POICourse4User, error) {
	now := time.Now()
	course, _ := models.QueryCourseById(courseId)
	length := course.Length * renewCount
	timeUnit := course.TimeUnit
	userToCourse, _ := models.QueryUserToCourse(courseId, userId)
	status := userToCourse.Status
	var timeFrom time.Time
	if status == models.COURSE_EXPIRED {
		timeFrom = now
	} else {
		timeFrom = userToCourse.TimeFrom
	}
	var timeTo time.Time
	if timeUnit == "D" {
		timeTo = userToCourse.TimeTo.AddDate(0, 0, int(length+1))
	} else if timeUnit == "M" {
		timeTo = userToCourse.TimeTo.AddDate(0, int(length), 1)
	} else if timeUnit == "Y" {
		timeTo = userToCourse.TimeTo.AddDate(int(length), 0, 1)
	}
	updateInfo := map[string]interface{}{
		"Status":   models.COURSE_SERVING,
		"TimeFrom": timeFrom.Format(TIME_FORMAT),
		"TimeTo":   timeTo.Format(TIME_FORMAT),
	}
	models.UpdateUserCourseInfo(userId, courseId, updateInfo)

	purchaseRecordInfo := map[string]interface{}{
		"Status": models.COURSE_PURCHASE_COMPLETED,
		"TimeTo": timeTo.Format(TIME_FORMAT),
		"Count":  renewCount,
	}
	models.UpdatePurchaseRecord(userId, courseId, purchaseRecordInfo)

	course4User, _ := QueryUserCourse(userId, courseId)
	return course4User, nil
}
