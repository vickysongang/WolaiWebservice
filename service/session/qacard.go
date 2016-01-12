package session

import (
	"errors"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"

	"WolaiWebservice/models"
)

func QueryQACardCatalog(pid int64) ([]*models.QACardCatalog, error) {
	var err error

	o := orm.NewOrm()

	var catalogs []*models.QACardCatalog
	_, err = o.QueryTable(new(models.QACardCatalog).TableName()).
		Filter("pid", pid).
		All(&catalogs)
	if err != nil {
		seelog.Error("%s | Pid: %d", err.Error(), pid)
		return nil, errors.New("无法找到子节点")
	}

	return catalogs, nil
}

func QueryQACardAttach(catalogId int64) ([]*models.QACardAttach, error) {
	var err error

	o := orm.NewOrm()

	var attachs []*models.QACardAttach
	_, err = o.QueryTable(new(models.QACardAttach).TableName()).
		Filter("catalog_id", catalogId).
		OrderBy("rank").
		All(&attachs)
	if err != nil {
		seelog.Error("%s | CatalogId: %d", err.Error(), catalogId)
		return nil, errors.New("无法找到对应附件")
	}

	return attachs, nil
}
