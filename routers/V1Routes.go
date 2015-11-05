// V1Routes
package routers

import (
	"POIWolaiWebService/handlers"
	"POIWolaiWebService/websocket"
)

var V1Routes = Routes{
	// Websocket
	Route{
		"V1WebSocketHandler",
		"GET",
		"/v1/ws",
		websocket.V1WebSocketHandler,
	},

	// 1.1 Login
	Route{
		"V1LoginPOST",
		"POST",
		"/v1/login",
		handlers.V1Login,
	},
	Route{
		"V1LoginGET",
		"GET",
		"/v1/login",
		handlers.V1Login,
	},
	Route{
		"V1LoginGETURL",
		"GET",
		"/v1/login/{phone}",
		handlers.V1LoginGETURL,
	},

	// 1.2 Update profile
	Route{
		"V1UpdateProfilePOST",
		"POST",
		"/v1/update_profile",
		handlers.V1UpdateProfile,
	},
	Route{
		"V1UpdateProfileGET",
		"GET",
		"/v1/update_profile",
		handlers.V1UpdateProfile,
	},
	Route{
		"V1UpdateProfileGETURL",
		"GET",
		"/v1/update_profile/{userId}/{nickname}/{avatar}/{gender}",
		handlers.V1UpdateProfileGETURL,
	},

	// 1.3 Oauth Login
	Route{
		"V1OauthLoginPOST",
		"POST",
		"/v1/oauth/qq/login",
		handlers.V1OauthLogin,
	},
	Route{
		"V1OauthLoginGET",
		"GET",
		"/v1/oauth/qq/login",
		handlers.V1OauthLogin,
	},

	// 1.4 Oauth Register
	Route{
		"V1OauthRegisterPOST",
		"POST",
		"/v1/oauth/qq/register",
		handlers.V1OauthRegister,
	},
	Route{
		"V1OauthRegisterGET",
		"GET",
		"/v1/oauth/qq/register",
		handlers.V1OauthRegister,
	},

	// 1.6 Teacher Recommendation
	Route{
		"V1TeacherRecommendationPOST",
		"POST",
		"/v1/teacher/recommendation",
		handlers.V1TeacherRecommendation,
	},
	Route{
		"V1TeacherRecommendationGET",
		"GET",
		"/v1/teacher/recommendation",
		handlers.V1TeacherRecommendation,
	},

	// 1.7 Teacher Profile
	Route{
		"V1TeacherProfilePOST",
		"POST",
		"/v1/teacher/profile",
		handlers.V1TeacherProfile,
	},
	Route{
		"V1TeacherProfileGET",
		"GET",
		"/v1/teacher/profile",
		handlers.V1TeacherProfile,
	},
	//1.8 Teacher post
	Route{
		"V1InsertTeacherPost",
		"POST",
		"/v1/teacher/insert",
		handlers.V1TeacherPost,
	},
	//1.9 Check Phone
	Route{
		"V1CheckPhoneGET",
		"GET",
		"/v1/oauth/qq/checkphone",
		handlers.V1CheckPhoneBindWithQQ,
	},
	Route{
		"V1CheckPhonePost",
		"POST",
		"/v1/oauth/qq/checkphone",
		handlers.V1CheckPhoneBindWithQQ,
	},
	// 1.10 support and teacher list
	Route{
		"V1SupportAndTeacherListPOST",
		"POST",
		"/v1/support/teacher/list",
		handlers.V1SupportAndTeacherList,
	},
	Route{
		"V1TSupportAndTeacherListGET",
		"GET",
		"/v1/support/teacher/list",
		handlers.V1SupportAndTeacherList,
	},

	// 2.1 Atrium
	Route{
		"V1AtriumPOST",
		"POST",
		"/v1/atrium",
		handlers.V1Atrium,
	},
	Route{
		"V1AtriumGET",
		"GET",
		"/v1/atrium",
		handlers.V1Atrium,
	},

	// 2.2 Feed Post
	Route{
		"V1FeedPostPOST",
		"POST",
		"/v1/feed/post",
		handlers.V1FeedPost,
	},
	Route{
		"V1FeedPostGET",
		"GET",
		"/v1/feed/post",
		handlers.V1FeedPost,
	},

	// 2.3 Feed Detail
	Route{
		"V1FeedDetailPOST",
		"POST",
		"/v1/feed/detail",
		handlers.V1FeedDetail,
	},
	Route{
		"V1FeedDetailGET",
		"GET",
		"/v1/feed/detail",
		handlers.V1FeedDetail,
	},

	// 2.4 Feed Like
	Route{
		"V1FeedLikePOST",
		"POST",
		"/v1/feed/like",
		handlers.V1FeedLike,
	},
	Route{
		"V1FeedLikeGET",
		"GET",
		"/v1/feed/like",
		handlers.V1FeedLike,
	},

	// 2.5 Feed Favorite
	/*
		Route{
			"V1FeedFavPOST",
			"POST",
			"/v1/feed/favorite",
			V1FeedFav,
		},
		Route{
			"V1FeedFavGET",
			"GET",
			"/v1/feed/favorite",
			V1FeedFav,
		},
	*/

	// 2.6 Feed Comment
	Route{
		"V1FeedCommentPOST",
		"POST",
		"/v1/feed/comment",
		handlers.V1FeedComment,
	},
	Route{
		"V1FeedCommentGET",
		"GET",
		"/v1/feed/comment",
		handlers.V1FeedComment,
	},

	// 2.7 Feed Comment Like
	/*
		Route{
			"V1FeedCommentLikePOST",
			"POST",
			"/v1/feed/comment/like",
			V1FeedCommentLike,
		},
		Route{
			"V1FeedCommentLikeGET",
			"GET",
			"/v1/feed/comment/like",
			V1FeedCommentLike,
		},
	*/
	// 2.8 Feed mark
	Route{
		"V1MarkFeedPOST",
		"POST",
		"/v1/feed/mark",
		handlers.V1FeedMark,
	},
	Route{
		"V1MarkFeedGET",
		"GET",
		"/v1/feed/mark",
		handlers.V1FeedMark,
	},

	// 2.9 Top Feed
	Route{
		"V1TopFeedPOST",
		"POST",
		"/v1/feed/top/get",
		handlers.V1GetTopFeed,
	},
	Route{
		"V1TopFeedGET",
		"GET",
		"/v1/feed/top/get",
		handlers.V1GetTopFeed,
	},
	// 2.10 Delete Feed
	Route{
		"V1DeleteFeedPOST",
		"POST",
		"/v1/feed/delete",
		handlers.V1FeedDelete,
	},
	Route{
		"V1DeleteFeedGET",
		"GET",
		"/v1/feed/delete",
		handlers.V1FeedDelete,
	},
	// 2.11 Recover Feed
	Route{
		"V1RecoverFeedPOST",
		"POST",
		"/v1/feed/recover",
		handlers.V1FeedRecover,
	},
	Route{
		"V1RecoverFeedGET",
		"GET",
		"/v1/feed/recover",
		handlers.V1FeedRecover,
	},
	// 2.12 Top Feed
	Route{
		"V1TopFeedPOST",
		"POST",
		"/v1/feed/top",
		handlers.V1FeedTop,
	},
	Route{
		"V1TopFeedGET",
		"GET",
		"/v1/feed/top",
		handlers.V1FeedTop,
	},

	// 3.1 User Info
	Route{
		"V1UserInfoPOST",
		"POST",
		"/v1/user/info",
		handlers.V1UserInfo,
	},
	Route{
		"V1UserInfoGET",
		"GET",
		"/v1/user/info",
		handlers.V1UserInfo,
	},

	// 3.2 User Wallet
	Route{
		"V1UserMyWalletPOST",
		"POST",
		"/v1/user/mywallet",
		handlers.V1UserMyWallet,
	},
	Route{
		"V1UserMyWalletGET",
		"GET",
		"/v1/user/mywallet",
		handlers.V1UserMyWallet,
	},

	// 3.3 User MyFeed
	Route{
		"V1UserMyFeedPOST",
		"POST",
		"/v1/user/myfeed",
		handlers.V1UserMyFeed,
	},
	Route{
		"V1UserMyFeedGET",
		"GET",
		"/v1/user/myfeed",
		handlers.V1UserMyFeed,
	},

	// 3.4 User MyFollowing
	Route{
		"V1UserMyFollowPOST",
		"POST",
		"/v1/user/myfollow",
		handlers.V1UserMyFollowing,
	},
	Route{
		"V1UserMyFollowGET",
		"GET",
		"/v1/user/myfollow",
		handlers.V1UserMyFollowing,
	},

	// 3.5 User MyLike
	Route{
		"V1UserMyLikePOST",
		"POST",
		"/v1/user/mylike",
		handlers.V1UserMyLike,
	},
	Route{
		"V1UserMyLikeGET",
		"GET",
		"/v1/user/mylike",
		handlers.V1UserMyLike,
	},

	// 3.6 User Follow
	Route{
		"V1UserFollowPOST",
		"POST",
		"/v1/user/follow",
		handlers.V1UserFollow,
	},
	Route{
		"V1UserFollowGET",
		"GET",
		"/v1/user/follow",
		handlers.V1UserFollow,
	},

	// 3.7 User Unfollow
	Route{
		"V1UserUnfollowPOST",
		"POST",
		"/v1/user/unfollow",
		handlers.V1UserUnfollow,
	},
	Route{
		"V1UserUnfollowGET",
		"GET",
		"/v1/user/unfollow",
		handlers.V1UserUnfollow,
	},

	// 3.8 User Order
	Route{
		"V1MyOrdersGET",
		"GET",
		"/v1/user/myorders",
		handlers.V1OrderInSession,
	},
	Route{
		"V1MyOrdersPost",
		"POST",
		"/v1/user/myorders",
		handlers.V1OrderInSession,
	},

	// 4.1 Get Conversation ID
	Route{
		"V1GetConversationIDPOST",
		"POST",
		"/v1/conversation/get",
		handlers.V1GetConversationID,
	},
	Route{
		"V1GetConversationIDGET",
		"GET",
		"/v1/conversation/get",
		handlers.V1GetConversationID,
	},
	// 4.2 Get Conversation Participants
	Route{
		"V1ConversationParticipantsPOST",
		"POST",
		"/v1/conversation/participant",
		handlers.V1GetConversationParticipants,
	},

	//5.1 Grade List
	Route{
		"V1GradeListPOST",
		"POST",
		"/v1/grade/list",
		handlers.V1GradeList,
	},
	Route{
		"V1GradeListGET",
		"GET",
		"/v1/grade/list",
		handlers.V1GradeList,
	},

	//5.2 Subject List
	Route{
		"V1SubjectListPOST",
		"POST",
		"/v1/subject/list",
		handlers.V1SubjectList,
	},
	Route{
		"V1SubjectListGET",
		"GET",
		"/v1/subject/list",
		handlers.V1SubjectList,
	},

	// 5.3 Create Order
	Route{
		"V1OrderCreatePOST",
		"POST",
		"/v1/order/create",
		handlers.V1OrderCreate,
	},
	Route{
		"V1OrderCreateGET",
		"GET",
		"/v1/order/create",
		handlers.V1OrderCreate,
	},

	//5.4 Personal Order Confirm
	Route{
		"V1OrderPersonalConfirmPOST",
		"POST",
		"/v1/order/personal/confirm",
		handlers.V1OrderPersonalConfirm,
	},
	Route{
		"V1OrderPersonalConfirmGET",
		"GET",
		"/v1/order/personal/confirm",
		handlers.V1OrderPersonalConfirm,
	},

	//5.5 Teacher Expect Price
	Route{
		"V1TeacherExpectPost",
		"POST",
		"/v1/teacher/expect",
		handlers.V1TeacherExpect,
	},
	Route{
		"V1TeacherExpectGET",
		"GET",
		"/v1/teacher/expect",
		handlers.V1TeacherExpect,
	},
	// 5.6 Create RealTime Order
	Route{
		"V1RealTimeOrderCreatePOST",
		"POST",
		"/v1/order/realtime/create",
		handlers.V1RealTimeOrderCreate,
	},
	Route{
		"V1RealTimeOrderCreateGET",
		"GET",
		"/v1/order/realtime/create",
		handlers.V1RealTimeOrderCreate,
	},
	//5.7 RealTime Order Confirm
	Route{
		"V1RealTimeOrderConfirmPOST",
		"POST",
		"/v1/order/realtime/confirm",
		handlers.V1RealTimeOrderConfirm,
	},
	Route{
		"V1RealTimeOrderConfirmGET",
		"GET",
		"/v1/order/realtime/confirm",
		handlers.V1RealTimeOrderConfirm,
	},

	//6.1 Trade Charge
	Route{
		"V1TradeChargePOST",
		"POST",
		"/v1/trade/charge",
		handlers.V1TradeCharge,
	},
	Route{
		"V1TradeChargeGET",
		"GET",
		"/v1/trade/charge",
		handlers.V1TradeCharge,
	},

	//6.2 Trade Withdraw
	Route{
		"V1TradeWithdrawPOST",
		"POST",
		"/v1/trade/withdraw",
		handlers.V1TradeWithdraw,
	},
	Route{
		"V1TradeWithdrawGET",
		"GET",
		"/v1/trade/withdraw",
		handlers.V1TradeWithdraw,
	},

	//6.3 Trade Award
	Route{
		"V1TradeAwardPOST",
		"POST",
		"/v1/trade/award",
		handlers.V1TradeAward,
	},
	Route{
		"V1TradeAwardGET",
		"GET",
		"/v1/trade/award",
		handlers.V1TradeAward,
	},

	//6.4 Trade Promotion
	Route{
		"V1TradePromotionPOST",
		"POST",
		"/v1/trade/promotion",
		handlers.V1TradePromotion,
	},
	Route{
		"V1TradePromotionGET",
		"GET",
		"/v1/trade/promotion",
		handlers.V1TradePromotion,
	},

	// 6.5 User Trade Record
	Route{
		"V1MyTradeRecordGET",
		"GET",
		"/v1/trade/traderecord",
		handlers.V1TradeRecord,
	},
	Route{
		"V1MyTradeRecordPOST",
		"POST",
		"/v1/trade/traderecord",
		handlers.V1TradeRecord,
	},

	// 7.1 Complain
	Route{
		"V1ComplainPOST",
		"POST",
		"/v1/complaint/complain",
		handlers.V1Complain,
	},
	Route{
		"V1ComplainGET",
		"GET",
		"/v1/complaint/complain",
		handlers.V1Complain,
	},

	// 7.2 Handle Complaint
	Route{
		"V1HandleComplaintPOST",
		"POST",
		"/v1/complaint/handle",
		handlers.V1HandleComplaint,
	},
	Route{
		"V1HandleComplaintGET",
		"GET",
		"/v1/complaint/handle",
		handlers.V1HandleComplaint,
	},
	// 7.3 Check Complaint Exsits
	Route{
		"V1CheckComplaintExsitsPOST",
		"POST",
		"/v1/complaint/check",
		handlers.V1CheckComplaintExsits,
	},
	Route{
		"V1CheckComplaintExsitsGET",
		"GET",
		"/v1/complaint/check",
		handlers.V1CheckComplaintExsits,
	},

	// 8.1 Search Teachers
	Route{
		"V1SearchTeachersPOST",
		"POST",
		"/v1/search/teacher",
		handlers.V1SearchTeachers,
	},
	Route{
		"V1SearchTeachersGET",
		"GET",
		"/v1/search/teacher",
		handlers.V1SearchTeachers,
	},
	// 8.2 Search Users
	Route{
		"V1SearchUsersPOST",
		"POST",
		"/v1/search/user",
		handlers.V1SearchUsers,
	},
	Route{
		"V1SearchUsersGET",
		"GET",
		"/v1/search/user",
		handlers.V1SearchUsers,
	},

	// 9.1 Insert Evaluation
	Route{
		"V1EvaluatePOST",
		"POST",
		"/v1/evaluation/insert",
		handlers.V1Evaluate,
	},

	Route{
		"V1EvaluateGET",
		"GET",
		"/v1/evaluation/insert",
		handlers.V1Evaluate,
	},

	// 9.2 Query Evaluation
	Route{
		"V1GetEvaluationPOST",
		"POST",
		"/v1/evaluation/query",
		handlers.V1GetEvaluation,
	},
	Route{
		"V1GetEvaluationGET",
		"GET",
		"/v1/evaluation/query",
		handlers.V1GetEvaluation,
	},

	// 9.3 Query Evaluation Labels
	Route{
		"V1GetEvaluationLabelPOST",
		"POST",
		"/v1/evaluation/label",
		handlers.V1GetEvaluationLabels,
	},
	Route{
		"V1GetEvaluationLabelGET",
		"GET",
		"/v1/evaluation/label",
		handlers.V1GetEvaluationLabels,
	},

	// 10.1 Activity Notification
	Route{
		"V1GetActivitiesPOST",
		"POST",
		"/v1/activity/notification",
		handlers.V1ActivityNotification,
	},
	Route{
		"V1GetActivitiesGET",
		"GET",
		"/v1/activity/notification",
		handlers.V1ActivityNotification,
	},

	// 11.1 Bind User with InvitatoinCode
	Route{
		"V1BindInvitationCodePOST",
		"POST",
		"/v1/invitation/bind",
		handlers.V1BindUserWithInvitationCode,
	},
	Route{
		"V1BindInvitationCodeGET",
		"GET",
		"/v1/invitation/bind",
		handlers.V1BindUserWithInvitationCode,
	},

	// 11.2 Check User Has binded with invitationCode
	Route{
		"V1CheckUserBindWithInvitationCodePOST",
		"POST",
		"/v1/invitation/check",
		handlers.V1CheckUserHasBindWithInvitationCode,
	},
	Route{
		"V1CheckUserBindWithInvitationCodeGET",
		"GET",
		"/v1/invitation/check",
		handlers.V1CheckUserHasBindWithInvitationCode,
	},

	// 12.1 GetCourses
	Route{
		"V1GetCoursesPOST",
		"POST",
		"/v1/course/list",
		handlers.V1GetCourses,
	},
	Route{
		"V1GetCoursesGET",
		"GET",
		"/v1/course/list",
		handlers.V1GetCourses,
	},
	// 12.2 Join course
	Route{
		"V1JoinCoursesPOST",
		"POST",
		"/v1/course/join",
		handlers.V1JoinCourse,
	},
	Route{
		"V1JoinCoursesGET",
		"GET",
		"/v1/course/join",
		handlers.V1JoinCourse,
	},
	// 12.3 Active course
	Route{
		"V1ActiveCoursesPOST",
		"POST",
		"/v1/course/active",
		handlers.V1ActiveCourse,
	},
	Route{
		"V1ActiveCourseGET",
		"GET",
		"/v1/course/active",
		handlers.V1ActiveCourse,
	},
	// 12.4 User Renew course
	Route{
		"V1RenewCoursePOST",
		"POST",
		"/v1/course/renew",
		handlers.V1RenewCourse,
	},
	Route{
		"V1RenewCoursesGET",
		"GET",
		"/v1/course/renew",
		handlers.V1RenewCourse,
	},
	// 12.5 Support Renew course
	Route{
		"V1SupportRenewCoursePOST",
		"POST",
		"/v1/course/support/renew",
		handlers.V1SupportRenewCourse,
	},
	Route{
		"V1SupportRenewCourseGET",
		"GET",
		"/v1/course/support/renew",
		handlers.V1SupportRenewCourse,
	},
	// 12.6 Support reject course apply
	Route{
		"V1SupportRejectCoursePOST",
		"POST",
		"/v1/course/support/reject",
		handlers.V1SupportRejectCourse,
	},
	Route{
		"V1SupportRejectCourseGET",
		"GET",
		"/v1/course/support/reject",
		handlers.V1SupportRejectCourse,
	},

	// 13.1 Insert experience
	Route{
		"V1InsertExperiencePOST",
		"POST",
		"/v1/experience/insert",
		handlers.V1InsertExperience,
	},
	Route{
		"V1InsertExperienceGET",
		"GET",
		"/v1/experience/insert",
		handlers.V1InsertExperience,
	},

	// 14.1 pingpp pay
	Route{
		"V1PayByPingppPOST",
		"POST",
		"/v1/pingpp/pay",
		handlers.V1PayByPingpp,
	},
	Route{
		"V1PayByPingppGET",
		"GET",
		"/v1/pingpp/pay",
		handlers.V1PayByPingpp,
	},
	// 14.2 pingpp refund
	Route{
		"V1RefundByPingppPOST",
		"POST",
		"/v1/pingpp/refund",
		handlers.V1RefundByPingpp,
	},
	Route{
		"V1RefundByPingppGET",
		"GET",
		"/v1/pingpp/refund",
		handlers.V1RefundByPingpp,
	},
	// 14.3 pingpp query payment
	Route{
		"V1QueryPaymentByPingppPOST",
		"POST",
		"/v1/pingpp/pay/query",
		handlers.V1QueryPaymentByPingpp,
	},
	Route{
		"V1QueryPaymentByPingppGET",
		"GET",
		"/v1/pingpp/pay/query",
		handlers.V1QueryPaymentByPingpp,
	},
	// 14.4 pingpp query payment list
	Route{
		"V1QueryPaymentListByPingppPOST",
		"POST",
		"/v1/pingpp/pay/list",
		handlers.V1QueryPaymentListByPingpp,
	},
	Route{
		"V1QueryPaymentListByPingppGET",
		"GET",
		"/v1/pingpp/pay/list",
		handlers.V1QueryPaymentListByPingpp,
	},
	// 14.5 pingpp query refund
	Route{
		"V1QueryRefundByPingppPOST",
		"POST",
		"/v1/pingpp/refund/query",
		handlers.V1QueryRefundByPingpp,
	},
	Route{
		"V1QueryRefundByPingppGET",
		"GET",
		"/v1/pingpp/refund/query",
		handlers.V1QueryRefundByPingpp,
	},
	// 14.6 pingpp query refund list
	Route{
		"V1QueryRefundListByPingppPOST",
		"POST",
		"/v1/pingpp/refund/list",
		handlers.V1QueryRefundListByPingpp,
	},
	Route{
		"V1QueryRefundListByPingppGET",
		"GET",
		"/v1/pingpp/refund/list",
		handlers.V1QueryRefundListByPingpp,
	},
	// 14.7 pingpp webhook
	Route{
		"V1WebhookByPingppPOST",
		"POST",
		"/v1/pingpp/webhook",
		handlers.V1WebhookByPingpp,
	},
	Route{
		"V1WebhookByPingppGET",
		"GET",
		"/v1/pingpp/webhook",
		handlers.V1WebhookByPingpp,
	},

	// 15.1 location insert
	Route{
		"V1InsertLocationPOST",
		"POST",
		"/v1/location/insert",
		handlers.V1InsertLocation,
	},
	Route{
		"V1InsertLocatoinGET",
		"GET",
		"/v1/location/insert",
		handlers.V1InsertLocation,
	},

	Route{
		"V1BannerPOST",
		"POST",
		"/v1/banner",
		handlers.V1Banner,
	},
	Route{
		"V1BannerGET",
		"GET",
		"/v1/banner",
		handlers.V1Banner,
	},

	Route{
		"V1StatusLivePOST",
		"POST",
		"/v1/status/live",
		handlers.V1StatusLive,
	},
	Route{
		"V1StatusLiveGET",
		"GET",
		"/v1/status/live",
		handlers.V1StatusLive,
	},

	Route{
		"V1SendAdvMessagePOST",
		"POST",
		"/v1/send/adv",
		handlers.V1SendAdvMessage,
	},
	Route{
		"V1SendAdvMessageGET",
		"GET",
		"/v1/send/adv",
		handlers.V1SendAdvMessage,
	},
	Route{
		"V1GetHelpItemsPOST",
		"POST",
		"/v1/help/get",
		handlers.V1GetHelpItems,
	},
	Route{
		"V1GetHelpItemsGET",
		"GET",
		"/v1/help/get",
		handlers.V1GetHelpItems,
	},

	//---------POI Monitor------//
	Route{
		"V1MonitorUserPOST",
		"POST",
		"/v1/monitor/user",
		handlers.GetUserMonitorInfo,
	},
	Route{
		"V1MonitorUserGET",
		"GET",
		"/v1/monitor/user",
		handlers.GetUserMonitorInfo,
	},
	Route{
		"V1MonitorOrderPOST",
		"POST",
		"/v1/monitor/order",
		handlers.GetOrderMonitorInfo,
	},
	Route{
		"V1MonitorOrderGET",
		"GET",
		"/v1/monitor/order",
		handlers.GetOrderMonitorInfo,
	},
	Route{
		"V1MonitorSessionPOST",
		"POST",
		"/v1/monitor/session",
		handlers.GetSessionMonitorInfo,
	},
	Route{
		"V1MonitorSessionGET",
		"GET",
		"/v1/monitor/session",
		handlers.GetSessionMonitorInfo,
	},
	Route{
		"V1PostSeekHelpPOST",
		"POST",
		"/v1/seekhelp/post",
		handlers.V1SetSeekHelp,
	},
	Route{
		"V1PostSeekHelpGET",
		"GET",
		"/v1/seekhelp/post",
		handlers.V1SetSeekHelp,
	},
	Route{
		"V1GetSeekHelpsPOST",
		"POST",
		"/v1/seekhelp/get",
		handlers.V1GetSeekHelps,
	},
	Route{
		"V1GetSeekHelpsGET",
		"GET",
		"/v1/seekhelp/get",
		handlers.V1GetSeekHelps,
	},

	// Dummy
	Route{
		"Dummy",
		"GET",
		"/dummy",
		handlers.Dummy,
	},
	Route{
		"Dummy2",
		"GET",
		"/dummy2",
		handlers.Dummy2,
	},
	Route{
		"TestGET",
		"GET",
		"/test",
		handlers.Test,
	},
}
