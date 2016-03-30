// dictionary
package models

import "github.com/astaxie/beego/orm"

const (
	DICT_TYPE_EVALUATION = "evaluation"
)

type Dictionary struct {
	Id      int64  `json:"id" orm:"pk"`
	Code    string `json:"code"`
	Meaning string `json:"meaning"`
	Type    string `json:"type"`
}

func init() {
	orm.RegisterModel(new(Dictionary))
}

func (d *Dictionary) TableName() string {
	return "dictionary"
}

func ReadDictionary(id int64) (*Dictionary, error) {
	o := orm.NewOrm()
	dictionary := Dictionary{Id: id}
	err := o.Read(&dictionary)
	if err != nil {
		return nil, err
	}
	return &dictionary, nil
}
