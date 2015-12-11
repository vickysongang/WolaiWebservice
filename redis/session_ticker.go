package redis

import (
	"encoding/json"
	"strconv"
	"time"

	seelog "github.com/cihub/seelog"
	"gopkg.in/redis.v3"

	"WolaiWebservice/models"
)

func (rm *POIRedisManager) SetSessionTicker(timestamp int64, tickerInfo string) {
	tickerZ := redis.Z{Member: tickerInfo, Score: float64(timestamp)}

	_ = rm.RedisClient.ZAdd(SESSION_TICKER, tickerZ)
}

func (rm *POIRedisManager) GetSessionTicks(timestamp int64) []string {
	ticks, err := rm.RedisClient.ZRangeByScore(SESSION_TICKER,
		redis.ZRangeByScore{
			Min:    "-inf",
			Max:    strconv.FormatInt(timestamp, 10),
			Offset: 0,
			Count:  10,
		}).Result()
	if err == redis.Nil {
		return nil
	}

	for i := range ticks {
		_ = rm.RedisClient.ZRem(SESSION_TICKER, ticks[i])
	}

	return ticks
}

/*
 * 将老师的计划开始时间和预计结束时间存入redis
 */
func (rm *POIRedisManager) SetSessionUserTick(sessionId int64) bool {
	session, err := models.ReadSession(sessionId)
	if err != nil {
		return false
	}

	_, err = models.ReadOrder(session.OrderId)
	if err != nil {
		return false
	}

	planTimeStr := session.PlanTime
	planTime, _ := time.Parse(time.RFC3339, planTimeStr)
	blockDuration := 30 * time.Minute

	timeFrom := planTime.Add(-blockDuration)

	teacherStartMap := map[string]int64{
		"userId":    session.Tutor,
		"sessionId": sessionId,
		"lock":      1,
	}
	studentStartMap := map[string]int64{
		"userId":    session.Creator,
		"sessionId": sessionId,
		"lock":      1,
	}

	teacherStartStr, _ := json.Marshal(teacherStartMap)
	studentStartStr, _ := json.Marshal(studentStartMap)

	teacherIdStr := strconv.FormatInt(session.Tutor, 10)
	studentIdStr := strconv.FormatInt(session.Creator, 10)

	teacherTimeFromZ := redis.Z{Member: string(teacherStartStr), Score: float64(timeFrom.Unix())}
	studentTimeFromZ := redis.Z{Member: string(studentStartStr), Score: float64(timeFrom.Unix())}

	rm.RedisClient.ZAdd(SESSION_USER_LOCK+teacherIdStr, teacherTimeFromZ)
	rm.RedisClient.ZAdd(SESSION_USER_LOCK+studentIdStr, studentTimeFromZ)

	rm.RedisClient.ZAdd(SESSION_USER_TICKER, teacherTimeFromZ)
	rm.RedisClient.ZAdd(SESSION_USER_TICKER, studentTimeFromZ)

	seelog.Debug("SetSessionLock: sessionId:", sessionId, "teacherId:", session.Tutor, " studentId:", session.Creator)

	if time.Now().Unix() > timeFrom.Unix() {
		return true
	}
	return false
}

/*
 * 获取特定时间段内的用户事件
 */
func (rm *POIRedisManager) GetSessionUserTicks(timestamp int64) []models.POITickInfo {
	ticks, err := rm.RedisClient.ZRangeByScoreWithScores(SESSION_USER_TICKER,
		redis.ZRangeByScore{
			Min:    "-inf",
			Max:    strconv.FormatInt(timestamp+5, 10),
			Offset: 0,
			Count:  0,
		}).Result()
	if err == redis.Nil {
		return nil
	}

	tickInfo := make([]models.POITickInfo, 0)
	for i := range ticks {
		_ = rm.RedisClient.ZRem(SESSION_USER_TICKER, ticks[i].Member.(string))
		tickInfo = append(tickInfo, models.POITickInfo{
			Timestamp: int64(ticks[i].Score),
			Content:   ticks[i].Member.(string),
		})
	}

	return tickInfo
}

/*
 * 判断老师在某一时间段内是否处于忙碌状态
 */
func (rm *POIRedisManager) IsUserAvailable(userId int64, startTime time.Time) bool {
	blockDuration := 30 * time.Minute
	timeFrom := startTime.Add(-blockDuration)
	timestampFrom := timeFrom.Unix()

	userIdStr := strconv.FormatInt(userId, 10)
	items, err := rm.RedisClient.ZRangeByScore(SESSION_USER_LOCK+userIdStr,
		redis.ZRangeByScore{
			Min:    strconv.FormatInt(timestampFrom-125, 10),
			Max:    strconv.FormatInt(timestampFrom+125, 10),
			Offset: 0,
			Count:  10,
		}).Result()
	if err == redis.Nil {
		return true
	}
	if len(items) > 0 {
		return false
	}
	return true
}
