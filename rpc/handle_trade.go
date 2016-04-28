package rpc

import (
	"strconv"

	"WolaiWebservice/handlers/response"
	"WolaiWebservice/models"
	"WolaiWebservice/service/trade"
	"WolaiWebservice/utils/leancloud/lcmessage"
)

func (watcher *RpcWatcher) HandleTradeVoucher(request *RpcRequest, resp *RpcResponse) error {
	var err error

	userId, err := strconv.ParseInt(request.Args["userId"], 10, 64)
	if err != nil {
		*resp = NewRpcResponse(2, "无效的用户ID", response.NullObject)
		return err
	}

	amount, err := strconv.ParseInt(request.Args["amount"], 10, 64)
	if err != nil {
		*resp = NewRpcResponse(2, "无效的金额", response.NullObject)
		return err
	}

	comment := request.Args["comment"]

	err = trade.HandleTradeVoucher(userId, amount, comment)
	if err != nil {
		*resp = NewRpcResponse(2, "交易失败", response.NullObject)
		return err
	}

	*resp = NewRpcResponse(0, "", response.NullObject)
	return nil
}

func (watcher *RpcWatcher) HandleTradePromotion(request *RpcRequest, resp *RpcResponse) error {
	var err error

	userId, err := strconv.ParseInt(request.Args["userId"], 10, 64)
	if err != nil {
		*resp = NewRpcResponse(2, "无效的用户ID", response.NullObject)
		return err
	}

	amount, err := strconv.ParseInt(request.Args["amount"], 10, 64)
	if err != nil {
		*resp = NewRpcResponse(2, "无效的金额", response.NullObject)
		return err
	}

	comment := request.Args["comment"]

	err = trade.HandleTradePromotion(userId, amount, comment)
	if err != nil {
		*resp = NewRpcResponse(2, "交易失败", response.NullObject)
		return err
	}

	*resp = NewRpcResponse(0, "", response.NullObject)
	return nil
}

func (watcher *RpcWatcher) HandleTradeDeduction(request *RpcRequest, resp *RpcResponse) error {
	var err error

	userId, err := strconv.ParseInt(request.Args["userId"], 10, 64)
	if err != nil {
		*resp = NewRpcResponse(2, "无效的用户ID", response.NullObject)
		return err
	}

	amount, err := strconv.ParseInt(request.Args["amount"], 10, 64)
	if err != nil {
		*resp = NewRpcResponse(2, "无效的金额", response.NullObject)
		return err
	}

	comment := request.Args["comment"]

	err = trade.HandleTradeDeduction(userId, amount, comment)
	if err != nil {
		*resp = NewRpcResponse(2, "交易失败", response.NullObject)
		return err
	}

	*resp = NewRpcResponse(0, "", response.NullObject)
	return nil
}

func (watcher *RpcWatcher) HandleTradeWithdraw(request *RpcRequest, resp *RpcResponse) error {
	var err error

	userId, err := strconv.ParseInt(request.Args["userId"], 10, 64)
	if err != nil {
		*resp = NewRpcResponse(2, "无效的用户ID", response.NullObject)
		return err
	}

	amount, err := strconv.ParseInt(request.Args["amount"], 10, 64)
	if err != nil {
		*resp = NewRpcResponse(2, "无效的金额", response.NullObject)
		return err
	}

	err = trade.HandleTradeWithdraw(userId, amount)
	if err != nil {
		*resp = NewRpcResponse(2, "交易失败", response.NullObject)
		return err
	}

	*resp = NewRpcResponse(0, "", response.NullObject)
	return nil
}

func (watcher *RpcWatcher) HandleCourseEarning(request *RpcRequest, resp *RpcResponse) error {
	var err error

	recordId, err := strconv.ParseInt(request.Args["recordId"], 10, 64)
	if err != nil {
		*resp = NewRpcResponse(2, "无效的购买记录ID", response.NullObject)
		return err
	}
	chapterId, err := strconv.ParseInt(request.Args["chapterId"], 10, 64)
	if err != nil {
		*resp = NewRpcResponse(2, "无效的课时ID", response.NullObject)
		return err
	}
	period, err := strconv.ParseInt(request.Args["period"], 10, 64)
	if err != nil {
		*resp = NewRpcResponse(2, "无效的课时号", response.NullObject)
		return err
	}

	chapter, _ := models.ReadCourseCustomChapter(chapterId)
	course, err := models.ReadCourse(chapter.CourseId)
	if err != nil {
		*resp = NewRpcResponse(2, "无效的课程信息", response.NullObject)
		return err
	}
	if course.Type == models.COURSE_TYPE_DELUXE {
		go lcmessage.SendCourseChapterCompleteMsg(recordId, chapterId)

		err = trade.HandleCourseEarning(recordId, period, chapterId)
		if err != nil {
			*resp = NewRpcResponse(2, "交易失败", response.NullObject)
			return err
		}
	} else if course.Type == models.COURSE_TYPE_AUDITION {
		go lcmessage.SendAuditionCourseChapterCompleteMsg(recordId, chapterId)

		err = trade.HandleAuditionCourseEarning(recordId, period, chapterId)
		if err != nil {
			*resp = NewRpcResponse(2, "交易失败", response.NullObject)
			return err
		}
	}

	*resp = NewRpcResponse(0, "", response.NullObject)
	return nil
}
