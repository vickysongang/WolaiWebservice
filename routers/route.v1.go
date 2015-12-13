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
}
