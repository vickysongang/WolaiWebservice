// POIInvitation
package models

import "github.com/astaxie/beego/orm"

type POIInvitation struct {
	Id             int64 `orm:"pk"`
	Name           string
	Pid            int64
	InvitationCode string
}

type POIUserToInvitation struct {
	Id             int64  `json:"id" orm:"pk"`
	UserId         int64  `json:"userId"`
	InvitationCode string `json:"invitationCode"`
}

func (userToInvitation *POIUserToInvitation) TableName() string {
	return "user_to_invitation"
}

func (invitation *POIInvitation) TableName() string {
	return "invitation"
}

func init() {
	orm.RegisterModel(new(POIUserToInvitation), new(POIInvitation))
}

/*
 * 检查验证码是否有效
 */
func CheckInvitationCodeValid(invitationCode string) bool {
	o := orm.NewOrm()
	count, err := o.QueryTable("invitation").Filter("invitation_code", invitationCode).Count()
	if err != nil {
		return false
	}
	if count == 0 {
		return false
	}
	return true
}

/*
 * 判断用户是否绑定过邀请码
 */
func CheckUserHasBindWithInvitationCode(userId int64) bool {
	o := orm.NewOrm()
	count, err := o.QueryTable("user_to_invitation").Filter("user_id", userId).Count()
	if err != nil {
		return false
	}
	if count > 0 {
		return true
	}
	return false
}

/*
 * 将用户和邀请码绑定
 */
func InsertUserToInvitation(userToInvitation *POIUserToInvitation) (*POIUserToInvitation, error) {
	o := orm.NewOrm()
	id, err := o.Insert(userToInvitation)
	if err != nil {
		return nil, err
	}
	userToInvitation.Id = id
	return userToInvitation, nil
}
