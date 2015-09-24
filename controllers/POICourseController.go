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
		newUserToCourse.Status = "unopen"
		course4User.UserToCourse = &newUserToCourse
	} else {
		course4User.UserToCourse = userToCourse
	}
	course4User.CurrentTime = time.Now()
	return course4User, nil
}

func JoinCourse(userId, courseId int64) (models.POICourse4User, error) {
	userToCourse := models.POIUserToCourse{UserId: userId, CourseId: courseId, Status: models.COURSE_OPENING}
	models.InsertUserToCourse(&userToCourse)
	course4User, _ := QueryUserCourse(userId, courseId)
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
	course4User, _ := QueryUserCourse(userId, courseId)
	return course4User, nil
}

func RenewUserCourse(userId, courseId int64) (models.POICourse4User, error) {
	now := time.Now()
	course, _ := models.QueryCourseById(courseId)
	length := course.Length
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
		timeTo = now.AddDate(0, 0, int(length+1))
	} else if timeUnit == "M" {
		timeTo = now.AddDate(0, int(length), 1)
	} else if timeUnit == "Y" {
		timeTo = now.AddDate(int(length), 0, 1)
	}
	updateInfo := map[string]interface{}{
		"Status":   models.COURSE_SERVING,
		"TimeFrom": timeFrom.Format(TIME_FORMAT),
		"TimeTo":   timeTo.Format(TIME_FORMAT),
	}
	models.UpdateUserCourseInfo(userId, courseId, updateInfo)
	course4User, _ := QueryUserCourse(userId, courseId)
	return course4User, nil
}
