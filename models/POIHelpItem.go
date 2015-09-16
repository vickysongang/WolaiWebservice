// POIHelpItem.go
package models

import (
	"POIWolaiWebService/utils"

	"github.com/astaxie/beego/orm"
)

type POIHelpItem struct {
	Id    int64  `json:"-" orm:"pk"`
	Title string `json:"title"`
	Url   string `json:"url"`
	Rank  int64  `json:"-"`
}

type POIHelpItems []*POIHelpItem

func (help *POIHelpItem) TableName() string {
	return "help_item"
}

func init() {
	orm.RegisterModel(new(POIHelpItem))
}

func QueryHelpItems() (POIHelpItems, error) {
	helpItems := make(POIHelpItems, 0)
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(utils.DB_TYPE)
	qb.Select("id,title,url,rank").From("help_item").OrderBy("rank").Asc()
	sql := qb.String()
	_, err := o.Raw(sql).QueryRows(&helpItems)
	if err != nil {
		return nil, err
	}
	return helpItems, nil
}
