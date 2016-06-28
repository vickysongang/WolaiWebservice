package course

import (
	"WolaiWebservice/models"

	courseService "WolaiWebservice/service/course"
)

func GetCourseBanners() (int64, []*models.CourseBanners) {
	banners, err := courseService.QueryCourseBanners()
	if err != nil {
		return 2, nil
	}

	return 0, banners
}
