// user_list_upgrade
package user

import (
	userService "WolaiWebservice/service/user"
	"WolaiWebservice/websocket"
	"sort"
)

type IntArray []int64

func (array IntArray) Len() int {
	return len(array)
}

func (array IntArray) Less(i, j int) bool {
	if array[i] < array[j] {
		return true
	}
	return false
}

func (array IntArray) Swap(i, j int) {
	var temp int64 = array[i]
	array[i] = array[j]
	array[j] = temp
}

func GetTeacherRecommendationUpgrade(userId, page, count int64) (int64, error, []*UserListItem) {
	resultTeacherIds := make([]int64, 0)
	result := make([]*UserListItem, 0)
	freeTeacherIds := websocket.UserManager.GetFreeTeachers()
	sort.Sort(IntArray(freeTeacherIds))
	busyTeacherIds := websocket.UserManager.GetBusyTeachers()
	sort.Sort(IntArray(busyTeacherIds))
	onlineTeacherIds := websocket.UserManager.GetOnlineTeachers(true)
	sort.Sort(IntArray(onlineTeacherIds))

	var tempTeacherIds []int64
	tempTeacherIds = append(tempTeacherIds, freeTeacherIds...)
	tempTeacherIds = append(tempTeacherIds, busyTeacherIds...)
	tempTeacherIds = append(tempTeacherIds, onlineTeacherIds...)

	tempTeacherIdsLen := int64(len(tempTeacherIds))
	if tempTeacherIdsLen/((page+1)*count) > 0 {
		resultTeacherIds = append(resultTeacherIds, tempTeacherIds[(page*count):(page+1)*count]...)
	} else {
		leftTempLen := tempTeacherIdsLen % count
		tempPageSize := tempTeacherIdsLen / count
		if leftTempLen > 0 {
			var offset, limitCount int64
			if page == tempPageSize {
				offset = 0
				limitCount = count - leftTempLen
				resultTeacherIds = append(resultTeacherIds, tempTeacherIds[(page*count):]...)
			} else {
				offset = count - leftTempLen + (page-tempPageSize-1)*count
				limitCount = count
			}
			teacherIds, err := userService.QueryTeacherRecommendationExcludeOnline(userId, limitCount, offset, tempTeacherIds)
			if err != nil {
				return 0, nil, result
			}
			resultTeacherIds = append(resultTeacherIds, teacherIds...)
		} else {
			offset := (page - tempPageSize) * count
			teacherIds, err := userService.QueryTeacherRecommendationExcludeOnline(userId, count, offset, tempTeacherIds)
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

func GetTeacherRecentUpgrade(userId, page, count int64) (int64, error, []*UserListItem) {
	var err error

	result := make([]*UserListItem, 0)

	teacherIds, err := userService.QueryTeacherBySessionFreq(userId, page, count)
	if err != nil {
		return 0, nil, result
	}

	for _, teacherId := range teacherIds {
		item, err := AssembleUserListItemUpgrade(teacherId)
		if err != nil {
			continue
		}

		result = append(result, item)
	}

	return 0, nil, result
}
