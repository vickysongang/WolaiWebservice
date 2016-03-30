package order

import (
	"errors"

	"WolaiWebservice/models"
	qapkgService "WolaiWebservice/service/qapkg"
)

type priceInfo struct {
	Price            int64 `json:"price"`
	LeftQaTimeLength int64 `json:"leftQaTimeLength"`
}

func CalculateOrderExpect(userId, teacherId, teacherTier, gradeId, subjectId int64) (int64, error, *priceInfo) {
	var price int64
	if teacherId != 0 {
		teacher, err := models.ReadTeacherProfile(teacherId)
		if err != nil {
			return 2, errors.New("导师资料异常"), nil
		}

		tier, err := models.ReadTeacherTierHourly(teacher.TierId)
		if err != nil {
			return 2, errors.New("导师资料异常"), nil
		}

		price = tier.QAPriceHourly
	} else if teacherTier != 0 {
		tier, err := models.ReadTeacherTierHourly(teacherTier)
		if err != nil {
			return 2, errors.New("导师资料异常"), nil
		}

		price = tier.QAPriceHourly
	} else {
		tier, err := models.ReadTeacherTierHourly(models.LOWEST_TEACHER_TIER)
		if err != nil {
			return 2, errors.New("导师资料异常"), nil
		}

		price = tier.QAPriceHourly
	}

	price = ((price / 10) / 60) * 10
	leftQaTimeLength := qapkgService.GetLeftQaTimeLength(userId)

	info := priceInfo{
		Price:            price,
		LeftQaTimeLength: leftQaTimeLength,
	}

	return 0, nil, &info
}
