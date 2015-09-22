package routers

import (
	"encoding/json"
	"net/http"
	"time"

	"POIWolaiWebService/handlers"
	"POIWolaiWebService/websocket"

	seelog "github.com/cihub/seelog"
	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = WebLogger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}

func WebLogger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		formData, _ := json.Marshal(r.Form)
		seelog.Info("[", r.Method, "] ", r.RequestURI, "\t", name, "\t", time.Since(start),
			"\t", string(formData))
	})
}

var routes = Routes{
	// Websocket
	Route{
		"V1WebSocketHandler",
		"GET",
		"/v1/ws",
		websocket.V1WebSocketHandler,
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
		"v1InsertTeacherPost",
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
		"v1CheckPhonePost",
		"POST",
		"/v1/oauth/qq/checkphone",
		handlers.V1CheckPhoneBindWithQQ,
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
		"v1MyOrdersGET",
		"GET",
		"/v1/user/myorders",
		handlers.V1OrderInSession,
	},
	Route{
		"v1MyOrdersPost",
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
		"v1MyTradeRecordGET",
		"GET",
		"/v1/trade/traderecord",
		handlers.V1TradeRecord,
	},
	Route{
		"v1MyTradeRecordPOST",
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

	// 8.1 Search Teachers
	Route{
		"V1SearchTeacherPOST",
		"POST",
		"/v1/search/teacher",
		handlers.V1SearchTeacher,
	},
	Route{
		"V1SearchTeacherGET",
		"GET",
		"/v1/search/teacher",
		handlers.V1SearchTeacher,
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
		"V1ConversationParticipantsPOST",
		"POST",
		"/v1/conversation/participant",
		handlers.V1GetConversationParticipants,
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
		"V1GetHelpCheatsPOST",
		"POST",
		"/v1/help/get",
		handlers.V1GetHelpItems,
	},
	Route{
		"V1GetHelpCheatsGET",
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
}
