package order

import (
	"time"

	"WolaiWebservice/models"
)

func CreateOrder(userId, teacherId, teacherTier, gradeId, subjectId int64) (int64, *models.Order) {

	var orderType string
	var pricePerHour, realPricePerHour int64
	orderDate := time.Now().Format(time.RFC3339)

	if teacherId != 0 {
		orderType = models.ORDER_TYPE_PERSONAL_INSTANT

		teacher, err := models.ReadTeacherProfile(teacherId)
		if err != nil {
			return 2, nil
		}
		pricePerHour = teacher.PricePerHour
		realPricePerHour = teacher.RealPricePerHour

	} else {
		orderType = models.ORDER_TYPE_GENERAL_INSTANT
	}

	order := models.Order{
		Creator:          userId,
		GradeId:          gradeId,
		SubjectId:        subjectId,
		Date:             orderDate,
		Type:             orderType,
		Status:           models.ORDER_STATUS_CREATED,
		TeacherId:        teacherId,
		PricePerHour:     pricePerHour,
		RealPricePerHour: realPricePerHour,
	}

	orderPtr, err := models.CreateOrder(&order)
	if err != nil {
		return 2, nil
	}

	return 0, orderPtr
}
