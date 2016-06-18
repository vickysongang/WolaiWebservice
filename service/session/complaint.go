// complaint
package session

import (
	"WolaiWebservice/config"

	"github.com/astaxie/beego/orm"
)

func GetComplaintStatus(userId, sessionId int64) string {
	o := orm.NewOrm()

	qb, _ := orm.NewQueryBuilder(config.Env.Database.Type)
	status := ""
	qb.Select("status").From("complaint").Where("user_id = ? and session_id = ?")
	sql := qb.String()
	o.Raw(sql, userId, sessionId).QueryRow(&status)
	return status
}
