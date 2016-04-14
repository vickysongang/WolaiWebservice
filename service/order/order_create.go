package order

import (
	"errors"
	"time"

	"WolaiWebservice/models"
)

func CreateOrder(creatorId, gradeId, subjectId, teacherId, tierId, recordId, chapterId int64,
	orderType string) (*models.Order, error) {
	var err error

	_, err = models.ReadUser(creatorId)
	if err != nil {
		return nil, err
	}

	_, err = models.ReadGrade(gradeId)
	if err != nil {
		return nil, err
	}

	_, err = models.ReadSubject(subjectId)
	if err != nil {
		return nil, err
	}

	if teacherId != 0 {
		profile, err := models.ReadTeacherProfile(teacherId)
		if err != nil {
			return nil, err
		}

		tierId = profile.TierId
	}

	var priceHourly, salaryHourly int64
	if tierId != 0 {
		tier, err := models.ReadTeacherTierHourly(tierId)
		if err != nil {
			return nil, err
		}

		priceHourly = tier.QAPriceHourly
		salaryHourly = tier.QASalaryHourly
	}

	var courseId int64
	if orderType == models.ORDER_TYPE_COURSE_INSTANT {
		if recordId != 0 {
			record, err := models.ReadCoursePurchaseRecord(recordId)
			if err != nil {
				return nil, err
			}

			courseId = record.CourseId
			priceHourly = record.PriceHourly
			salaryHourly = record.SalaryHourly
		}
	} else if orderType == models.ORDER_TYPE_AUDITION_COURSE_INSTANT {
		if recordId != 0 {
			record, err := models.ReadCourseAuditionRecord(recordId)
			if err != nil {
				return nil, err
			}
			courseId = record.CourseId
			priceHourly = record.PriceHourly
			salaryHourly = record.SalaryHourly
		}
	}

	if chapterId != 0 {
		chapter, err := models.ReadCourseCustomChapter(chapterId)
		if err != nil {
			return nil, err
		}

		if courseId != chapter.CourseId {
			return nil, errors.New("课程信息不匹配")
		}
	}

	date := time.Now().Format(time.RFC3339)

	order := models.Order{
		Creator:      creatorId,
		GradeId:      gradeId,
		SubjectId:    subjectId,
		Date:         date,
		Type:         orderType,
		Status:       models.ORDER_STATUS_CREATED,
		TeacherId:    teacherId,
		TierId:       tierId,
		PriceHourly:  priceHourly,
		SalaryHourly: salaryHourly,
		CourseId:     courseId,
		ChapterId:    chapterId,
	}

	orderPtr, err := models.CreateOrder(&order)
	if err != nil {
		return nil, err
	}

	return orderPtr, nil
}
