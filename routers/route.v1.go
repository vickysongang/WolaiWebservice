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
	// 15.3 sendcloud verify rand code
	route.Route{
		"V1SCSendMessagePOST",
		"POST",
		"/v1/sendcloud/verify",
		handlers.V1VerifyRandCode,
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
