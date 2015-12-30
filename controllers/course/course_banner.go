package course

import (
	"WolaiWebservice/models"

	"github.com/astaxie/beego/orm"
)

func GetCourseBanners() (int64, []*models.CourseBanners) {
	o := orm.NewOrm()

	var banners []*models.CourseBanners
	_, err := o.QueryTable("course_banners").OrderBy("rank").All(&banners)
	if err != nil {
		return 2, nil
	}

	return 0, banners
}
