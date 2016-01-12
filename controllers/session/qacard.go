package session

import (
	"WolaiWebservice/models"
	sessionService "WolaiWebservice/service/session"
)

func QACardCatelog(pid int64) (int64, error, []*models.QACardCatelog) {
	var err error

	catelogs, err := sessionService.QueryQACardCatelog(pid)
	if err != nil {
		return 2, err, nil
	}

	return 0, nil, catelogs
}

func QACardAttach(catelogId int64) (int64, error, []*models.QACardAttach) {
	var err error

	attachs, err := sessionService.QueryQACardAttach(catelogId)
	if err != nil {
		return 2, err, nil
	}

	return 0, nil, attachs
}
