package user

import (
	"WolaiWebservice/config/settings"
	"WolaiWebservice/models"
	"errors"
	"time"
)

type initialUserDataUsage struct {
	*models.UserDataUsage
	Freq int64 `json:"freq"`
}

type updateReturn struct {
	Freq int64 `json:"freq"`
}

func GetUserDataUsage(userId int64) (int64, error, *initialUserDataUsage) {
	var err error

	data, err := models.ReadUserDataUsage(userId)
	if err != nil {
		newData := models.UserDataUsage{
			UserId:         userId,
			LastUpdateTime: time.Now(),
		}
		data, err = models.CreateUserDataUsage(&newData)
		if err != nil {
			return 2, err, nil
		}
	}
	freq := settings.FreqSyncDataUsage()

	initialData := initialUserDataUsage{
		UserDataUsage: data,
		Freq:          freq,
	}

	return 0, nil, &initialData
}

func UpdateUserDataUsage(userId, data, dataClass, dataLog int64) (int64, error, *updateReturn) {
	var err error

	dataUsage, err := models.ReadUserDataUsage(userId)
	if err != nil {
		return 2, err, nil
	}

	if dataUsage.Data > data || dataUsage.DataClass > dataClass || dataUsage.DataLog > dataLog {
		return 2, errors.New("更新的流量怎么会小啊！"), nil
	}

	totalClaimAdd := data - dataUsage.Data
	totalClassClaimAdd := dataClass - dataUsage.DataClass
	logDiff := dataLog - dataUsage.DataLog

	logTarget := settings.LogDataTarget()
	totalClaimAdd += logDiff

	if logTarget == models.CONST_CLAIM_TYPE_CLASS {
		totalClassClaimAdd += logDiff
	}

	err = models.HandleDataClaimUpdate(userId, data, dataClass, dataLog, totalClaimAdd, totalClassClaimAdd)

	if err != nil {
		return 2, err, nil
	}

	freq := settings.FreqSyncDataUsage()

	result := updateReturn{
		Freq: freq,
	}

	return 0, nil, &result
}

type ReImbstRecordByMonth struct {
	Month   time.Month                    `json:"month"`
	Year    int64                         `json:"year"`
	Records []*models.DataReimbsmtRecords `json:"records"`
}

func GetReimbstRecords(userId int64, page, count int64) (int64, error, []*ReImbstRecordByMonth) {
	var err error
	records, err := models.QueryUserReimbsmtRecords(userId, page, count)

	if err != nil {
		return 2, err, nil
	}

	resultRecords := make([]*ReImbstRecordByMonth, 0)
	curIndex := 0
	for _, record := range records {
		record.Type = models.ReIMBSMTMap[record.Type]
		record.Status = models.ReIMBSMTMap[record.Status]

		if len(resultRecords) == 0 || resultRecords[curIndex-1].Month != record.CreateTime.Month() || resultRecords[curIndex-1].Year != int64(record.CreateTime.Year()) {
			recordMonth := ReImbstRecordByMonth{
				Month:   record.CreateTime.Month(),
				Year:    int64(record.CreateTime.Year()),
				Records: make([]*models.DataReimbsmtRecords, 0),
			}
			recordMonth.Records = append(recordMonth.Records, record)
			resultRecords = append(resultRecords, &recordMonth)
			curIndex++
		} else {
			resultRecords[curIndex-1].Records = append(resultRecords[curIndex-1].Records, record)
		}
	}
	return 0, nil, resultRecords
}
