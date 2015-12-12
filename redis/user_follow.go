package redis

import (
	"strconv"

	"gopkg.in/redis.v3"

	"WolaiWebservice/models"
)

func (rm *POIRedisManager) SetUserFollow(userId, followId int64) {
	userIdStr := strconv.FormatInt(userId, 10)
	followIdStr := strconv.FormatInt(followId, 10)

	_ = rm.RedisClient.HSet(USER_FOLLOWING+userIdStr, followIdStr, "0")
	_ = rm.RedisClient.HSet(USER_FOLLOWER+followIdStr, userIdStr, "0")

	return
}

func (rm *POIRedisManager) RemoveUserFollow(userId, followId int64) {
	userIdStr := strconv.FormatInt(userId, 10)
	followIdStr := strconv.FormatInt(followId, 10)

	_ = rm.RedisClient.HDel(USER_FOLLOWING+userIdStr, followIdStr)
	_ = rm.RedisClient.HDel(USER_FOLLOWER+followIdStr, userIdStr)

	return
}

func (rm *POIRedisManager) HasFollowedUser(userId, followId int64) bool {
	userIdStr := strconv.FormatInt(userId, 10)
	followIdStr := strconv.FormatInt(followId, 10)

	var result bool
	_, err := rm.RedisClient.HGet(USER_FOLLOWING+userIdStr, followIdStr).Result()
	if err == redis.Nil {
		result = false
	} else {
		result = true
	}

	return result
}

func (rm *POIRedisManager) GetUserFollowList(userId, pageNum, pageCount int64) models.POITeachers {
	userIdStr := strconv.FormatInt(userId, 10)
	userIds := rm.RedisClient.HKeys(USER_FOLLOWING + userIdStr).Val()
	start := pageNum * pageCount
	//	teachers := make(POITeachers, len(userIds))
	teachers := make(models.POITeachers, 0)
	//	for i := range userIds {
	//		userIdtmp, _ := strconv.ParseInt(userIds[i], 10, 64)
	//		teachers[i] = *(QueryTeacher(userIdtmp))
	//	}
	length := int64(len(userIds))
	for i := start; i < (start + pageCount); i++ {
		if i < length {
			userIdtmp, _ := strconv.ParseInt(userIds[i], 10, 64)
			teacher := *(models.QueryTeacher(userIdtmp))
			teacher.HasFollowed = true
			teachers = append(teachers, teacher)
		}
	}
	return teachers
}
