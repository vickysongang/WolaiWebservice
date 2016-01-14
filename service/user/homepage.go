package user

import (
	"errors"
	"fmt"

	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
)

func GetUserBroadcast(userId int64) ([]*models.Broadcast, error) {
	o := orm.NewOrm()

	user, err := models.ReadUser(userId)
	if err != nil {
		return nil, err
	}

	var broadcasts []*models.Broadcast
	_, err = o.QueryTable(new(models.Broadcast).TableName()).
		Filter("access_right", user.AccessRight).
		Filter("active", "Y").
		OrderBy("rank").All(&broadcasts)

	if err != nil {
		return nil, errors.New("获取通知信息失败")
	}

	return broadcasts, nil
}

func AssembleUserGreeting(userId int64) (*string, error) {
	var err error

	var greeting string

	user, err := models.ReadUser(userId)
	if err != nil {
		return nil, err
	}

	if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER {
		profile, err := models.ReadTeacherProfile(userId)
		if err != nil {
			greeting = "我来和你一起开启名师之旅"
		} else {
			hours := profile.ServiceTime / 3600
			if hours == 0 {
				greeting = "我来和你一起开启名师之旅"
			} else {
				greeting = fmt.Sprintf("你和我来已经一起努力了%d个小时", hours)
			}
		}
	} else if user.AccessRight == models.USER_ACCESSRIGHT_STUDENT {
		o := orm.NewOrm()

		var sessions []models.Session
		_, err = o.QueryTable(new(models.Session).TableName()).
			Filter("creator", userId).
			Filter("status", models.SESSION_STATUS_COMPLETE).
			All(&sessions)
		if err != nil {
			greeting = "“我来”陪你开启成长之旅"
		} else {
			var sum int64
			for _, session := range sessions {
				sum = sum + session.Length
			}

			hours := sum / 3600
			if hours == 0 {
				greeting = "“我来”陪你开启成长之旅"
			} else {
				greeting = fmt.Sprintf("“我来”已经陪伴你%d小时了", hours)
			}
		}
	}

	return &greeting, err
}
