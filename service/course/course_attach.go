// course_attach
package course

import (
	"WolaiWebservice/models"

	"github.com/astaxie/beego/orm"
)

func QueryChapterAttachPics(attachId int64) ([]*models.CourseChapterAttachPic, error) {
	o := orm.NewOrm()
	var attachPics []*models.CourseChapterAttachPic
	_, err := o.QueryTable("course_chapter_attach_pic").
		Filter("attach_id", attachId).All(&attachPics)
	return attachPics, err
}
