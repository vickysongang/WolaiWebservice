package rpc

import (
	"strconv"

	courseController "WolaiWebservice/controllers/course"
	"WolaiWebservice/handlers/response"
	"WolaiWebservice/models"
	"WolaiWebservice/service/trade"
	tradeService "WolaiWebservice/service/trade"

	"WolaiWebservice/utils/leancloud/lcmessage"
	"WolaiWebservice/utils/wechat/wcmessage"
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

		err = trade.HandleCourseEarning(recordId, period, chapterId)
		if err != nil {
			*resp = NewRpcResponse(2, "交易失败", response.NullObject)
			return err
		}
	} else if course.Type == models.COURSE_TYPE_AUDITION {

		err = trade.HandleAuditionCourseEarning(recordId, period, chapterId)
		if err != nil {
			*resp = NewRpcResponse(2, "交易失败", response.NullObject)
			return err
		}
	}

	go wcmessage.SendChapterCompleteNotification(recordId, chapterId, course.Name, course.Type)

	*resp = NewRpcResponse(0, "", response.NullObject)
	return nil
}

func (watcher *RpcWatcher) HandleTradeRewardRegistration(request *RpcRequest, resp *RpcResponse) error {
	var err error

	userId, err := strconv.ParseInt(request.Args["userId"], 10, 64)
	if err != nil {
		*resp = NewRpcResponse(2, "无效的用户Id", response.NullObject)
		return err
	}
	tradeService.HandleTradeRewardGivenQaPkg(userId, tradeService.COMMENT_QA_PKG_GIVEN_INVITATION)
	go lcmessage.SendWelcomeMessageStudent(userId)
	*resp = NewRpcResponse(0, "", response.NullObject)
	return nil
}

func (watcher *RpcWatcher) HandleTradeQapkgGiven(request *RpcRequest, resp *RpcResponse) error {
	var err error

	userId, err := strconv.ParseInt(request.Args["userId"], 10, 64)
	if err != nil {
		*resp = NewRpcResponse(2, "无效的用户ID", response.NullObject)
		return err
	}

	qapkgId, err := strconv.ParseInt(request.Args["qapkgId"], 10, 64)
	if err != nil {
		*resp = NewRpcResponse(2, "无效的家教时间包id", response.NullObject)
		return err
	}

	comment := request.Args["comment"]

	err = trade.HandleTradeQapkgGiven(userId, qapkgId, comment)
	if err != nil {
		*resp = NewRpcResponse(2, "交易失败", response.NullObject)
		return err
	}

	*resp = NewRpcResponse(0, "", response.NullObject)
	return nil
}

func (watcher *RpcWatcher) HandleTradeCoursePurchase(request *RpcRequest, resp *RpcResponse) error {
	var err error

	recordId, err := strconv.ParseInt(request.Args["recordId"], 10, 64)
	if err != nil {
		*resp = NewRpcResponse(2, "无效的课程购买记录Id", response.NullObject)
		return err
	}
	comment := request.Args["comment"]
	err = tradeService.HandleCoursePurchaseTradeRecord(recordId, 0, comment)
	if err != nil {
		*resp = NewRpcResponse(2, "交易失败", response.NullObject)
		return err
	}
	*resp = NewRpcResponse(0, "", response.NullObject)
	return nil
}

func (watcher *RpcWatcher) HandleTradeCourseQuotaPurchase(request *RpcRequest, resp *RpcResponse) error {
	var err error

	recordId, err := strconv.ParseInt(request.Args["recordId"], 10, 64)
	amount, err := strconv.ParseInt(request.Args["amount"], 10, 64)
	if err != nil {
		*resp = NewRpcResponse(2, "无效的通用课时购买记录Id", response.NullObject)
		return err
	}
	comment := request.Args["comment"]
	err = tradeService.HandleCourseQuotaPurchaseTradeRecord(recordId, amount, 0, comment)
	if err != nil {
		*resp = NewRpcResponse(2, "交易失败", response.NullObject)
		return err
	}
	*resp = NewRpcResponse(0, "", response.NullObject)
	return nil
}

func (watcher *RpcWatcher) HandleTradeCourseQuotaRefund(request *RpcRequest, resp *RpcResponse) error {
	var err error

	recordId, err := strconv.ParseInt(request.Args["recordId"], 10, 64)
	amount, err := strconv.ParseInt(request.Args["amount"], 10, 64)
	if err != nil {
		*resp = NewRpcResponse(2, "无效的通用课时退款记录Id", response.NullObject)
		return err
	}
	comment := request.Args["comment"]
	err = tradeService.HandleCourseQuotaRefundTradeRecord(recordId, amount, 0, comment)
	if err != nil {
		*resp = NewRpcResponse(2, "交易失败", response.NullObject)
		return err
	}
	*resp = NewRpcResponse(0, "", response.NullObject)
	return nil
}

func (watcher *RpcWatcher) HandleTradeCourseRefundToWallet(request *RpcRequest, resp *RpcResponse) error {
	var err error

	recordId, err := strconv.ParseInt(request.Args["recordId"], 10, 64)
	amount, err := strconv.ParseInt(request.Args["amount"], 10, 64)
	if err != nil {
		*resp = NewRpcResponse(2, "无效的课程购买记录Id", response.NullObject)
		return err
	}
	comment := request.Args["comment"]
	err = tradeService.HandleCourseRefundToWalletTradeRecord(recordId, amount, 0, comment)
	if err != nil {
		*resp = NewRpcResponse(2, "交易失败", response.NullObject)
		return err
	}
	*resp = NewRpcResponse(0, "", response.NullObject)
	return nil
}

func (watcher *RpcWatcher) HandleTradeCourseRefundToQuota(request *RpcRequest, resp *RpcResponse) error {
	var err error

	recordId, err := strconv.ParseInt(request.Args["recordId"], 10, 64)
	amount, err := strconv.ParseInt(request.Args["amount"], 10, 64)
	if err != nil {
		*resp = NewRpcResponse(2, "无效的课程退款到通用课时记录Id", response.NullObject)
		return err
	}
	comment := request.Args["comment"]
	err = tradeService.HandleCourseRefundToQuotaTradeRecord(recordId, amount, 0, comment)
	if err != nil {
		*resp = NewRpcResponse(2, "交易失败", response.NullObject)
		return err
	}
	*resp = NewRpcResponse(0, "", response.NullObject)
	return nil
}

func (watcher *RpcWatcher) HandleDeluxeCoursePayByQuota(request *RpcRequest, resp *RpcResponse) error {
	var err error

	userId, err := strconv.ParseInt(request.Args["userId"], 10, 64)
	courseId, err := strconv.ParseInt(request.Args["courseId"], 10, 64)
	if err != nil {
		*resp = NewRpcResponse(2, "无效的课程Id", response.NullObject)
		return err
	}
	comment := request.Args["comment"]
	_, err = courseController.HandleDeluxeCoursePayByQuota(userId, courseId, comment)
	if err != nil {
		*resp = NewRpcResponse(2, "交易失败", response.NullObject)
		return err
	}
	*resp = NewRpcResponse(0, "", response.NullObject)
	return nil
}
