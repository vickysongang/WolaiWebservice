package course

import (
	"WolaiWebservice/models"

	"github.com/astaxie/beego/orm"
)

type coursePreview struct {
	Id             int64  `json:"id"`
	Name           string `json:"name"`
	Cover          string `json:"cover"`
	LongCover      string `json:"longcover"`
	RecommendIntro string `json:"recommendIntro"`
}

type moduleInfo struct {
	Id         int64            `json:"id"`
	Name       string           `json:"name"`
	Type       int64            `json:"type"`
	CourseList []*coursePreview `json:"courseList"`
}

type courseHomePage struct {
	Master  *moduleInfo `json:"master"`
	Synchro *moduleInfo `json:"synchro"`
	Hot     *moduleInfo `json:"hot"`
}

func GetCourseHomePage() (int64, *courseHomePage) {
	o := orm.NewOrm()

	var masterCourses []*models.Course
	_, err := o.QueryTable("course").Filter("type", "1").All(&masterCourses)
	if err != nil {
		return 2, nil
	}

	masterModule, err := models.ReadCourseModule(1)
	if err != nil {
		return 2, nil
	}

	masterPreviews := make([]*coursePreview, 0)
	for _, course := range masterCourses {
		preview := coursePreview{
			Id:             course.Id,
			Name:           course.Name,
			Cover:          course.Cover,
			LongCover:      course.LongCover,
			RecommendIntro: "",
		}

		masterPreviews = append(masterPreviews, &preview)
	}

	masterInfo := moduleInfo{
		Id:         masterModule.Id,
		Name:       masterModule.Name,
		Type:       masterModule.Type,
		CourseList: masterPreviews,
	}

	var synchCourses []*models.Course
	_, err = o.QueryTable("course").Filter("type", "2").All(&synchCourses)
	if err != nil {
		return 2, nil
	}

	synchModule, err := models.ReadCourseModule(2)
	if err != nil {
		return 2, nil
	}

	synchPreviews := make([]*coursePreview, 0)
	for _, course := range synchCourses {
		preview := coursePreview{
			Id:             course.Id,
			Name:           course.Name,
			Cover:          course.Cover,
			LongCover:      course.LongCover,
			RecommendIntro: "",
		}

		synchPreviews = append(synchPreviews, &preview)
	}

	synchInfo := moduleInfo{
		Id:         synchModule.Id,
		Name:       synchModule.Name,
		Type:       synchModule.Type,
		CourseList: synchPreviews,
	}

	var hotCourses []*models.Course
	_, err = o.QueryTable("course").Filter("type", "3").All(&hotCourses)
	if err != nil {
		return 2, nil
	}

	hotModule, err := models.ReadCourseModule(3)
	if err != nil {
		return 2, nil
	}

	hotPreviews := make([]*coursePreview, 0)
	for _, course := range hotCourses {
		preview := coursePreview{
			Id:             course.Id,
			Name:           course.Name,
			Cover:          course.Cover,
			LongCover:      course.LongCover,
			RecommendIntro: "",
		}

		hotPreviews = append(hotPreviews, &preview)
	}

	hotInfo := moduleInfo{
		Id:         hotModule.Id,
		Name:       hotModule.Name,
		Type:       hotModule.Type,
		CourseList: hotPreviews,
	}

	homePage := courseHomePage{
		Master:  &masterInfo,
		Synchro: &synchInfo,
		Hot:     &hotInfo,
	}

	return 0, &homePage
}
