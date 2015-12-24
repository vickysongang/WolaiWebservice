package user

import (
	"errors"
	"fmt"

	"github.com/astaxie/beego/orm"

	"WolaiWebservice/models"
)

type greetingInfo struct {
	Greeting string `json:"greeting"`
}

func AssembleUserGreeting(userId int64) (int64, error, *greetingInfo) {
	var info greetingInfo

	user, err := models.ReadUser(userId)
	if err != nil {
		return 2, errors.New("用户信息异常"), nil
	}

	if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER {
		profile, err := models.ReadTeacherProfile(userId)
		if err != nil {
			return 2, errors.New("用户信息异常"), nil
		}

		hours := profile.ServiceTime / 3600
		if hours == 0 {
			info.Greeting = "我来和你一起开启名师之旅"
		} else {
			info.Greeting = fmt.Sprintf("你和我来已经一起努力了%d个小时", hours)
		}
	} else if user.AccessRight == models.USER_ACCESSRIGHT_STUDENT {
		o := orm.NewOrm()

		var sessions []models.Session
		_, err = o.QueryTable("sessions").Filter("creator", userId).Filter("status", models.SESSION_STATUS_COMPLETE).
			All(&sessions)
		if err != nil {
			info.Greeting = "“我来”陪你开启成长之旅"
		} else {
			var sum int64
			for _, session := range sessions {
				sum = sum + session.Length
			}

			hours := sum / 3600
			if hours == 0 {
				info.Greeting = "“我来”陪你开启成长之旅"
			} else {
				info.Greeting = fmt.Sprintf("“我来”已经陪伴你%d小时了", hours)
			}
		}
	}

	return 0, nil, &info
}
