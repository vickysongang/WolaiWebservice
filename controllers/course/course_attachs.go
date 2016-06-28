package course

import (
	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
)

type attachInfo struct {
	*models.CourseChapterAttach
	AttachPics []*models.CourseChapterAttachPic `json:"attachPics"`
}

type chapterAttachInfo struct {
	models.CourseCustomChapter
	AttachList []*attachInfo `json:"attachInfo"`
}

func GetCourseAttachs(relId int64) (int64, []*chapterAttachInfo) {
	courseCustomChapters, err := courseService.QueryCustomChaptersByRelId(relId)
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

		attachPics, err := courseService.QueryChapterAttachPics(courseChapterAttach.Id)
		if err != nil {
			continue
		}
		for _, attachPic := range attachPics {
			attachPic.ChapterId = courseChapter.Id
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

	attachPics, err := courseService.QueryChapterAttachPics(courseChapterAttach.Id)
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
