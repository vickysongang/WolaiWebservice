// qa_pkg
package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

const (
	QA_PKG_TYPE_PERMANENT = "permanent"
	QA_PKG_TYPE_MONTHLY   = "monthly"
)

type QaPkg struct {
	Id            int64     `json:"id" orm:"pk"`
	TimeLength    int64     `json:"timeLength"`
	OriginalPrice int64     `json:"originalPrice"`
	DiscountPrice int64     `json:"discountPrice"`
	Type          string    `json:"type"`
	Month         int64     `json:"month"`
	Comment       string    `json:"comment"`
	CreateTime    time.Time `json:"createTime" orm:"auto_now_add;type(datetime)"`
	ModuleId      int64     `json:"-"`
}

func (pkg *QaPkg) TableName() string {
	return "qa_pkg"
}

func init() {
	orm.RegisterModel(new(QaPkg))
}

func ReadQaPkg(pkgId int64) (*QaPkg, error) {
	o := orm.NewOrm()
	pkg := QaPkg{Id: pkgId}
	err := o.Read(&pkg)
	if err != nil {
		return nil, err
	}
	return &pkg, nil
}
