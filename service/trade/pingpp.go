// pingpp
package trade

import (
	"WolaiWebservice/models"

	"github.com/astaxie/beego/orm"
)

func QueryPingppRecordByChargeId(chargeId string) (*models.PingppRecord, error) {
	o := orm.NewOrm()
	record := models.PingppRecord{}
	err := o.QueryTable("pingpp_record").Filter("charge_id", chargeId).One(&record)
	return &record, err
}
