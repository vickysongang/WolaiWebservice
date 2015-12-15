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

		courseChapters, _ := queryCourseChapters(course.Id)
		item.ChapterCount = int64(len(courseChapters))
		item.StudentCount = queryCourseStudentCount(course.Id)

		courses = append(courses, &item)
	}
	// courses, _ := queryModuleCourses(moduleType, page, count)

	// moduleCourses := make([]*courseItem, 0)

	// for _, course := range courses {
	// 	moduleCourseDisplayInfo := courseItem{}
	// 	moduleCourseDisplayInfo.Id = course.Id
	// 	moduleCourseDisplayInfo.Name = course.Name
	// 	moduleCourseDisplayInfo.ImgCover = course.ImgCover
	// 	moduleCourseDisplayInfo.ImgLongCover = course.ImgLongCover

	// 	//获取课程的在学学生数
	// 	studentCount := queryCourseStudentCount(course.Id)
	// 	moduleCourseDisplayInfo.StudentCount = studentCount

	// 	//获取课时数
	// 	courseChapters, _ := queryCourseChapters(course.Id)
	// 	chapterCount := len(courseChapters)
	// 	moduleCourseDisplayInfo.ChapterCount = int64(chapterCount)
	// 	moduleCourses = append(moduleCourses, &moduleCourseDisplayInfo)
	// }
	return 0, courses
}

//查询模块的全部课程
// func queryModuleCourses(moduleType, page, count int64) ([]models.Course, error) {
// 	start := page * count
// 	o := orm.NewOrm()
// 	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
// 	qb.Select("course.id,course.name,course.type,course.grade_id,course.subject_id," +
// 		"course.time_from,course.time_to,course.cover,course.long_cover,course.intro,course.create_time,course.creator").
// 		From("course").InnerJoin("course_to_module").On("course.id = course_to_module.course_id").
// 		InnerJoin("course_module").On("course_to_module.module_id =  course_module.id").
// 		Where("course_module.type = ?").Limit(int(count)).Offset(int(start))
// 	sql := qb.String()
// 	courses := make([]models.Course, 0)
// 	_, err := o.Raw(sql, moduleType).QueryRows(&courses)
// 	return courses, err
// }

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
