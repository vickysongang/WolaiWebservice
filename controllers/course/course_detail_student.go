package course

import (
	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
	courseService "WolaiWebservice/service/course"
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
	StudentCount           int64                       `json:"studentCount"`
	ChapterCount           int64                       `json:"chapterCount"`
	AuditionStatus         string                      `json:"auditionStatus"`
	PurchaseStatus         string                      `json:"purchaseStatus"`
	ChapterCompletedPeriod int64                       `json:"chapterCompletePeriod"`
	CharacteristicList     []models.CourseContentIntro `json:"characteristicList"`
	ChapterList            []*courseChapterStatus      `json:"chapterList"`
	TeacherList            []*teacherItem              `json:"teacherList"`
	AuditionCourseId       int64                       `json:"auditionCourseId,omitempty"`
}

func GetCourseDetailStudent(userId int64, courseId int64) (int64, *courseDetailStudent) {
	var err error
	var course *models.Course
	if courseId == 0 { //代表试听课，从H5页面跳转过来的
		course = courseService.QueryAuditionCourse()
		if course == nil {
			return 2, nil
		}
	} else {
		course, err = models.ReadCourse(courseId)
		if err != nil {
			return 2, nil
		}
	}
	if course.Type == models.COURSE_TYPE_DELUXE {
		status, course := GetDeluxeCourseDetail(userId, course)
		return status, course
	} else if course.Type == models.COURSE_TYPE_AUDITION {
		status, course := GetAuditionCourseDetail(userId, course)
		return status, course
	}
	return 0, nil
}

func GetDeluxeCourseDetail(userId int64, course *models.Course) (int64, *courseDetailStudent) {
	var err error
	o := orm.NewOrm()
	courseId := course.Id
	studentCount := courseService.GetCourseStudentCount(courseId)

	detail := courseDetailStudent{
		Course:       *course,
		StudentCount: studentCount,
	}

	characteristicList, _ := courseService.QueryCourseContentIntros(courseId)
	detail.CharacteristicList = characteristicList

	var purchaseRecord models.CoursePurchaseRecord
	err = o.QueryTable(new(models.CoursePurchaseRecord).TableName()).Filter("user_id", userId).Filter("course_id", courseId).
		One(&purchaseRecord)
	if err != nil && err != orm.ErrNoRows {
		return 2, nil
	}

	purchaseFlag := (err != orm.ErrNoRows) //判断是否购买或者试听
	if !purchaseFlag {
		detail.AuditionStatus = models.PURCHASE_RECORD_STATUS_IDLE
		detail.PurchaseStatus = models.PURCHASE_RECORD_STATUS_IDLE
		detail.TeacherList, _ = queryCourseTeacherList(courseId)
		detail.ChapterList, _ = queryCourseChapterStatus(courseId, 0)
		chapterCount := courseService.GetCourseChapterCount(courseId)
		detail.ChapterCount = chapterCount - 1
	} else {
		detail.AuditionStatus = purchaseRecord.AuditionStatus
		detail.PurchaseStatus = purchaseRecord.PurchaseStatus
		detail.TeacherList, _ = queryCourseCurrentTeacher(purchaseRecord.TeacherId)
		detail.ChapterCount = purchaseRecord.ChapterCount
		if purchaseRecord.TeacherId == 0 {
			detail.ChapterList, _ = queryCourseChapterStatus(courseId, 0)
		} else {
			detail.ChapterCompletedPeriod, err = courseService.QueryLatestCourseChapterPeriod(courseId, userId)
			if err != nil {
				detail.ChapterList, _ = queryCourseCustomChapterStatus(courseId, detail.ChapterCompletedPeriod, userId, purchaseRecord.TeacherId)
			} else {
				detail.ChapterList, _ = queryCourseCustomChapterStatus(courseId, detail.ChapterCompletedPeriod+1, userId, purchaseRecord.TeacherId)
			}
		}
	}
	auditionCourse := courseService.QueryAuditionCourse()
	if auditionCourse != nil {
		detail.AuditionCourseId = auditionCourse.Id
	}
	return 0, &detail
}

func GetAuditionCourseDetail(userId int64, course *models.Course) (int64, *courseDetailStudent) {
	o := orm.NewOrm()
	courseId := course.Id
	studentCount := courseService.GetAuditionCourseStudentCount(courseId)
	detail := courseDetailStudent{
		Course:       *course,
		ChapterCount: 1,
		StudentCount: studentCount,
	}
	characteristicList, _ := courseService.QueryCourseContentIntros(courseId)
	detail.CharacteristicList = characteristicList

	var auditionRecord models.CourseAuditionRecord
	err := o.QueryTable(new(models.CourseAuditionRecord).TableName()).
		Filter("course_id", courseId).Filter("user_id", userId).
		Exclude("status", models.AUDITION_RECORD_STATUS_COMPLETE).
		One(&auditionRecord)
	if err != nil && err != orm.ErrNoRows {
		return 2, nil
	}
	purchaseFlag := (err != orm.ErrNoRows)
	if !purchaseFlag {
		detail.AuditionStatus = models.PURCHASE_RECORD_STATUS_IDLE
		detail.PurchaseStatus = models.PURCHASE_RECORD_STATUS_IDLE
		detail.TeacherList = make([]*teacherItem, 0)
		detail.ChapterList, _ = queryCourseChapterStatus(courseId, 0)
	} else {
		detail.AuditionStatus = auditionRecord.Status
		detail.PurchaseStatus = auditionRecord.Status
		detail.TeacherList, _ = queryCourseCurrentTeacher(auditionRecord.TeacherId)
		if auditionRecord.TeacherId == 0 {
			detail.ChapterList, _ = queryCourseCustomChapterStatus(courseId, 0, userId, auditionRecord.TeacherId)
		} else {
			detail.ChapterList, _ = queryCourseCustomChapterStatus(courseId, 1, userId, auditionRecord.TeacherId)
		}
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
