// POICourseController
package controllers

import (
	"POIWolaiWebService/models"
	"POIWolaiWebService/utils"
	"time"
)

const (
	TIME_TO_FORMAT = "2006-01-02"
	TIME_HHMMSS    = " 23:59:59"
)

/*
 * 获取课程列表
 */
func QueryUserCourses(userId int64) (models.POICourseInfos, error) {
	courseInfos := make(models.POICourseInfos, 0)
	var courseType int64
	if models.IsUserJoinCourse(userId) {
		_, err := models.QueryGiveCourse4User(userId)
		if err == nil {
			courseType = 0
		} else {
			courseType = 1
		}
	} else {
		courseType = 0
	}
	courses, _ := models.QueryCourses(courseType)
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
		course4User.JoinCount = models.GetCourseJoinCount() + 555
		courseInfos = append(courseInfos, &course4User)
	}
	return courseInfos, nil
}

/*
 * 获取课程信息
 */
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
	course4User.JoinCount = models.GetCourseJoinCount() + 555
	course4User.CurrentTime = time.Now()
	return course4User, nil
}

/*
 * 加入课程
 */
func JoinCourse(userId, courseId int64) (models.POICourse4User, error) {
	userToCourse := models.POIUserToCourse{UserId: userId, CourseId: courseId, Status: models.COURSE_OPENING}
	models.InsertUserToCourse(&userToCourse)
	course4User, _ := QueryUserCourse(userId, courseId)

	purchaseRecord := models.POICoursePurchaseRecord{UserId: userId, CourseId: courseId, Type: models.COURSE_JOIN, Status: models.COURSE_PURCHASE_PENDING}
	models.InsertCoursePurchaseRecord(&purchaseRecord)

	return course4User, nil
}

/*
 * 客服通过用户的加入课程申请
 */
func ActiveUserCourse(userId, courseId int64) (models.POICourse4User, error) {
	now := time.Now()
	course, _ := models.QueryCourseById(courseId)
	var timeTo time.Time
	if course.TimeUnit == "D" {
		timeTo = now.AddDate(0, 0, int(course.Length))
	} else if course.TimeUnit == "M" {
		timeTo = now.AddDate(0, int(course.Length), 0)
	} else if course.TimeUnit == "Y" {
		timeTo = now.AddDate(int(course.Length), 0, 0)
	}

	updateInfo := map[string]interface{}{
		"Status":   models.COURSE_SERVING,
		"TimeFrom": now.Format(utils.TIME_FORMAT),
		"TimeTo":   timeTo.Format(TIME_TO_FORMAT) + TIME_HHMMSS,
	}
	models.UpdateUserCourseInfo(userId, courseId, updateInfo)

	purchaseRecordInfo := map[string]interface{}{
		"Status": models.COURSE_PURCHASE_COMPLETED,
		"TimeTo": timeTo.Format(TIME_TO_FORMAT) + TIME_HHMMSS,
	}
	models.UpdatePurchaseRecord(userId, courseId, purchaseRecordInfo)
	course4User, _ := QueryUserCourse(userId, courseId)

	return course4User, nil
}

/*
 * 用户续期课程
 */
func UserRenewCourse(userId, courseId int64) (models.POICourse4User, error) {
	purchaseRecord, _ := models.QueryPendingPurchaseRecord(userId, courseId)
	if purchaseRecord == nil {
		var purchaseType string
		course, _ := models.QueryCourseById(courseId)
		if course.Type == models.COURSE_GIVE_TYPE {
			purchaseType = models.COURSE_UPGRADE
		} else {
			purchaseType = models.COURSE_RENEW
		}
		newPurchaseRecord := models.POICoursePurchaseRecord{UserId: userId, CourseId: courseId, Type: purchaseType, Status: models.COURSE_PURCHASE_PENDING}
		models.InsertCoursePurchaseRecord(&newPurchaseRecord)
	}
	course4User, _ := QueryUserCourse(userId, courseId)
	return course4User, nil
}

/*
 * 客服处理用户的续期申请
 */
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
	newCourseId := courseId
	var timeTo time.Time
	if course.Type == models.COURSE_GIVE_TYPE {
		timeTo = timeFrom.AddDate(0, 1, 0)
		newCourseId, _ = models.QueryDefaultCourseId()
	} else {
		if timeUnit == "D" {
			timeTo = userToCourse.TimeTo.AddDate(0, 0, int(length))
		} else if timeUnit == "M" {
			timeTo = userToCourse.TimeTo.AddDate(0, int(length), 0)
		} else if timeUnit == "Y" {
			timeTo = userToCourse.TimeTo.AddDate(int(length), 0, 0)
		}
	}
	updateInfo := map[string]interface{}{
		"CourseId": newCourseId,
		"Status":   models.COURSE_SERVING,
		"TimeFrom": timeFrom.Format(utils.TIME_FORMAT),
		"TimeTo":   timeTo.Format(TIME_TO_FORMAT) + TIME_HHMMSS,
	}
	models.UpdateUserCourseInfo(userId, courseId, updateInfo)

	purchaseRecordInfo := map[string]interface{}{
		"Status": models.COURSE_PURCHASE_COMPLETED,
		"TimeTo": timeTo.Format(TIME_TO_FORMAT) + TIME_HHMMSS,
		"Count":  renewCount,
	}
	models.UpdatePurchaseRecord(userId, courseId, purchaseRecordInfo)

	course4User, _ := QueryUserCourse(userId, newCourseId)
	return course4User, nil
}

/*
 * 客服拒绝用户的续期申请
 */
func SupportRejectUserCourse(userId, courseId int64) (models.POICourse4User, error) {
	course4User, _ := QueryUserCourse(userId, courseId)
	return course4User, nil
}
