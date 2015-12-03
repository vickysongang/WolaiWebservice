package routers

import (
	"WolaiWebservice/handlers"
	"WolaiWebservice/routers/route"
	"WolaiWebservice/websocket"
)

var routesV1 = route.Routes{
	// Websocket
	route.Route{
		"V1WebSocketHandler",
		"GET",
		"/v1/ws",
		websocket.V1WebSocketHandler,
	},

	// 1.1 Login
	route.Route{
		"V1LoginPOST",
		"POST",
		"/v1/login",
		handlers.V1Login,
	},
	route.Route{
		"V1LoginGET",
		"GET",
		"/v1/login",
		handlers.V1Login,
	},
	route.Route{
		"V1LoginGETURL",
		"GET",
		"/v1/login/{phone}",
		handlers.V1LoginGETURL,
	},

	// 1.2 Update profile
	route.Route{
		"V1UpdateProfilePOST",
		"POST",
		"/v1/update_profile",
		handlers.V1UpdateProfile,
	},
	route.Route{
		"V1UpdateProfileGET",
		"GET",
		"/v1/update_profile",
		handlers.V1UpdateProfile,
	},
	route.Route{
		"V1UpdateProfileGETURL",
		"GET",
		"/v1/update_profile/{userId}/{nickname}/{avatar}/{gender}",
		handlers.V1UpdateProfileGETURL,
	},

	// 1.3 Oauth Login
	route.Route{
		"V1OauthLoginPOST",
		"POST",
		"/v1/oauth/qq/login",
		handlers.V1OauthLogin,
	},
	route.Route{
		"V1OauthLoginGET",
		"GET",
		"/v1/oauth/qq/login",
		handlers.V1OauthLogin,
	},

	// 1.4 Oauth Register
	route.Route{
		"V1OauthRegisterPOST",
		"POST",
		"/v1/oauth/qq/register",
		handlers.V1OauthRegister,
	},
	route.Route{
		"V1OauthRegisterGET",
		"GET",
		"/v1/oauth/qq/register",
		handlers.V1OauthRegister,
	},

	// 1.6 Teacher Recommendation
	route.Route{
		"V1TeacherRecommendationPOST",
		"POST",
		"/v1/teacher/recommendation",
		handlers.V1TeacherRecommendation,
	},
	route.Route{
		"V1TeacherRecommendationGET",
		"GET",
		"/v1/teacher/recommendation",
		handlers.V1TeacherRecommendation,
	},

	// 1.7 Teacher Profile
	route.Route{
		"V1TeacherProfilePOST",
		"POST",
		"/v1/teacher/profile",
		handlers.V1TeacherProfile,
	},
	route.Route{
		"V1TeacherProfileGET",
		"GET",
		"/v1/teacher/profile",
		handlers.V1TeacherProfile,
	},
	//1.8 Teacher post
	route.Route{
		"V1InsertTeacherPost",
		"POST",
		"/v1/teacher/insert",
		handlers.V1TeacherPost,
	},
	//1.9 Check Phone
	route.Route{
		"V1CheckPhoneGET",
		"GET",
		"/v1/oauth/qq/checkphone",
		handlers.V1CheckPhoneBindWithQQ,
	},
	route.Route{
		"V1CheckPhonePost",
		"POST",
		"/v1/oauth/qq/checkphone",
		handlers.V1CheckPhoneBindWithQQ,
	},
	// 1.10 support and teacher list
	route.Route{
		"V1SupportAndTeacherListPOST",
		"POST",
		"/v1/support/teacher/list",
		handlers.V1SupportAndTeacherList,
	},
	route.Route{
		"V1TSupportAndTeacherListGET",
		"GET",
		"/v1/support/teacher/list",
		handlers.V1SupportAndTeacherList,
	},
	// 1.11 user login info insert
	route.Route{
		"V1InsertUserLoginInfoPOST",
		"POST",
		"/v1/user/logininfo/insert",
		handlers.V1InsertUserLoginInfo,
	},
	route.Route{
		"V1InsertUserLoginInfoGET",
		"GET",
		"/v1/user/logininfo/insert",
		handlers.V1InsertUserLoginInfo,
	},

	// 2.1 Atrium
	route.Route{
		"V1AtriumPOST",
		"POST",
		"/v1/atrium",
		handlers.V1Atrium,
	},
	route.Route{
		"V1AtriumGET",
		"GET",
		"/v1/atrium",
		handlers.V1Atrium,
	},

	// 2.2 Feed Post
	route.Route{
		"V1FeedPostPOST",
		"POST",
		"/v1/feed/post",
		handlers.V1FeedPost,
	},
	route.Route{
		"V1FeedPostGET",
		"GET",
		"/v1/feed/post",
		handlers.V1FeedPost,
	},

	// 2.3 Feed Detail
	route.Route{
		"V1FeedDetailPOST",
		"POST",
		"/v1/feed/detail",
		handlers.V1FeedDetail,
	},
	route.Route{
		"V1FeedDetailGET",
		"GET",
		"/v1/feed/detail",
		handlers.V1FeedDetail,
	},

	// 2.4 Feed Like
	route.Route{
		"V1FeedLikePOST",
		"POST",
		"/v1/feed/like",
		handlers.V1FeedLike,
	},
	route.Route{
		"V1FeedLikeGET",
		"GET",
		"/v1/feed/like",
		handlers.V1FeedLike,
	},

	// 2.5 Feed Favorite
	/*
		route.Route{
			"V1FeedFavPOST",
			"POST",
			"/v1/feed/favorite",
			V1FeedFav,
		},
		route.Route{
			"V1FeedFavGET",
			"GET",
			"/v1/feed/favorite",
			V1FeedFav,
		},
	*/

	// 2.6 Feed Comment
	route.Route{
		"V1FeedCommentPOST",
		"POST",
		"/v1/feed/comment",
		handlers.V1FeedComment,
	},
	route.Route{
		"V1FeedCommentGET",
		"GET",
		"/v1/feed/comment",
		handlers.V1FeedComment,
	},

	// 2.7 Feed Comment Like
	/*
		route.Route{
			"V1FeedCommentLikePOST",
			"POST",
			"/v1/feed/comment/like",
			V1FeedCommentLike,
		},
		route.Route{
			"V1FeedCommentLikeGET",
			"GET",
			"/v1/feed/comment/like",
			V1FeedCommentLike,
		},
	*/
	// 2.8 Feed mark
	route.Route{
		"V1MarkFeedPOST",
		"POST",
		"/v1/feed/mark",
		handlers.V1FeedMark,
	},
	route.Route{
		"V1MarkFeedGET",
		"GET",
		"/v1/feed/mark",
		handlers.V1FeedMark,
	},

	// 2.9 Top Feed
	route.Route{
		"V1TopFeedPOST",
		"POST",
		"/v1/feed/top/get",
		handlers.V1GetTopFeed,
	},
	route.Route{
		"V1TopFeedGET",
		"GET",
		"/v1/feed/top/get",
		handlers.V1GetTopFeed,
	},
	// 2.10 Delete Feed
	route.Route{
		"V1DeleteFeedPOST",
		"POST",
		"/v1/feed/delete",
		handlers.V1FeedDelete,
	},
	route.Route{
		"V1DeleteFeedGET",
		"GET",
		"/v1/feed/delete",
		handlers.V1FeedDelete,
	},
	// 2.11 Recover Feed
	route.Route{
		"V1RecoverFeedPOST",
		"POST",
		"/v1/feed/recover",
		handlers.V1FeedRecover,
	},
	route.Route{
		"V1RecoverFeedGET",
		"GET",
		"/v1/feed/recover",
		handlers.V1FeedRecover,
	},
	// 2.12 Top Feed
	route.Route{
		"V1TopFeedPOST",
		"POST",
		"/v1/feed/top",
		handlers.V1FeedTop,
	},
	route.Route{
		"V1TopFeedGET",
		"GET",
		"/v1/feed/top",
		handlers.V1FeedTop,
	},

	// 3.1 User Info
	route.Route{
		"V1UserInfoPOST",
		"POST",
		"/v1/user/info",
		handlers.V1UserInfo,
	},
	route.Route{
		"V1UserInfoGET",
		"GET",
		"/v1/user/info",
		handlers.V1UserInfo,
	},

	// 3.2 User Wallet
	route.Route{
		"V1UserMyWalletPOST",
		"POST",
		"/v1/user/mywallet",
		handlers.V1UserMyWallet,
	},
	route.Route{
		"V1UserMyWalletGET",
		"GET",
		"/v1/user/mywallet",
		handlers.V1UserMyWallet,
	},

	// 3.3 User MyFeed
	route.Route{
		"V1UserMyFeedPOST",
		"POST",
		"/v1/user/myfeed",
		handlers.V1UserMyFeed,
	},
	route.Route{
		"V1UserMyFeedGET",
		"GET",
		"/v1/user/myfeed",
		handlers.V1UserMyFeed,
	},

	// 3.4 User MyFollowing
	route.Route{
		"V1UserMyFollowPOST",
		"POST",
		"/v1/user/myfollow",
		handlers.V1UserMyFollowing,
	},
	route.Route{
		"V1UserMyFollowGET",
		"GET",
		"/v1/user/myfollow",
		handlers.V1UserMyFollowing,
	},

	// 3.5 User MyLike
	route.Route{
		"V1UserMyLikePOST",
		"POST",
		"/v1/user/mylike",
		handlers.V1UserMyLike,
	},
	route.Route{
		"V1UserMyLikeGET",
		"GET",
		"/v1/user/mylike",
		handlers.V1UserMyLike,
	},

	// 3.6 User Follow
	route.Route{
		"V1UserFollowPOST",
		"POST",
		"/v1/user/follow",
		handlers.V1UserFollow,
	},
	route.Route{
		"V1UserFollowGET",
		"GET",
		"/v1/user/follow",
		handlers.V1UserFollow,
	},

	// 3.7 User Unfollow
	route.Route{
		"V1UserUnfollowPOST",
		"POST",
		"/v1/user/unfollow",
		handlers.V1UserUnfollow,
	},
	route.Route{
		"V1UserUnfollowGET",
		"GET",
		"/v1/user/unfollow",
		handlers.V1UserUnfollow,
	},

	// 3.8 User Order
	route.Route{
		"V1MyOrdersGET",
		"GET",
		"/v1/user/myorders",
		handlers.V1OrderInSession,
	},
	route.Route{
		"V1MyOrdersPost",
		"POST",
		"/v1/user/myorders",
		handlers.V1OrderInSession,
	},

	// 4.1 Get Conversation ID
	route.Route{
		"V1GetConversationIDPOST",
		"POST",
		"/v1/conversation/get",
		handlers.V1GetConversationID,
	},
	route.Route{
		"V1GetConversationIDGET",
		"GET",
		"/v1/conversation/get",
		handlers.V1GetConversationID,
	},
	// 4.2 Get Conversation Participants
	route.Route{
		"V1ConversationParticipantsPOST",
		"POST",
		"/v1/conversation/participant",
		handlers.V1GetConversationParticipants,
	},

	//5.1 Grade List
	route.Route{
		"V1GradeListPOST",
		"POST",
		"/v1/grade/list",
		handlers.V1GradeList,
	},
	route.Route{
		"V1GradeListGET",
		"GET",
		"/v1/grade/list",
		handlers.V1GradeList,
	},

	//5.2 Subject List
	route.Route{
		"V1SubjectListPOST",
		"POST",
		"/v1/subject/list",
		handlers.V1SubjectList,
	},
	route.Route{
		"V1SubjectListGET",
		"GET",
		"/v1/subject/list",
		handlers.V1SubjectList,
	},

	// 5.3 Create Order
	route.Route{
		"V1OrderCreatePOST",
		"POST",
		"/v1/order/create",
		handlers.V1OrderCreate,
	},
	route.Route{
		"V1OrderCreateGET",
		"GET",
		"/v1/order/create",
		handlers.V1OrderCreate,
	},

	//5.4 Personal Order Confirm
	route.Route{
		"V1OrderPersonalConfirmPOST",
		"POST",
		"/v1/order/personal/confirm",
		handlers.V1OrderPersonalConfirm,
	},
	route.Route{
		"V1OrderPersonalConfirmGET",
		"GET",
		"/v1/order/personal/confirm",
		handlers.V1OrderPersonalConfirm,
	},

	//5.5 Teacher Expect Price
	route.Route{
		"V1TeacherExpectPost",
		"POST",
		"/v1/teacher/expect",
		handlers.V1TeacherExpect,
	},
	route.Route{
		"V1TeacherExpectGET",
		"GET",
		"/v1/teacher/expect",
		handlers.V1TeacherExpect,
	},
	// 5.6 Create RealTime Order
	route.Route{
		"V1RealTimeOrderCreatePOST",
		"POST",
		"/v1/order/realtime/create",
		handlers.V1RealTimeOrderCreate,
	},
	route.Route{
		"V1RealTimeOrderCreateGET",
		"GET",
		"/v1/order/realtime/create",
		handlers.V1RealTimeOrderCreate,
	},
	//5.7 RealTime Order Confirm
	route.Route{
		"V1RealTimeOrderConfirmPOST",
		"POST",
		"/v1/order/realtime/confirm",
		handlers.V1RealTimeOrderConfirm,
	},
	route.Route{
		"V1RealTimeOrderConfirmGET",
		"GET",
		"/v1/order/realtime/confirm",
		handlers.V1RealTimeOrderConfirm,
	},

	//6.1 Trade Charge
	route.Route{
		"V1TradeChargePOST",
		"POST",
		"/v1/trade/charge",
		handlers.V1TradeCharge,
	},
	route.Route{
		"V1TradeChargeGET",
		"GET",
		"/v1/trade/charge",
		handlers.V1TradeCharge,
	},

	//6.2 Trade Withdraw
	route.Route{
		"V1TradeWithdrawPOST",
		"POST",
		"/v1/trade/withdraw",
		handlers.V1TradeWithdraw,
	},
	route.Route{
		"V1TradeWithdrawGET",
		"GET",
		"/v1/trade/withdraw",
		handlers.V1TradeWithdraw,
	},

	//6.3 Trade Award
	route.Route{
		"V1TradeAwardPOST",
		"POST",
		"/v1/trade/award",
		handlers.V1TradeAward,
	},
	route.Route{
		"V1TradeAwardGET",
		"GET",
		"/v1/trade/award",
		handlers.V1TradeAward,
	},

	//6.4 Trade Promotion
	route.Route{
		"V1TradePromotionPOST",
		"POST",
		"/v1/trade/promotion",
		handlers.V1TradePromotion,
	},
	route.Route{
		"V1TradePromotionGET",
		"GET",
		"/v1/trade/promotion",
		handlers.V1TradePromotion,
	},

	// 6.5 User Trade Record
	route.Route{
		"V1MyTradeRecordGET",
		"GET",
		"/v1/trade/traderecord",
		handlers.V1TradeRecord,
	},
	route.Route{
		"V1MyTradeRecordPOST",
		"POST",
		"/v1/trade/traderecord",
		handlers.V1TradeRecord,
	},

	// 7.1 Complain
	route.Route{
		"V1ComplainPOST",
		"POST",
		"/v1/complaint/complain",
		handlers.V1Complain,
	},
	route.Route{
		"V1ComplainGET",
		"GET",
		"/v1/complaint/complain",
		handlers.V1Complain,
	},

	// 7.2 Handle Complaint
	route.Route{
		"V1HandleComplaintPOST",
		"POST",
		"/v1/complaint/handle",
		handlers.V1HandleComplaint,
	},
	route.Route{
		"V1HandleComplaintGET",
		"GET",
		"/v1/complaint/handle",
		handlers.V1HandleComplaint,
	},
	// 7.3 Check Complaint Exsits
	route.Route{
		"V1CheckComplaintExsitsPOST",
		"POST",
		"/v1/complaint/check",
		handlers.V1CheckComplaintExsits,
	},
	route.Route{
		"V1CheckComplaintExsitsGET",
		"GET",
		"/v1/complaint/check",
		handlers.V1CheckComplaintExsits,
	},

	// 8.1 Search Teachers
	route.Route{
		"V1SearchTeachersPOST",
		"POST",
		"/v1/search/teacher",
		handlers.V1SearchTeachers,
	},
	route.Route{
		"V1SearchTeachersGET",
		"GET",
		"/v1/search/teacher",
		handlers.V1SearchTeachers,
	},
	// 8.2 Search Users
	route.Route{
		"V1SearchUsersPOST",
		"POST",
		"/v1/search/user",
		handlers.V1SearchUsers,
	},
	route.Route{
		"V1SearchUsersGET",
		"GET",
		"/v1/search/user",
		handlers.V1SearchUsers,
	},

	// 9.1 Insert Evaluation
	route.Route{
		"V1EvaluatePOST",
		"POST",
		"/v1/evaluation/insert",
		handlers.V1Evaluate,
	},

	route.Route{
		"V1EvaluateGET",
		"GET",
		"/v1/evaluation/insert",
		handlers.V1Evaluate,
	},

	// 9.2 Query Evaluation
	route.Route{
		"V1GetEvaluationPOST",
		"POST",
		"/v1/evaluation/query",
		handlers.V1GetEvaluation,
	},
	route.Route{
		"V1GetEvaluationGET",
		"GET",
		"/v1/evaluation/query",
		handlers.V1GetEvaluation,
	},

	// 9.3 Query Evaluation Labels
	route.Route{
		"V1GetEvaluationLabelPOST",
		"POST",
		"/v1/evaluation/label",
		handlers.V1GetEvaluationLabels,
	},
	route.Route{
		"V1GetEvaluationLabelGET",
		"GET",
		"/v1/evaluation/label",
		handlers.V1GetEvaluationLabels,
	},

	// 10.1 Activity Notification
	route.Route{
		"V1GetActivitiesPOST",
		"POST",
		"/v1/activity/notification",
		handlers.V1ActivityNotification,
	},
	route.Route{
		"V1GetActivitiesGET",
		"GET",
		"/v1/activity/notification",
		handlers.V1ActivityNotification,
	},

	// 11.1 Bind User with InvitatoinCode
	route.Route{
		"V1BindInvitationCodePOST",
		"POST",
		"/v1/invitation/bind",
		handlers.V1BindUserWithInvitationCode,
	},
	route.Route{
		"V1BindInvitationCodeGET",
		"GET",
		"/v1/invitation/bind",
		handlers.V1BindUserWithInvitationCode,
	},

	// 11.2 Check User Has binded with invitationCode
	route.Route{
		"V1CheckUserBindWithInvitationCodePOST",
		"POST",
		"/v1/invitation/check",
		handlers.V1CheckUserHasBindWithInvitationCode,
	},
	route.Route{
		"V1CheckUserBindWithInvitationCodeGET",
		"GET",
		"/v1/invitation/check",
		handlers.V1CheckUserHasBindWithInvitationCode,
	},

	// 12.1 GetCourses
	route.Route{
		"V1GetCoursesPOST",
		"POST",
		"/v1/course/list",
		handlers.V1GetCourses,
	},
	route.Route{
		"V1GetCoursesGET",
		"GET",
		"/v1/course/list",
		handlers.V1GetCourses,
	},
	// 12.2 Join course
	route.Route{
		"V1JoinCoursesPOST",
		"POST",
		"/v1/course/join",
		handlers.V1JoinCourse,
	},
	route.Route{
		"V1JoinCoursesGET",
		"GET",
		"/v1/course/join",
		handlers.V1JoinCourse,
	},
	// 12.3 Active course
	route.Route{
		"V1ActiveCoursesPOST",
		"POST",
		"/v1/course/active",
		handlers.V1ActiveCourse,
	},
	route.Route{
		"V1ActiveCourseGET",
		"GET",
		"/v1/course/active",
		handlers.V1ActiveCourse,
	},
	// 12.4 User Renew course
	route.Route{
		"V1RenewCoursePOST",
		"POST",
		"/v1/course/renew",
		handlers.V1RenewCourse,
	},
	route.Route{
		"V1RenewCoursesGET",
		"GET",
		"/v1/course/renew",
		handlers.V1RenewCourse,
	},
	// 12.5 Support Renew course
	route.Route{
		"V1SupportRenewCoursePOST",
		"POST",
		"/v1/course/support/renew",
		handlers.V1SupportRenewCourse,
	},
	route.Route{
		"V1SupportRenewCourseGET",
		"GET",
		"/v1/course/support/renew",
		handlers.V1SupportRenewCourse,
	},
	// 12.6 Support reject course apply
	route.Route{
		"V1SupportRejectCoursePOST",
		"POST",
		"/v1/course/support/reject",
		handlers.V1SupportRejectCourse,
	},
	route.Route{
		"V1SupportRejectCourseGET",
		"GET",
		"/v1/course/support/reject",
		handlers.V1SupportRejectCourse,
	},

	// 13.1 Insert experience
	route.Route{
		"V1InsertExperiencePOST",
		"POST",
		"/v1/experience/insert",
		handlers.V1InsertExperience,
	},
	route.Route{
		"V1InsertExperienceGET",
		"GET",
		"/v1/experience/insert",
		handlers.V1InsertExperience,
	},

	// 14.1 pingpp pay
	route.Route{
		"V1PayByPingppPOST",
		"POST",
		"/v1/pingpp/pay",
		handlers.V1PayByPingpp,
	},
	route.Route{
		"V1PayByPingppGET",
		"GET",
		"/v1/pingpp/pay",
		handlers.V1PayByPingpp,
	},
	// 14.2 pingpp refund
	route.Route{
		"V1RefundByPingppPOST",
		"POST",
		"/v1/pingpp/refund",
		handlers.V1RefundByPingpp,
	},
	route.Route{
		"V1RefundByPingppGET",
		"GET",
		"/v1/pingpp/refund",
		handlers.V1RefundByPingpp,
	},
	// 14.3 pingpp query payment
	route.Route{
		"V1QueryPaymentByPingppPOST",
		"POST",
		"/v1/pingpp/pay/query",
		handlers.V1QueryPaymentByPingpp,
	},
	route.Route{
		"V1QueryPaymentByPingppGET",
		"GET",
		"/v1/pingpp/pay/query",
		handlers.V1QueryPaymentByPingpp,
	},
	// 14.4 pingpp query payment list
	route.Route{
		"V1QueryPaymentListByPingppPOST",
		"POST",
		"/v1/pingpp/pay/list",
		handlers.V1QueryPaymentListByPingpp,
	},
	route.Route{
		"V1QueryPaymentListByPingppGET",
		"GET",
		"/v1/pingpp/pay/list",
		handlers.V1QueryPaymentListByPingpp,
	},
	// 14.5 pingpp query refund
	route.Route{
		"V1QueryRefundByPingppPOST",
		"POST",
		"/v1/pingpp/refund/query",
		handlers.V1QueryRefundByPingpp,
	},
	route.Route{
		"V1QueryRefundByPingppGET",
		"GET",
		"/v1/pingpp/refund/query",
		handlers.V1QueryRefundByPingpp,
	},
	// 14.6 pingpp query refund list
	route.Route{
		"V1QueryRefundListByPingppPOST",
		"POST",
		"/v1/pingpp/refund/list",
		handlers.V1QueryRefundListByPingpp,
	},
	route.Route{
		"V1QueryRefundListByPingppGET",
		"GET",
		"/v1/pingpp/refund/list",
		handlers.V1QueryRefundListByPingpp,
	},
	// 14.7 pingpp webhook
	route.Route{
		"V1WebhookByPingppPOST",
		"POST",
		"/v1/pingpp/webhook",
		handlers.V1WebhookByPingpp,
	},
	route.Route{
		"V1WebhookByPingppGET",
		"GET",
		"/v1/pingpp/webhook",
		handlers.V1WebhookByPingpp,
	},
	// 14.8 pingpp result
	route.Route{
		"V1GetPingResultPOST",
		"POST",
		"/v1/pingpp/result",
		handlers.V1GetPingppResult,
	},
	route.Route{
		"V1GetPingResultGET",
		"GET",
		"/v1/pingpp/result",
		handlers.V1GetPingppResult,
	},
	// 15.1 sendcloud smshook
	route.Route{
		"V1WSMSHookPOST",
		"POST",
		"/v1/sendcloud/smshook",
		handlers.V1SmsHook,
	},
	// 15.2 senccloud sendmessage
	route.Route{
		"V1SCSendMessagePOST",
		"POST",
		"/v1/sendcloud/sendmessage",
		handlers.V1SendMessage,
	},

	route.Route{
		"V1BannerPOST",
		"POST",
		"/v1/banner",
		handlers.V1Banner,
	},
	route.Route{
		"V1BannerGET",
		"GET",
		"/v1/banner",
		handlers.V1Banner,
	},

	route.Route{
		"V1StatusLivePOST",
		"POST",
		"/v1/status/live",
		handlers.V1StatusLive,
	},
	route.Route{
		"V1StatusLiveGET",
		"GET",
		"/v1/status/live",
		handlers.V1StatusLive,
	},

	route.Route{
		"V1SendAdvMessagePOST",
		"POST",
		"/v1/send/adv",
		handlers.V1SendAdvMessage,
	},
	route.Route{
		"V1SendAdvMessageGET",
		"GET",
		"/v1/send/adv",
		handlers.V1SendAdvMessage,
	},
	route.Route{
		"V1GetHelpItemsPOST",
		"POST",
		"/v1/help/get",
		handlers.V1GetHelpItems,
	},
	route.Route{
		"V1GetHelpItemsGET",
		"GET",
		"/v1/help/get",
		handlers.V1GetHelpItems,
	},

	//---------POI Monitor------//
	route.Route{
		"V1MonitorUserPOST",
		"POST",
		"/v1/monitor/user",
		handlers.GetUserMonitorInfo,
	},
	route.Route{
		"V1MonitorUserGET",
		"GET",
		"/v1/monitor/user",
		handlers.GetUserMonitorInfo,
	},
	route.Route{
		"V1MonitorOrderPOST",
		"POST",
		"/v1/monitor/order",
		handlers.GetOrderMonitorInfo,
	},
	route.Route{
		"V1MonitorOrderGET",
		"GET",
		"/v1/monitor/order",
		handlers.GetOrderMonitorInfo,
	},
	route.Route{
		"V1MonitorSessionPOST",
		"POST",
		"/v1/monitor/session",
		handlers.GetSessionMonitorInfo,
	},
	route.Route{
		"V1MonitorSessionGET",
		"GET",
		"/v1/monitor/session",
		handlers.GetSessionMonitorInfo,
	},
	route.Route{
		"V1PostSeekHelpPOST",
		"POST",
		"/v1/seekhelp/post",
		handlers.V1SetSeekHelp,
	},
	route.Route{
		"V1PostSeekHelpGET",
		"GET",
		"/v1/seekhelp/post",
		handlers.V1SetSeekHelp,
	},
	route.Route{
		"V1GetSeekHelpsPOST",
		"POST",
		"/v1/seekhelp/get",
		handlers.V1GetSeekHelps,
	},
	route.Route{
		"V1GetSeekHelpsGET",
		"GET",
		"/v1/seekhelp/get",
		handlers.V1GetSeekHelps,
	},
	route.Route{
		"V1GetSeekHelpsCountGET",
		"GET",
		"/v1/seekhelp/count",
		handlers.V1GetSeekHelpsCount,
	},
	route.Route{
		"V1GetSeekHelpsCountPOST",
		"POST",
		"/v1/seekhelp/count",
		handlers.V1GetSeekHelpsCount,
	},
	route.Route{
		"V1GetMessageLogsPOST",
		"POST",
		"/v1/messagelog/list",
		handlers.V1GetMessageLogs,
	},
	route.Route{
		"V1GetMessageLogsGET",
		"GET",
		"/v1/messagelog/list",
		handlers.V1GetMessageLogs,
	},
	route.Route{
		"V1GetMessageLogsCountPOST",
		"POST",
		"/v1/messagelog/count",
		handlers.V1GetMessageLogsCount,
	},
	route.Route{
		"V1GetMessageLogsCountGET",
		"GET",
		"/v1/messagelog/count",
		handlers.V1GetMessageLogsCount,
	},

	// Dummy
	route.Route{
		"Dummy",
		"GET",
		"/dummy",
		handlers.Dummy,
	},
	route.Route{
		"Dummy2",
		"GET",
		"/dummy2",
		handlers.Dummy2,
	},
	route.Route{
		"TestGET",
		"GET",
		"/test",
		handlers.Test,
	},
}
