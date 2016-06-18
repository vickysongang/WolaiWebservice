package course

import (
	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
)

type courseItem struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	ImgCover     string `json:"imgCover"`
	ImgLongCover string `json:"imgLongCover"`
	StudentCount int64  `json:"studentCount"`
	ChapterCount int64  `json:"chapterCount"`
}

func GetCourseModuleList(moduleId, page, count int64) (int64, []*courseItem) {

	courseModules, err := courseService.QueryCourseModules(moduleId, page, count)
	if err != nil {
		return 2, nil
	}

	courses := make([]*courseItem, 0)
	for _, courseModule := range courseModules {
		course, err := models.ReadCourse(courseModule.CourseId)
		if err != nil {
			continue
		}

		item := courseItem{
			Id:           course.Id,
			Name:         course.Name,
			ImgCover:     course.ImgCover,
			ImgLongCover: course.ImgLongCover,
		}

		chapterCount := courseService.GetCourseChapterCount(courseModule.CourseId)
		item.ChapterCount = chapterCount
		item.StudentCount = courseService.GetCourseStudentCount(course.Id)

		courses = append(courses, &item)
	}

	return 0, courses
}
