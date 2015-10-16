// V2Routes
package routers

import (
	"POIWolaiWebService/handlers"
	"POIWolaiWebService/websocket"
)

var V2Routes = Routes{
	// Websocket
	Route{
		"V2WebSocketHandler",
		"GET",
		"/v2/ws",
		websocket.V1WebSocketHandler,
	},

	// 1.1 Login
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

	// 1.2 Update profile
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

	// 1.3 Oauth Login
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

	// 1.4 Oauth Register
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

	// 1.6 Teacher Recommendation
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

	// 1.7 Teacher Profile
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
	//1.8 Teacher post
	Route{
		"V2InsertTeacherPost",
		"POST",
		"/v2/teacher/insert",
		handlers.V2TeacherPost,
	},
	//1.9 Check Phone
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

	// 2.1 Atrium
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

	// 2.2 Feed Post
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

	// 2.3 Feed Detail
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

	// 2.4 Feed Like
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

	// 2.5 Feed Favorite
	/*
		Route{
			"V2FeedFavPOST",
			"POST",
			"/v2/feed/favorite",
			V2FeedFav,
		},
		Route{
			"V2FeedFavGET",
			"GET",
			"/v2/feed/favorite",
			V2FeedFav,
		},
	*/

	// 2.6 Feed Comment
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

	// 2.7 Feed Comment Like
	/*
		Route{
			"V2FeedCommentLikePOST",
			"POST",
			"/v2/feed/comment/like",
			V2FeedCommentLike,
		},
		Route{
			"V2FeedCommentLikeGET",
			"GET",
			"/v2/feed/comment/like",
			V2FeedCommentLike,
		},
	*/
	// 2.8 Feed mark
	Route{
		"V2MarkFeedPOST",
		"POST",
		"/v2/feed/mark",
		handlers.V2FeedMark,
	},
	Route{
		"V2MarkFeedGET",
		"GET",
		"/v2/feed/mark",
		handlers.V2FeedMark,
	},

	// 3.1 User Info
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

	// 3.2 User Wallet
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

	// 3.3 User MyFeed
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

	// 3.4 User MyFollowing
	Route{
		"V2UserMyFollowPOST",
		"POST",
		"/v2/user/myfollow",
		handlers.V2UserMyFollowing,
	},
	Route{
		"V2UserMyFollowGET",
		"GET",
		"/v2/user/myfollow",
		handlers.V2UserMyFollowing,
	},

	// 3.5 User MyLike
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

	// 3.6 User Follow
	Route{
		"V2UserFollowPOST",
		"POST",
		"/v2/user/follow",
		handlers.V2UserFollow,
	},
	Route{
		"V2UserFollowGET",
		"GET",
		"/v2/user/follow",
		handlers.V2UserFollow,
	},

	// 3.7 User Unfollow
	Route{
		"V2UserUnfollowPOST",
		"POST",
		"/v2/user/unfollow",
		handlers.V2UserUnfollow,
	},
	Route{
		"V2UserUnfollowGET",
		"GET",
		"/v2/user/unfollow",
		handlers.V2UserUnfollow,
	},

	// 3.8 User Order
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

	// 4.1 Get Conversation ID
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

	//5.1 Grade List
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

	//5.2 Subject List
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

	// 5.3 Create Order
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

	//5.4 Personal Order Confirm
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

	//5.5 Teacher Expect Price
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

	//6.1 Trade Charge
	Route{
		"V2TradeChargePOST",
		"POST",
		"/v2/trade/charge",
		handlers.V2TradeCharge,
	},
	Route{
		"V2TradeChargeGET",
		"GET",
		"/v2/trade/charge",
		handlers.V2TradeCharge,
	},

	//6.2 Trade Withdraw
	Route{
		"V2TradeWithdrawPOST",
		"POST",
		"/v2/trade/withdraw",
		handlers.V2TradeWithdraw,
	},
	Route{
		"V2TradeWithdrawGET",
		"GET",
		"/v2/trade/withdraw",
		handlers.V2TradeWithdraw,
	},

	//6.3 Trade Award
	Route{
		"V2TradeAwardPOST",
		"POST",
		"/v2/trade/award",
		handlers.V2TradeAward,
	},
	Route{
		"V2TradeAwardGET",
		"GET",
		"/v2/trade/award",
		handlers.V2TradeAward,
	},

	//6.4 Trade Promotion
	Route{
		"V2TradePromotionPOST",
		"POST",
		"/v2/trade/promotion",
		handlers.V2TradePromotion,
	},
	Route{
		"V2TradePromotionGET",
		"GET",
		"/v2/trade/promotion",
		handlers.V2TradePromotion,
	},

	// 6.5 User Trade Record
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

	// 7.1 Complain
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

	// 7.2 Handle Complaint
	Route{
		"V2HandleComplaintPOST",
		"POST",
		"/v2/complaint/handle",
		handlers.V2HandleComplaint,
	},
	Route{
		"V2HandleComplaintGET",
		"GET",
		"/v2/complaint/handle",
		handlers.V2HandleComplaint,
	},

	// 8.1 Search Teachers
	Route{
		"V2SearchTeacherPOST",
		"POST",
		"/v2/search/teacher",
		handlers.V2SearchTeacher,
	},
	Route{
		"V2SearchTeacherGET",
		"GET",
		"/v2/search/teacher",
		handlers.V2SearchTeacher,
	},

	// 9.1 Insert Evaluation
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

	// 9.2 Query Evaluation
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

	// 9.3 Query Evaluation Labels
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

	// 10.1 Activity Notification
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

	// 11.1 Bind User with InvitatoinCode
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

	// 11.2 Check User Has binded with invitationCode
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

	// 12.1 GetCourses
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
	// 12.2 Join course
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
	// 12.3 Active course
	Route{
		"V2ActiveCoursesPOST",
		"POST",
		"/v2/course/active",
		handlers.V2ActiveCourse,
	},
	Route{
		"V2ActiveCourseGET",
		"GET",
		"/v2/course/active",
		handlers.V2ActiveCourse,
	},
	// 12.4 User Renew course
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
	// 12.5 Support Renew course
	Route{
		"V2SupportRenewCoursePOST",
		"POST",
		"/v2/course/support/renew",
		handlers.V2SupportRenewCourse,
	},
	Route{
		"V2SupportRenewCourseGET",
		"GET",
		"/v2/course/support/renew",
		handlers.V2SupportRenewCourse,
	},
	// 12.6 Support reject course apply
	Route{
		"V2SupportRejectCoursePOST",
		"POST",
		"/v2/course/support/reject",
		handlers.V2SupportRejectCourse,
	},
	Route{
		"V2SupportRejectCourseGET",
		"GET",
		"/v2/course/support/reject",
		handlers.V2SupportRejectCourse,
	},

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

	Route{
		"V2StatusLivePOST",
		"POST",
		"/v2/status/live",
		handlers.V2StatusLive,
	},
	Route{
		"V2StatusLiveGET",
		"GET",
		"/v2/status/live",
		handlers.V2StatusLive,
	},

	Route{
		"V2ConversationParticipantsPOST",
		"POST",
		"/v2/conversation/participant",
		handlers.V2GetConversationParticipants,
	},

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

	//---------POI Monitor------//
	Route{
		"V2MonitorUserPOST",
		"POST",
		"/v2/monitor/user",
		handlers.GetUserMonitorInfo,
	},
	Route{
		"V2MonitorUserGET",
		"GET",
		"/v2/monitor/user",
		handlers.GetUserMonitorInfo,
	},
	Route{
		"V2MonitorOrderPOST",
		"POST",
		"/v2/monitor/order",
		handlers.GetOrderMonitorInfo,
	},
	Route{
		"V2MonitorOrderGET",
		"GET",
		"/v2/monitor/order",
		handlers.GetOrderMonitorInfo,
	},
}
