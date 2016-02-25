package user

import (
	"WolaiWebservice/models"
	userService "WolaiWebservice/service/user"
	"WolaiWebservice/websocket"
)

func SearchUser(userId int64, keyword string, page, count int64) (int64, error, []*UserListItem) {
	var err error

	result := make([]*UserListItem, 0)

	users, err := userService.QueryUserByKeyword(keyword, page, count)
	if err != nil {
		return 0, nil, result
	}

	for _, user := range users {
		item, err := AssembleUserListItem(user.Id)
		if err != nil {
			continue
		}

		result = append(result, item)
	}

	return 0, nil, result
}

func GetTeacherRecent(userId, page, count int64) (int64, error, []*UserListItem) {
	var err error

	result := make([]*UserListItem, 0)

	teacherIds, err := userService.QueryTeacherBySessionFreq(userId, page, count)
	if err != nil {
		return 0, nil, result
	}

	for _, teacherId := range teacherIds {
		item, err := AssembleUserListItem(teacherId)
		if err != nil {
			continue
		}

		result = append(result, item)
	}

	return 0, nil, result
}

func GetTeacherRecommendation(userId, page, count int64) (int64, error, []*UserListItem) {
	var err error

	result := make([]*UserListItem, 0)

	onlineTeacherIds := websocket.UserManager.GetOnlineTeachers()
	teacherIds, err := userService.QueryTeacherRecommendationExcludeOnline(userId, page, count, onlineTeacherIds)
	if err != nil {
		return 0, nil, result
	}
	resultTeacherIds := make([]int64, 0)
	if page == 0 {
		resultTeacherIds = append(resultTeacherIds, onlineTeacherIds...)
		resultTeacherIds = append(resultTeacherIds, teacherIds...)
	} else {
		resultTeacherIds = teacherIds
	}
	for _, teacherId := range resultTeacherIds {
		item, err := AssembleUserListItem(teacherId)
		if err != nil {
			continue
		}

		result = append(result, item)
	}

	return 0, nil, result
}

func GetContactRecommendation(userId, page, count int64) (int64, error, []*UserListItem) {
	var err error

	result := make([]*UserListItem, 0)

	// 如果是第一页，加入我来团队和助教
	if page == 0 {
		wolaiItem, err := AssembleUserListItem(models.USER_WOLAI_SUPPORT)
		if err == nil {
			result = append(result, wolaiItem)
		}

		assistants, err := userService.QueryUserByAccessRight(models.USER_ACCESSRIGHT_ASSISTANT, 0, 10)
		if err == nil {
			for _, assistant := range assistants {
				assistItem, err := AssembleUserListItem(assistant.Id)
				if err != nil {
					continue
				}

				result = append(result, assistItem)
			}
		}

	}

	teacherIds, err := userService.QueryTeacherRecommendation(userId, page, count)
	if err != nil {
		return 0, nil, result
	}

	for _, teacherId := range teacherIds {
		item, err := AssembleUserListItem(teacherId)
		if err != nil {
			continue
		}

		result = append(result, item)
	}

	return 0, nil, result
}
