// session
package session

import (
	"WolaiWebservice/models"

	"github.com/astaxie/beego/orm"
)

func GetUserSessions(userId, page, count int64) ([]*models.Session, error) {
	o := orm.NewOrm()
	cond := orm.NewCondition()
	cond = cond.Or("creator", userId).Or("tutor", userId)
	var sessions []*models.Session
	_, err := o.QueryTable("sessions").SetCond(cond).
		OrderBy("-id").Offset(page * count).
		Limit(count).All(&sessions)
	return sessions, err
}
