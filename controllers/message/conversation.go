package message

import (
	"WolaiWebservice/models"
	"WolaiWebservice/utils/leancloud"
)

func GetConversation(userId1, userId2 int64) (int64, error, string) {
	var err error

	_, err = models.ReadUser(userId1)
	if err != nil {
		return 2, err, ""
	}

	_, err = models.ReadUser(userId2)
	if err != nil {
		return 2, err, ""
	}

	convId, err := leancloud.GetConversation(userId1, userId2)
	if err != nil {
		return 2, err, ""
	}

	return 0, nil, convId
}
