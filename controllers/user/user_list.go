package user

import (
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

	onlineTeacherIds := websocket.UserManager.GetOnlineTeachers(false)
	teacherIds, err := userService.QueryTeacherRecommendationExcludeOnline(userId, count, page*count, onlineTeacherIds)
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
	assistants, err := userService.QueryAssistants()
	if page == 0 {
		if err == nil {
			for _, userId := range assistants {
				assistItem, err := AssembleUserListItem(userId)
				if err != nil {
					continue
				}
				result = append(result, assistItem)
			}
		}

	}

	teacherIds, err := userService.QueryTeacherRecommendation(userId, assistants, page, count)
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
