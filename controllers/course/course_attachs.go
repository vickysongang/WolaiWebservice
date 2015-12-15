package course

import (
	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
)

type attachInfo struct {
	models.CourseChapterAttach
	AttachPics []*models.CourseChapterAttachPic `json:"attachPics"`
}

type chapterAttachInfo struct {
	models.CourseChapter
	AttachList []*attachInfo `json:"attachInfo"`
}

func GetCourseAttachs(courseId int64) (int64, []*chapterAttachInfo) {
	o := orm.NewOrm()

	var courseChapters []*models.CourseChapter
	_, err := o.QueryTable("course_chapter").Filter("course_id", courseId).OrderBy("period").All(&courseChapters)
	if err != nil {
		return 2, nil
	}

	infos := make([]*chapterAttachInfo, 0)
	for _, courseChapter := range courseChapters {
		var courseChapterAttach models.CourseChapterAttach
		err := o.QueryTable("course_chapter_attach").Filter("chapter_id", courseChapter.Id).One(&courseChapterAttach)
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
			CourseChapter: *courseChapter,
			AttachList:    aInfos,
		}

		infos = append(infos, &cInfo)
	}

	return 0, infos
}

func GetCourseChapterAttachs(chapterId int64) (int64, *chapterAttachInfo) {
	o := orm.NewOrm()

	courseChapter, err := models.ReadCourseChapter(chapterId)
	if err != nil {
		return 2, nil
	}

	var courseChapterAttach models.CourseChapterAttach
	err = o.QueryTable("course_chapter_attach").Filter("chapter_id", courseChapter.Id).One(&courseChapterAttach)
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
		CourseChapter: *courseChapter,
		AttachList:    aInfos,
	}

	return 0, &cInfo
}
