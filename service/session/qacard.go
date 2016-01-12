package session

import (
	"errors"

	"github.com/astaxie/beego/orm"
	"github.com/cihub/seelog"

	"WolaiWebservice/models"
)

func QueryQACardCatelog(pid int64) ([]*models.QACardCatelog, error) {
	var err error

	o := orm.NewOrm()

	var catelogs []*models.QACardCatelog
	_, err = o.QueryTable(new(models.QACardCatelog).TableName()).
		Filter("pid", pid).
		All(&catelogs)
	if err != nil {
		seelog.Error("%s | Pid: %d", err.Error(), pid)
		return nil, errors.New("无法找到子节点")
	}

	return catelogs, nil
}

func QueryQACardAttach(catelogId int64) ([]*models.QACardAttach, error) {
	var err error

	o := orm.NewOrm()

	var attachs []*models.QACardAttach
	_, err = o.QueryTable(new(models.QACardAttach).TableName()).
		Filter("catelog_id", catelogId).
		OrderBy("rank").
		All(&attachs)
	if err != nil {
		seelog.Error("%s | CatelogId: %d", err.Error(), catelogId)
		return nil, errors.New("无法找到对应附件")
	}

	return attachs, nil
}
