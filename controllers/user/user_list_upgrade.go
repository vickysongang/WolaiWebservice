// user_list_upgrade
package user

import (
	userService "WolaiWebservice/service/user"
	"WolaiWebservice/websocket"
)

func GetTeacherRecommendationUpgrade(userId, page, count int64) (int64, error, []*UserListItem) {
	resultTeacherIds := make([]int64, 0)
	result := make([]*UserListItem, 0)
	freeTeacherIds := websocket.UserManager.GetFreeTeachers()
	busyTeacherIds := websocket.UserManager.GetBusyTeachers()
	onlineTeacherIds := websocket.UserManager.GetOnlineTeachers(true)

	var tempTeacheherIds []int64
	tempTeacheherIds = append(tempTeacheherIds, freeTeacherIds...)
	tempTeacheherIds = append(tempTeacheherIds, busyTeacherIds...)
	tempTeacheherIds = append(tempTeacheherIds, onlineTeacherIds...)

	tempTeacherIdsLen := int64(len(tempTeacheherIds))
	if tempTeacherIdsLen/((page+1)*count) > 0 {
		resultTeacherIds = append(resultTeacherIds, tempTeacheherIds[(page*count):(page+1)*count]...)
	} else {
		leftTempLen := tempTeacherIdsLen % count
		tempPageSize := tempTeacherIdsLen / count
		if leftTempLen > 0 {
			var offset, limitCount int64
			if page == tempPageSize {
				offset = 0
				limitCount = count - leftTempLen
				resultTeacherIds = append(resultTeacherIds, tempTeacheherIds[(page*count):]...)
			} else {
				offset = count - leftTempLen + (page-tempPageSize-1)*count
				limitCount = count
			}
			teacherIds, err := userService.QueryTeacherRecommendationExcludeOnline(userId, limitCount, offset, tempTeacheherIds)
			if err != nil {
				return 0, nil, result
			}
			resultTeacherIds = append(resultTeacherIds, teacherIds...)
		} else {
			offset := (page - tempPageSize) * count
			teacherIds, err := userService.QueryTeacherRecommendationExcludeOnline(userId, count, offset, tempTeacheherIds)
			if err != nil {
				return 0, nil, result
			}
			resultTeacherIds = append(resultTeacherIds, teacherIds...)
		}
	}

	for _, teacherId := range resultTeacherIds {
		item, err := AssembleUserListItemUpgrade(teacherId)
		if err != nil {
			continue
		}

		result = append(result, item)
	}

	return 0, nil, result
}
