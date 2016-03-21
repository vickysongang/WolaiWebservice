// qa_pkg_module
package models

import "github.com/astaxie/beego/orm"

type QaPkgModule struct {
	Id      int64  `json:"id" orm:"pk"`
	Name    string `json:"name"`
	Comment string `json:"comment"`
	Rank    int64  `json:"rank"`
}

func (pkgModule *QaPkgModule) TableName() string {
	return "qa_pkg_module"
}

func init() {
	orm.RegisterModel(new(QaPkgModule))
}

func ReadQaPkgModule(moduleId int64) (*QaPkgModule, error) {
	o := orm.NewOrm()
	module := QaPkgModule{Id: moduleId}
	err := o.Read(&module)
	if err != nil {
		return nil, err
	}
	return &module, nil
}
