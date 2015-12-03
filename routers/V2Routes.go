// V2Routes
package routers

import (
	"WolaiWebservice/handlers"
	"WolaiWebservice/websocket"
)

var V2Routes = Routes{
	/********************************Websocket******************************/
	Route{
		"V2WebSocketHandler",
		"GET",
		"/v2/ws",
		websocket.V1WebSocketHandler,
	},

	/********************************1.注册&登陆******************************/
	// 1.1 注册时绑定QQ
	Route{
		"V2OauthRegisterPOST",
		"POST",
		"/v2/oauth/qq/register",
		handlers.V2OauthRegister,
	},
	Route{
		"V2OauthRegisterGET",
		"GET",
		"/v2/oauth/qq/register",
		handlers.V2OauthRegister,
	},

	//1.2 检查手机号码是否已经绑定过QQ
	Route{
		"V2CheckPhoneGET",
		"GET",
		"/v2/oauth/qq/checkphone",
		handlers.V2CheckPhoneBindWithQQ,
	},
	Route{
		"V2CheckPhonePost",
		"POST",
		"/v2/oauth/qq/checkphone",
		handlers.V2CheckPhoneBindWithQQ,
	},

	// 1.3 用户通过手机号码登陆
	Route{
		"V2LoginPOST",
		"POST",
		"/v2/login",
		handlers.V2Login,
	},
	Route{
		"V2LoginGET",
		"GET",
		"/v2/login",
		handlers.V2Login,
	},
	Route{
		"V2LoginGETURL",
		"GET",
		"/v2/login/{phone}",
		handlers.V2LoginGETURL,
	},

	// 1.4 用户通过QQ绑定登陆
	Route{
		"V2OauthLoginPOST",
		"POST",
		"/v2/oauth/qq/login",
		handlers.V2OauthLogin,
	},
	Route{
		"V2OauthLoginGET",
		"GET",
		"/v2/oauth/qq/login",
		handlers.V2OauthLogin,
	},

	// 1.5 用户登陆时获取用户的定位信息
	Route{
		"V2InsertUserLoginInfoPOST",
		"POST",
		"/v2/user/logininfo/insert",
		handlers.V2InsertUserLoginInfo,
	},
	Route{
		"V2InsertUserLoginInfoGET",
		"GET",
		"/v2/user/logininfo/insert",
		handlers.V2InsertUserLoginInfo,
	},

	/********************************2.侧边栏相关******************************/
	// 2.1 修改用户的个人资料
	Route{
		"V2UpdateProfilePOST",
		"POST",
		"/v2/update_profile",
		handlers.V2UpdateProfile,
	},
	Route{
		"V2UpdateProfileGET",
		"GET",
		"/v2/update_profile",
		handlers.V2UpdateProfile,
	},
	Route{
		"V2UpdateProfileGETURL",
		"GET",
		"/v2/update_profile/{userId}/{nickname}/{avatar}/{gender}",
		handlers.V2UpdateProfileGETURL,
	},

	// 2.2 获取用户的个人资料
	Route{
		"V2UserInfoPOST",
		"POST",
		"/v2/user/info",
		handlers.V2UserInfo,
	},
	Route{
		"V2UserInfoGET",
		"GET",
		"/v2/user/info",
		handlers.V2UserInfo,
	},

	// 2.3 我的课程
	Route{
		"V2MyOrdersGET",
		"GET",
		"/v2/user/myorders",
		handlers.V2OrderInSession,
	},
	Route{
		"V2MyOrdersPost",
		"POST",
		"/v2/user/myorders",
		handlers.V2OrderInSession,
	},

	// 2.4 我的钱包
	Route{
		"V2UserMyWalletPOST",
		"POST",
		"/v2/user/mywallet",
		handlers.V2UserMyWallet,
	},
	Route{
		"V2UserMyWalletGET",
		"GET",
		"/v2/user/mywallet",
		handlers.V2UserMyWallet,
	},

	// 2.5 钱包明细
	Route{
		"V2MyTradeRecordGET",
		"GET",
		"/v2/trade/traderecord",
		handlers.V2TradeRecord,
	},
	Route{
		"V2MyTradeRecordPOST",
		"POST",
		"/v2/trade/traderecord",
		handlers.V2TradeRecord,
	},

	// 2.6 我的动态
	Route{
		"V2UserMyFeedPOST",
		"POST",
		"/v2/user/myfeed",
		handlers.V2UserMyFeed,
	},
	Route{
		"V2UserMyFeedGET",
		"GET",
		"/v2/user/myfeed",
		handlers.V2UserMyFeed,
	},

	// 2.7 我喜欢的
	Route{
		"V2UserMyLikePOST",
		"POST",
		"/v2/user/mylike",
		handlers.V2UserMyLike,
	},
	Route{
		"V2UserMyLikeGET",
		"GET",
		"/v2/user/mylike",
		handlers.V2UserMyLike,
	},

	// 2.8 绑定邀请码
	Route{
		"V2BindInvitationCodePOST",
		"POST",
		"/v2/invitation/bind",
		handlers.V2BindUserWithInvitationCode,
	},
	Route{
		"V2BindInvitationCodeGET",
		"GET",
		"/v2/invitation/bind",
		handlers.V2BindUserWithInvitationCode,
	},

	// 2.9 检查用户是否已经绑定过邀请码
	Route{
		"V2CheckUserBindWithInvitationCodePOST",
		"POST",
		"/v2/invitation/check",
		handlers.V2CheckUserHasBindWithInvitationCode,
	},
	Route{
		"V2CheckUserBindWithInvitationCodeGET",
		"GET",
		"/v2/invitation/check",
		handlers.V2CheckUserHasBindWithInvitationCode,
	},

	/********************************3.导师信息 **********************************/
	// 3.1 获取推荐的导师
	Route{
		"V2TeacherRecommendationPOST",
		"POST",
		"/v2/teacher/recommendation",
		handlers.V2TeacherRecommendation,
	},
	Route{
		"V2TeacherRecommendationGET",
		"GET",
		"/v2/teacher/recommendation",
		handlers.V2TeacherRecommendation,
	},

	// 3.2 查看导师的详细资料
	Route{
		"V2TeacherProfilePOST",
		"POST",
		"/v2/teacher/profile",
		handlers.V2TeacherProfile,
	},
	Route{
		"V2TeacherProfileGET",
		"GET",
		"/v2/teacher/profile",
		handlers.V2TeacherProfile,
	},

	// 3.3 获取老师的课时费
	Route{
		"V2TeacherExpectPost",
		"POST",
		"/v2/teacher/expect",
		handlers.V2TeacherExpect,
	},
	Route{
		"V2TeacherExpectGET",
		"GET",
		"/v2/teacher/expect",
		handlers.V2TeacherExpect,
	},

	// 3.4 获取我来团队、我来助教、导师列表
	Route{
		"V2SupportAndTeacherListPOST",
		"POST",
		"/v2/support/teacher/list",
		handlers.V2SupportAndTeacherList,
	},
	Route{
		"V2TSupportAndTeacherListGET",
		"GET",
		"/v2/support/teacher/list",
		handlers.V2SupportAndTeacherList,
	},

	/********************************4.动态&干货*******************************/
	// 4.1 获取所有的动态
	Route{
		"V2AtriumPOST",
		"POST",
		"/v2/atrium",
		handlers.V2Atrium,
	},
	Route{
		"V2AtriumGET",
		"GET",
		"/v2/atrium",
		handlers.V2Atrium,
	},

	// 4.2 发布动态
	Route{
		"V2FeedPostPOST",
		"POST",
		"/v2/feed/post",
		handlers.V2FeedPost,
	},
	Route{
		"V2FeedPostGET",
		"GET",
		"/v2/feed/post",
		handlers.V2FeedPost,
	},

	// 4.3 获取动态或干货详情
	Route{
		"V2FeedDetailPOST",
		"POST",
		"/v2/feed/detail",
		handlers.V2FeedDetail,
	},
	Route{
		"V2FeedDetailGET",
		"GET",
		"/v2/feed/detail",
		handlers.V2FeedDetail,
	},

	// 4.4 给动态标注喜欢
	Route{
		"V2FeedLikePOST",
		"POST",
		"/v2/feed/like",
		handlers.V2FeedLike,
	},
	Route{
		"V2FeedLikeGET",
		"GET",
		"/v2/feed/like",
		handlers.V2FeedLike,
	},

	// 4.5 评论动态或干货
	Route{
		"V2FeedCommentPOST",
		"POST",
		"/v2/feed/comment",
		handlers.V2FeedComment,
	},
	Route{
		"V2FeedCommentGET",
		"GET",
		"/v2/feed/comment",
		handlers.V2FeedComment,
	},

	// 4.6 获取置顶的干货
	Route{
		"V2TopFeedPOST",
		"POST",
		"/v2/feed/top/get",
		handlers.V2GetTopFeed,
	},
	Route{
		"V2TopFeedGET",
		"GET",
		"/v2/feed/top/get",
		handlers.V2GetTopFeed,
	},

	/********************************5.对话********************************/
	// 5.1 获取对话的ID
	Route{
		"V2GetConversationIDPOST",
		"POST",
		"/v2/conversation/get",
		handlers.V2GetConversationID,
	},
	Route{
		"V2GetConversationIDGET",
		"GET",
		"/v2/conversation/get",
		handlers.V2GetConversationID,
	},

	/********************************6.订单模块********************************/
	//6.1 获取年级列表
	Route{
		"V2GradeListPOST",
		"POST",
		"/v2/grade/list",
		handlers.V2GradeList,
	},
	Route{
		"V2GradeListGET",
		"GET",
		"/v2/grade/list",
		handlers.V2GradeList,
	},

	//6.2 获取科目列表
	Route{
		"V2SubjectListPOST",
		"POST",
		"/v2/subject/list",
		handlers.V2SubjectList,
	},
	Route{
		"V2SubjectListGET",
		"GET",
		"/v2/subject/list",
		handlers.V2SubjectList,
	},

	// 6.3 创建订单
	Route{
		"V2OrderCreatePOST",
		"POST",
		"/v2/order/create",
		handlers.V2OrderCreate,
	},
	Route{
		"V2OrderCreateGET",
		"GET",
		"/v2/order/create",
		handlers.V2OrderCreate,
	},

	// 6.4 确认点对点订单请求
	Route{
		"V2OrderPersonalConfirmPOST",
		"POST",
		"/v2/order/personal/confirm",
		handlers.V2OrderPersonalConfirm,
	},
	Route{
		"V2OrderPersonalConfirmGET",
		"GET",
		"/v2/order/personal/confirm",
		handlers.V2OrderPersonalConfirm,
	},

	/********************************7.搜索模块********************************/
	// 7.1 首页搜索导师
	Route{
		"V2SearchTeachersPOST",
		"POST",
		"/v2/search/teacher",
		handlers.V2SearchTeachers,
	},
	Route{
		"V2SearchTeachersGET",
		"GET",
		"/v2/search/teacher",
		handlers.V2SearchTeachers,
	},

	// 7.2 搜索导师或学生
	Route{
		"V2SearchUsersPOST",
		"POST",
		"/v2/search/user",
		handlers.V2SearchUsers,
	},
	Route{
		"V2SearchUsersGET",
		"GET",
		"/v2/search/user",
		handlers.V2SearchUsers,
	},

	/********************************8.课程评价********************************/
	// 8.1 获取系统评价标签列表
	Route{
		"V2GetEvaluationPOST",
		"POST",
		"/v2/evaluation/query",
		handlers.V2GetEvaluation,
	},
	Route{
		"V2GetEvaluationGET",
		"GET",
		"/v2/evaluation/query",
		handlers.V2GetEvaluation,
	},

	// 8.2 评价对方
	Route{
		"V2EvaluatePOST",
		"POST",
		"/v2/evaluation/insert",
		handlers.V2Evaluate,
	},
	Route{
		"V2EvaluateGET",
		"GET",
		"/v2/evaluation/insert",
		handlers.V2Evaluate,
	},

	// 8.3 获取评价后的标签
	Route{
		"V2GetEvaluationLabelPOST",
		"POST",
		"/v2/evaluation/label",
		handlers.V2GetEvaluationLabels,
	},
	Route{
		"V2GetEvaluationLabelGET",
		"GET",
		"/v2/evaluation/label",
		handlers.V2GetEvaluationLabels,
	},

	/******************************9.课程包月模块******************************/
	// 9.1 获取包月课程列表
	Route{
		"V2GetCoursesPOST",
		"POST",
		"/v2/course/list",
		handlers.V2GetCourses,
	},
	Route{
		"V2GetCoursesGET",
		"GET",
		"/v2/course/list",
		handlers.V2GetCourses,
	},

	// 9.2 学生申请加入包月课程
	Route{
		"V2JoinCoursesPOST",
		"POST",
		"/v2/course/join",
		handlers.V2JoinCourse,
	},
	Route{
		"V2JoinCoursesGET",
		"GET",
		"/v2/course/join",
		handlers.V2JoinCourse,
	},

	// 9.3 学生申请续费包月课程
	Route{
		"V2RenewCoursePOST",
		"POST",
		"/v2/course/renew",
		handlers.V2RenewCourse,
	},
	Route{
		"V2RenewCoursesGET",
		"GET",
		"/v2/course/renew",
		handlers.V2RenewCourse,
	},

	/********************************10.用户投诉*******************************/
	// 10.1 学生发起对导师的投诉
	Route{
		"V2ComplainPOST",
		"POST",
		"/v2/complaint/complain",
		handlers.V2Complain,
	},
	Route{
		"V2ComplainGET",
		"GET",
		"/v2/complaint/complain",
		handlers.V2Complain,
	},

	// 10.2 检查是否存在未处理的投诉
	Route{
		"V2CheckComplaintExsitsPOST",
		"POST",
		"/v2/complaint/check",
		handlers.V2CheckComplaintExsits,
	},
	Route{
		"V2CheckComplaintExsitsGET",
		"GET",
		"/v2/complaint/check",
		handlers.V2CheckComplaintExsits,
	},

	/******************************11.支付模块********************************/
	// 11.1 支付
	Route{
		"V2PayByPingppPOST",
		"POST",
		"/v2/pingpp/pay",
		handlers.V2PayByPingpp,
	},
	Route{
		"V2PayByPingppGET",
		"GET",
		"/v2/pingpp/pay",
		handlers.V2PayByPingpp,
	},

	// 11.2 退款
	Route{
		"V2RefundByPingppPOST",
		"POST",
		"/v2/pingpp/refund",
		handlers.V2RefundByPingpp,
	},
	Route{
		"V2RefundByPingppGET",
		"GET",
		"/v2/pingpp/refund",
		handlers.V2RefundByPingpp,
	},

	// 11.3 查询单笔支付记录
	Route{
		"V2QueryPaymentByPingppPOST",
		"POST",
		"/v2/pingpp/pay/query",
		handlers.V2QueryPaymentByPingpp,
	},
	Route{
		"V2QueryPaymentByPingppGET",
		"GET",
		"/v2/pingpp/pay/query",
		handlers.V2QueryPaymentByPingpp,
	},

	// 11.4 查询支付记录列表
	Route{
		"V2QueryPaymentListByPingppPOST",
		"POST",
		"/v2/pingpp/pay/list",
		handlers.V2QueryPaymentListByPingpp,
	},
	Route{
		"V2QueryPaymentListByPingppGET",
		"GET",
		"/v2/pingpp/pay/list",
		handlers.V2QueryPaymentListByPingpp,
	},

	// 11.5 查询单笔退款记录
	Route{
		"V2QueryRefundByPingppPOST",
		"POST",
		"/v2/pingpp/refund/query",
		handlers.V2QueryRefundByPingpp,
	},
	Route{
		"V2QueryRefundByPingppGET",
		"GET",
		"/v2/pingpp/refund/query",
		handlers.V2QueryRefundByPingpp,
	},

	// 11.6 查询退款记录列表
	Route{
		"V2QueryRefundListByPingppPOST",
		"POST",
		"/v2/pingpp/refund/list",
		handlers.V2QueryRefundListByPingpp,
	},
	Route{
		"V2QueryRefundListByPingppGET",
		"GET",
		"/v2/pingpp/refund/list",
		handlers.V2QueryRefundListByPingpp,
	},

	// 11.7 异步通知
	Route{
		"V2WebhookByPingppPOST",
		"POST",
		"/v2/pingpp/webhook",
		handlers.V2WebhookByPingpp,
	},
	Route{
		"V2WebhookByPingppGET",
		"GET",
		"/v2/pingpp/webhook",
		handlers.V2WebhookByPingpp,
	},

	// 11.8 查看支付结果
	Route{
		"V2GetPingResultPOST",
		"POST",
		"/v2/pingpp/result",
		handlers.V2GetPingppResult,
	},
	Route{
		"V2GetPingResultGET",
		"GET",
		"/v2/pingpp/result",
		handlers.V2GetPingppResult,
	},

	/******************************12.短信验证码*******************************/
	// 12.1 短信发送异步通知
	Route{
		"V2WSMSHookPOST",
		"POST",
		"/v2/sendcloud/smshook",
		handlers.V2SmsHook,
	},

	// 12.2 发送验证码
	Route{
		"V2SCSendMessagePOST",
		"POST",
		"/v2/sendcloud/sendmessage",
		handlers.V2SendMessage,
	},

	/******************************13.其他********************************/
	// 13.1 获取首页Banner
	Route{
		"V2BannerPOST",
		"POST",
		"/v2/banner",
		handlers.V2Banner,
	},
	Route{
		"V2BannerGET",
		"GET",
		"/v2/banner",
		handlers.V2Banner,
	},

	// 13.2 推送广告
	Route{
		"V2SendAdvMessagePOST",
		"POST",
		"/v2/send/adv",
		handlers.V2SendAdvMessage,
	},
	Route{
		"V2SendAdvMessageGET",
		"GET",
		"/v2/send/adv",
		handlers.V2SendAdvMessage,
	},

	// 13.3 获取帮助秘籍
	Route{
		"V2GetHelpItemsPOST",
		"POST",
		"/v2/help/get",
		handlers.V2GetHelpItems,
	},
	Route{
		"V2GetHelpItemsGET",
		"GET",
		"/v2/help/get",
		handlers.V2GetHelpItems,
	},

	// 13.4 求助客服
	Route{
		"V2PostSeekHelpPOST",
		"POST",
		"/v2/seekhelp/post",
		handlers.V2SetSeekHelp,
	},
	Route{
		"V2PostSeekHelpGET",
		"GET",
		"/v2/seekhelp/post",
		handlers.V2SetSeekHelp,
	},
	// 13.5 将活动信息发送给用户
	Route{
		"V2GetActivitiesPOST",
		"POST",
		"/v2/activity/notification",
		handlers.V2ActivityNotification,
	},
	Route{
		"V2GetActivitiesGET",
		"GET",
		"/v2/activity/notification",
		handlers.V2ActivityNotification,
	},
}
