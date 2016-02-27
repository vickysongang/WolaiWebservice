package rpc

import (
	"strconv"

	"WolaiWebservice/handlers/response"
	"WolaiWebservice/models"
	"WolaiWebservice/websocket"
)

func (watcher *RpcWatcher) GetStatusLive(request *RpcRequest, resp *RpcResponse) error {
	allOnlineUsers := len(websocket.UserManager.OnlineUserMap)
	onlineStudentsCount := 0
	onlineTeachersCount := 0
	for userId, _ := range websocket.UserManager.OnlineUserMap {
		user, _ := models.ReadUser(userId)
		if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER {
			onlineTeachersCount++
		}
	}
	onlineStudentsCount = allOnlineUsers - onlineTeachersCount
	liveTeachersCount := len(websocket.TeacherManager.GetLiveTeachers())
	assignOnTeachersCount := len(websocket.TeacherManager.GetAssignOnTeachers())
	content := map[string]interface{}{
		"onlineStudentsCount":   onlineStudentsCount,
		"onlineTeachersCount":   onlineTeachersCount,
		"liveTeachersCount":     liveTeachersCount,
		"assignOnTeachersCount": assignOnTeachersCount,
	}
	*resp = NewRpcResponse(0, "", content)
	return nil
}

func (watcher *RpcWatcher) GetOnlineTeacher(request *RpcRequest, resp *RpcResponse) error {
	teacherList := make([]*models.User, 0)
	for userId, _ := range websocket.UserManager.OnlineUserMap {
		user, err := models.ReadUser(userId)
		if err != nil {
			continue
		}

		if user.AccessRight == models.USER_ACCESSRIGHT_TEACHER {
			teacherList = append(teacherList, user)
		}
	}

	*resp = NewRpcResponse(0, "", teacherList)
	return nil
}

func (watcher *RpcWatcher) GetOnlineStudent(request *RpcRequest, resp *RpcResponse) error {
	studentList := make([]*models.User, 0)
	for userId, _ := range websocket.UserManager.OnlineUserMap {
		user, err := models.ReadUser(userId)
		if err != nil {
			continue
		}

		if user.AccessRight == models.USER_ACCESSRIGHT_STUDENT {
			studentList = append(studentList, user)
		}
	}

	*resp = NewRpcResponse(0, "", studentList)
	return nil
}

func (watcher *RpcWatcher) GetDispatchableTeacher(request *RpcRequest, resp *RpcResponse) error {
	dispatchableTeacherList := websocket.TeacherManager.GetLiveTeachers()

	teacherList := make([]*models.User, 0)
	for _, userId := range dispatchableTeacherList {
		user, err := models.ReadUser(userId)
		if err != nil {
			continue
		}

		teacherList = append(teacherList, user)
	}

	*resp = NewRpcResponse(0, "", teacherList)
	return nil
}

func (watcher *RpcWatcher) GetAssignableTeacher(request *RpcRequest, resp *RpcResponse) error {
	assignableTeacherList := websocket.TeacherManager.GetAssignOnTeachers()

	teacherList := make([]*models.User, 0)
	for _, userId := range assignableTeacherList {
		user, err := models.ReadUser(userId)
		if err != nil {
			continue
		}

		teacherList = append(teacherList, user)
	}

	*resp = NewRpcResponse(0, "", teacherList)
	return nil
}

func (watcher *RpcWatcher) GetUserStatus(request *RpcRequest, resp *RpcResponse) error {
	var err error

	userId, err := strconv.ParseInt(request.Args["userId"], 10, 64)
	if err != nil {
		*resp = NewRpcResponse(2, "无效的用户ID", response.NullObject)
		return err
	}

	status := websocket.UserManager.GetUserStatus(userId)

	*resp = NewRpcResponse(0, "", status)
	return nil
}
