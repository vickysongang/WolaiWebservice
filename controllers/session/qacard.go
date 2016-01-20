package session

import (
	"WolaiWebservice/models"
	sessionService "WolaiWebservice/service/session"
)

func QACardCatalog(pid int64) (int64, error, []*models.QACardCatalog) {
	var err error

	catalogs, err := sessionService.QueryQACardCatalog(pid)
	if err != nil {
		return 2, err, nil
	}

	return 0, nil, catalogs
}

func QACardAttach(catalogId int64) (int64, error, []*models.QACardAttach) {
	var err error

	attachs, err := sessionService.QueryQACardAttach(catalogId)
	if err != nil {
		return 2, err, nil
	}

	return 0, nil, attachs
}
