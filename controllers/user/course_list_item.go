package user

import (
	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
	userService "WolaiWebservice/service/user"
)

type CourseListItem struct {
	models.Course
	StudentCount int64 `json:"studentCount"`
	ChapterCount int64 `json:"chapterCount"`
}

func AssembleTeacherCourseList(teacherId, page, count int64) ([]*CourseListItem, error) {
	var err error

	result := make([]*CourseListItem, 0)

	courses, err := userService.GetTeacherCourses(teacherId, page, count)
	if err != nil {
		return result, nil
	}

	for _, course := range courses {
		item := CourseListItem{
			Course:       *course,
			StudentCount: courseService.GetCourseStudentCount(course.Id),
			ChapterCount: courseService.GetCourseChapterCount(course.Id),
		}

		result = append(result, &item)
	}

	return result, nil
}
