package order

import (
	"time"

	"WolaiWebservice/models"
	"WolaiWebservice/websocket"
)

func CreateOrder(userId, teacherId, teacherTier, gradeId, subjectId int64) (int64, *models.Order) {

	var orderType string
	var priceHourly, salaryHourly int64
	orderDate := time.Now().Format(time.RFC3339)

	if teacherId != 0 {
		orderType = models.ORDER_TYPE_PERSONAL_INSTANT

		teacher, err := models.ReadTeacherProfile(teacherId)
		if err != nil {
			return 2, nil
		}

		tier, err := models.ReadTeacherTierHourly(teacher.TierId)
		if err != nil {
			return 2, nil
		}

		priceHourly = tier.QAPriceHourly
		salaryHourly = tier.QASalaryHourly
	} else {
		orderType = models.ORDER_TYPE_GENERAL_INSTANT
	}

	order := models.Order{
		Creator:      userId,
		GradeId:      gradeId,
		SubjectId:    subjectId,
		Date:         orderDate,
		Type:         orderType,
		Status:       models.ORDER_STATUS_CREATED,
		TeacherId:    teacherId,
		PriceHourly:  priceHourly,
		SalaryHourly: salaryHourly,
	}

	orderPtr, err := models.CreateOrder(&order)
	if err != nil {
		return 2, nil
	}

	if orderPtr.Type == models.ORDER_TYPE_PERSONAL_INSTANT {
		go websocket.InitOrderMonitor(orderPtr.Id, teacherId)
	}

	return 0, orderPtr
}
