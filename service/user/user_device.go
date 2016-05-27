// user_device
package user

import (
	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
)

func QueryIosUserDevices() []models.UserDevice {
	o := orm.NewOrm()
	var devices []models.UserDevice
	o.QueryTable(new(models.UserDevice).TableName()).
		Filter("device_type", "ios").Filter("voip_token__isnull", false).
		All(&devices)
	return devices
}
