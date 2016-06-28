package course

import (
	"github.com/cihub/seelog"

	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
)

type coursePreview struct {
	Id             int64  `json:"id"`
	Name           string `json:"name"`
	ImgCover       string `json:"imgCover"`
	ImgLongCover   string `json:"imgLongCover"`
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
	masterCourses, err := courseService.QueryModuleCourses(1)
	if err != nil {
		seelog.Error(err.Error())
		return 2, nil
	}

	masterModule, err := models.ReadCourseModule(1)
	if err != nil {
		seelog.Error(err.Error())
		return 2, nil
	}

	masterPreviews := make([]*coursePreview, 0)
	for _, courseModule := range masterCourses {
		course, err := models.ReadCourse(courseModule.CourseId)
		if err != nil {
			continue
		}

		preview := coursePreview{
			Id:             course.Id,
			Name:           course.Name,
			ImgCover:       course.ImgCover,
			ImgLongCover:   course.ImgLongCover,
			RecommendIntro: course.RecommendIntro,
		}

		masterPreviews = append(masterPreviews, &preview)
	}

	masterInfo := moduleInfo{
		Id:         masterModule.Id,
		Name:       masterModule.Name,
		Type:       masterModule.Type,
		CourseList: masterPreviews,
	}

	synchCourses, err := courseService.QueryModuleCourses(2)
	if err != nil {
		seelog.Error(err.Error())
		return 2, nil
	}

	synchModule, err := models.ReadCourseModule(2)
	if err != nil {
		seelog.Error(err.Error())
		return 2, nil
	}

	synchPreviews := make([]*coursePreview, 0)
	for _, courseModule := range synchCourses {
		course, err := models.ReadCourse(courseModule.CourseId)
		if err != nil {
			continue
		}

		preview := coursePreview{
			Id:             course.Id,
			Name:           course.Name,
			ImgCover:       course.ImgCover,
			ImgLongCover:   course.ImgLongCover,
			RecommendIntro: course.RecommendIntro,
		}

		synchPreviews = append(synchPreviews, &preview)
	}

	synchInfo := moduleInfo{
		Id:         synchModule.Id,
		Name:       synchModule.Name,
		Type:       synchModule.Type,
		CourseList: synchPreviews,
	}

	hotCourses, err := courseService.QueryModuleCourses(3)
	if err != nil {
		seelog.Error(err.Error())
		return 2, nil
	}

	hotModule, err := models.ReadCourseModule(3)
	if err != nil {
		seelog.Error(err.Error())
		return 2, nil
	}

	hotPreviews := make([]*coursePreview, 0)
	for _, courseModule := range hotCourses {
		course, err := models.ReadCourse(courseModule.CourseId)
		if err != nil {
			continue
		}

		preview := coursePreview{
			Id:             course.Id,
			Name:           course.Name,
			ImgCover:       course.ImgCover,
			ImgLongCover:   course.ImgLongCover,
			RecommendIntro: course.RecommendIntro,
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
