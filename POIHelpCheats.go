// POIHelpCheats.go
package main

import (
	"github.com/astaxie/beego/orm"
)

type POIHelpCheats struct {
	Id    int64  `json:"-" orm:"pk"`
	Title string `json:"title"`
	Url   string `json:"url"`
	Rank  int64  `json:"-"`
}

type POIHelpCheatses []*POIHelpCheats

func (help *POIHelpCheats) TableName() string {
	return "help_cheats"
}

func init() {
	orm.RegisterModel(new(POIHelpCheats))
}

func QueryHelpCheats() (POIHelpCheatses, error) {
	helpCheatses := make(POIHelpCheatses, 0)
	o := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder(DB_TYPE)
	qb.Select("id,title,url,rank").From("help_cheats").OrderBy("rank").Asc()
	sql := qb.String()
	_, err := o.Raw(sql).QueryRows(&helpCheatses)
	if err != nil {
		return nil, err
	}
	return helpCheatses, nil
}
