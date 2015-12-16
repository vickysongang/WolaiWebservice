package course

import (
	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
	"WolaiWebservice/utils"
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
		item.ChapterCount = count
		item.StudentCount = queryCourseStudentCount(course.Id)

		courses = append(courses, &item)
	}

	return 0, courses
}

//查询课程在学的学生数,此处的判断逻辑为只要学生购买了该课程，就认为学生在学该课程
func queryCourseStudentCount(courseId int64) int64 {
	o := orm.NewOrm()
	count, _ := o.QueryTable("course_purchase_record").Filter("course_id", courseId).Filter("status__gte", 6).Count()
	return count
}

//查询课程的章节
func queryCourseChapters(courseId int64) ([]models.CourseChapter, error) {
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("id,course_id,title,abstract,period,create_time").
		From("course_chapter").
		Where("course_id = ?").
		OrderBy("period").Asc()
	sql := qb.String()
	courseChapters := make([]models.CourseChapter, 0)
	_, err := o.Raw(sql, courseId).QueryRows(&courseChapters)
	return courseChapters, err
}
