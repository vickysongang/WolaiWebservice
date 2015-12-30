package course

import (
	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
)

type courseItem struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	ImgCover     string `json:"imgCover"`
	ImgLongCover string `json:"imgLongCover"`
	StudentCount int64  `json:"studentCount"`
	ChapterCount int64  `json:"chapterCount"`
}

func GetCourseModuleList(moduleType, page, count int64) (int64, []*courseItem) {
	o := orm.NewOrm()

	var courseModules []*models.CourseToModule
	_, err := o.QueryTable("course_to_module").Filter("module_id", moduleType).
		OrderBy("rank").Offset(page * count).Limit(count).All(&courseModules)
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

		count, _ := o.QueryTable("course_chapter").Filter("course_id", courseModule.CourseId).Count()
		item.ChapterCount = count - 1
		item.StudentCount = queryCourseStudentCount(course.Id)

		courses = append(courses, &item)
	}

	return 0, courses
}
