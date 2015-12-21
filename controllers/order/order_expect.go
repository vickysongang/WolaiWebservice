package order

import (
	"errors"

	"WolaiWebservice/models"
)

// TODO: 查询最低价的档次
const (
	LOWEST_TEACHER_TIER = 3
)

type priceInfo struct {
	Price int64 `json:"price"`
}

func CalculateOrderExpect(userId, teacherId, teacherTier, gradeId, subjectId int64) (int64, error, *priceInfo) {
	_, err := models.ReadUser(userId)
	if err != nil {
		return 2, errors.New("用户资料异常"), nil
	}

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
		tier, err := models.ReadTeacherTierHourly(LOWEST_TEACHER_TIER)
		if err != nil {
			return 2, errors.New("导师资料异常"), nil
		}

		price = tier.QAPriceHourly
	}

	price = ((price / 10) / 60) * 10

	info := priceInfo{
		Price: price,
	}

	return 0, nil, &info
}
