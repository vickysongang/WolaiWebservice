package message

import (
	"WolaiWebservice/models"
	"WolaiWebservice/utils/leancloud"
)

func GetConversation(userId1, userId2 int64) (int64, string) {
	_, err1 := models.ReadUser(userId1)
	_, err2 := models.ReadUser(userId2)
	if err1 != nil || err2 != nil {
		return 2, ""
	}

	return 0, leancloud.GetConversation(userId1, userId2)
}
