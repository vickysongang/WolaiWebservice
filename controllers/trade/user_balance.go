// user_balance
package trade

import "WolaiWebservice/models"

func GetUserBalance(userId int64) (map[string]int64, error) {
	result := make(map[string]int64)
	user, err := models.ReadUser(userId)
	if err != nil {
		return result, err
	}
	profile, _ := models.ReadStudentProfile(userId)
	var gradeId, quotaQuantity, quotaGradeId int64
	if profile != nil {
		gradeId = profile.GradeId
		quotaQuantity = profile.QuotaQuantity
		quotaGradeId = profile.QuotaGradeId
	}
	result["balance"] = user.Balance
	result["gradeId"] = gradeId
	result["quotaQuantity"] = quotaQuantity
	result["quotaGradeId"] = quotaGradeId
	return result, nil
}
