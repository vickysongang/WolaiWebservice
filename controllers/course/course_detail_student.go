package course

import (
	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
	userService "WolaiWebservice/service/user"
)

type teacherItem struct {
	Id           int64    `json:"id"`
	Nickname     string   `json:"nickname"`
	Avatar       string   `json:"avatar"`
	Gender       int64    `json:"gender"`
	AccessRight  int64    `json:"accessRight"`
	School       string   `json:"school"`
	Major        string   `json:"major"`
	Intro        string   `json:"intro"`
	SubjectList  []string `json:"subjectList,omitempty"`
	OnlineStatus string   `json:"onlineStatus,omitempty"`
}

type courseDetailStudent struct {
	models.Course
	StudentCount           int64                                `json:"studentCount"`
	ChapterCount           int64                                `json:"chapterCount"`
	AuditionStatus         string                               `json:"auditionStatus"`
	PurchaseStatus         string                               `json:"purchaseStatus"`
	ChapterCompletedPeriod int64                                `json:"chapterCompletePeriod"`
	CharacteristicList     []models.CourseContentCharacteristic `json:"characteristicList"`
	ChapterList            []*courseChapterStatus               `json:"chapterList"`
	TeacherList            []*teacherItem                       `json:"teacherList"`
}

func GetCourseDetailStudent(userId int64, courseId int64) (int64, *courseDetailStudent) {
	o := orm.NewOrm()

	course, err := models.ReadCourse(courseId)
	if err != nil {
		return 2, nil
	}

	studentCount := queryCourseStudentCount(courseId)

	chapterCount, _ := o.QueryTable("course_chapter").Filter("course_id", courseId).Count()

	detail := courseDetailStudent{
		Course:       *course,
		StudentCount: studentCount,
		ChapterCount: chapterCount - 1,
	}
	characteristicList, _ := queryCourseContentCharacteristics(courseId)
	detail.CharacteristicList = characteristicList

	var purchaseRecord models.CoursePurchaseRecord
	err = o.QueryTable("course_purchase_record").Filter("user_id", userId).Filter("course_id", courseId).
		One(&purchaseRecord)
	if err != nil && err != orm.ErrNoRows {
		return 2, nil
	}

	purchaseFlag := (err != orm.ErrNoRows)
	if !purchaseFlag {
		detail.AuditionStatus = models.PURCHASE_RECORD_STATUS_IDLE
		detail.PurchaseStatus = models.PURCHASE_RECORD_STATUS_IDLE
		detail.TeacherList, _ = queryCourseTeacherList(courseId)
	} else {
		detail.AuditionStatus = purchaseRecord.AuditionStatus
		detail.PurchaseStatus = purchaseRecord.PurchaseStatus
		detail.TeacherList, _ = queryCourseCurrentTeacher(purchaseRecord.TeacherId)
	}

	detail.ChapterCompletedPeriod, err = queryLatestCourseChapterPeriod(courseId, userId)
	if err != nil {
		detail.ChapterList, _ = queryCourseChapterStatus(courseId, detail.ChapterCompletedPeriod)
	} else {
		detail.ChapterList, _ = queryCourseChapterStatus(courseId, detail.ChapterCompletedPeriod+1)
	}

	return 0, &detail
}

func queryCourseTeacherList(courseId int64) ([]*teacherItem, error) {
	o := orm.NewOrm()

	result := make([]*teacherItem, 0)

	var courseTeachers []*models.CourseToTeacher
	_, err := o.QueryTable("course_to_teachers").Filter("course_id", courseId).All(&courseTeachers)
	if err != nil {
		return result, nil
	}

	for _, courseTeacher := range courseTeachers {
		user, _ := models.ReadUser(courseTeacher.UserId)
		profile, _ := models.ReadTeacherProfile(courseTeacher.UserId)
		school, _ := models.ReadSchool(profile.SchoolId)

		subjectNames, err := userService.GetTeacherSubjectNameSlice(courseTeacher.UserId)
		if err != nil {
			return result, err
		}

		item := teacherItem{
			Id:           courseTeacher.UserId,
			Nickname:     user.Nickname,
			Avatar:       user.Avatar,
			Gender:       user.Gender,
			AccessRight:  user.AccessRight,
			School:       school.Name,
			Intro:        profile.Intro,
			SubjectList:  subjectNames,
			OnlineStatus: "online",
		}
		result = append(result, &item)
	}

	return result, nil
}

func queryCourseCurrentTeacher(teacherId int64) ([]*teacherItem, error) {
	result := make([]*teacherItem, 0)

	user, err := models.ReadUser(teacherId)
	if err != nil {
		return result, err
	}
	profile, _ := models.ReadTeacherProfile(teacherId)
	school, _ := models.ReadSchool(profile.SchoolId)

	subjectNames, err := userService.GetTeacherSubjectNameSlice(teacherId)
	if err != nil {
		return result, err
	}

	item := teacherItem{
		Id:           user.Id,
		Nickname:     user.Nickname,
		Avatar:       user.Avatar,
		Gender:       user.Gender,
		AccessRight:  user.AccessRight,
		School:       school.Name,
		Intro:        profile.Intro,
		SubjectList:  subjectNames,
		OnlineStatus: "online",
	}
	result = append(result, &item)

	return result, nil
}
