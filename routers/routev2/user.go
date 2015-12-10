package routev2

import (
	"github.com/gorilla/mux"

	"WolaiWebservice/handlers/handlerv2"
	"WolaiWebservice/routers/route"
	"WolaiWebservice/routers/wrapper"
)

func attachUserRoute(router *mux.Router) {
	for _, r := range userRoutes {
		router.
			Methods(r.Method).
			Path(r.Pattern).
			Name(r.Name).
			Handler(wrapper.HandlerWrapper(r.HandlerFunc, r.Name))
	}
}

var userRoutes = route.Routes{

	// 2.1.1
	route.Route{
		"LaunchApp",
		"POST",
		"/launch",
		handlerv2.Dummy,
	},

	// 2.1.2
	route.Route{
		"GetUserInfo",
		"POST",
		"/info",
		handlerv2.UserInfo,
	},

	// 2.1.3
	route.Route{
		"UpdateUserInfo",
		"POST",
		"/info/update",
		handlerv2.UserInfoUpdate,
	},

	// 2.1.4
	route.Route{
		"UserGreeting",
		"POST",
		"/greeting",
		handlerv2.UserGreeting,
	},

	// 2.1.5
	route.Route{
		"UserNotification",
		"POST",
		"/notification",
		handlerv2.UserNotification,
	},

	// 2.1.6
	route.Route{
		"PromotionOnLogin",
		"POST",
		"/promotion/onlogin",
		handlerv2.UserPromotionOnLogin,
	},

	// 2.2.1
	route.Route{
		"UserProfile",
		"POST",
		"/profile",
		handlerv2.Dummy,
	},

	// 2.2.2
	route.Route{
		"TeacherProfile",
		"POST",
		"/teacher/profile",
		handlerv2.UserTeacherProfile,
	},

	// 2.2.3
	route.Route{
		"StudentProfile",
		"POST",
		"/student/profile",
		handlerv2.Dummy,
	},

	// 2.3.1
	route.Route{
		"UserSearch",
		"POST",
		"/search",
		handlerv2.UserSearch,
	},

	// 2.3.2
	route.Route{
		"UserTeacherSearch",
		"POST",
		"/teacher/search",
		handlerv2.UserTeacherSearch,
	},

	// 2.3.3
	route.Route{
		"UserStudentSearch",
		"POST",
		"/student/search",
		handlerv2.Dummy,
	},

	// 2.3.4
	route.Route{
		"TeacherRecent",
		"POST",
		"/teacher/recent",
		handlerv2.UserTeacherRecent,
	},

	// 2.3.5
	route.Route{
		"TeacherRecommendation",
		"POST",
		"/teacher/recommendation",
		handlerv2.UserTeacherRecommendation,
	},

	// 2.3.6
	route.Route{
		"ContactRecommendation",
		"POST",
		"/contact/recommendation",
		handlerv2.UserContactRecommendation,
	},

	// 2.4.1
	route.Route{
		"GetInvitationCode",
		"POST",
		"/invitation/code",
		handlerv2.Dummy,
	},

	// 2.4.2
	route.Route{
		"GetInvitationRecord",
		"POST",
		"/invitation/record",
		handlerv2.Dummy,
	},
}
