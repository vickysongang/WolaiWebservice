package course

import (
	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
)

type attachInfo struct {
	*models.CourseChapterAttach
	AttachPics []*models.CourseChapterAttachPic `json:"attachPics"`
}

type chapterAttachInfo struct {
	models.CourseCustomChapter
	AttachList []*attachInfo `json:"attachInfo"`
}

func GetCourseAttachs(courseId int64) (int64, []*chapterAttachInfo) {
	o := orm.NewOrm()

	var courseCustomChapters []*models.CourseCustomChapter
	_, err := o.QueryTable("course_custom_chapter").
		Filter("rel_id", courseId).
		All(&courseCustomChapters)
	if err != nil {
		return 2, nil
	}

	infos := make([]*chapterAttachInfo, 0)
	for _, courseChapter := range courseCustomChapters {

		if courseChapter.AttachId == 0 {
			continue
		}

		courseChapterAttach, err := models.ReadCourseChapterAttach(courseChapter.AttachId)
		if err != nil {
			continue
		}

		var attachPics []*models.CourseChapterAttachPic
		_, err = o.QueryTable("course_chapter_attach_pic").Filter("attach_id", courseChapterAttach.Id).All(&attachPics)
		if err != nil {
			continue
		}

		aInfo := attachInfo{
			CourseChapterAttach: courseChapterAttach,
			AttachPics:          attachPics,
		}

		aInfos := make([]*attachInfo, 0)
		aInfos = append(aInfos, &aInfo)

		cInfo := chapterAttachInfo{
			CourseCustomChapter: *courseChapter,
			AttachList:          aInfos,
		}

		infos = append(infos, &cInfo)
	}

	return 0, infos
}

func GetCourseChapterAttachs(chapterId int64) (int64, *chapterAttachInfo) {
	o := orm.NewOrm()

	courseChapter, err := models.ReadCourseCustomChapter(chapterId)
	if err != nil {
		return 2, nil
	}

	if courseChapter.AttachId == 0 {
		return 2, nil
	}

	courseChapterAttach, err := models.ReadCourseChapterAttach(courseChapter.AttachId)
	if err != nil {
		return 2, nil
	}

	var attachPics []*models.CourseChapterAttachPic
	_, err = o.QueryTable("course_chapter_attach_pic").Filter("attach_id", courseChapterAttach.Id).All(&attachPics)
	if err != nil {
		return 2, nil
	}

	aInfo := attachInfo{
		CourseChapterAttach: courseChapterAttach,
		AttachPics:          attachPics,
	}

	aInfos := make([]*attachInfo, 0)
	aInfos = append(aInfos, &aInfo)

	cInfo := chapterAttachInfo{
		CourseCustomChapter: *courseChapter,
		AttachList:          aInfos,
	}

	return 0, &cInfo
}
