package user

import (
	"WolaiWebservice/models"
	userService "WolaiWebservice/service/user"
	"WolaiWebservice/websocket"
)

type UserListItem struct {
	Id           int64    `json:"id"`
	Nickname     string   `json:"nickname"`
	Avatar       string   `json:"avatar"`
	Gender       int64    `json:"gender"`
	AccessRight  int64    `json:"accessRight"`
	School       string   `json:"school"`
	Major        string   `json:"major"`
	SubjectList  []string `json:"subjectList,omitempty"`
	OnlineStatus string   `json:"onlineStatus,omitempty"`
}

func AssembleUserListItem(userId int64) (*UserListItem, error) {
	var err error

	user, err := models.ReadUser(userId)
	if err != nil {
		return nil, err
	}

	item := UserListItem{
		Id:           user.Id,
		Nickname:     user.Nickname,
		Avatar:       user.Avatar,
		Gender:       user.Gender,
		AccessRight:  user.AccessRight,
		OnlineStatus: websocket.WsManager.GetUserStatus(user.Id),
	}

	if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER {
		profile, err := models.ReadTeacherProfile(userId)
		if err != nil {
			return nil, err
		}

		school, err := models.ReadSchool(profile.SchoolId)
		if err == nil {
			item.School = school.Name
		}

		item.Major = profile.Major

		subjects, err := userService.GetTeacherSubjectNameSlice(userId)
		if err == nil {
			item.SubjectList = subjects
		}
	}

	return &item, nil
}
