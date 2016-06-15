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
			Handler(wrapper.HandlerWrapper(r.HandlerFunc, r.Name, r.LogFlag, r.AuthFlag))
	}
}

var userRoutes = route.Routes{

	// 2.1.1
	route.Route{
		"LaunchApp",
		"POST",
		"/launch",
		handlerv2.UserLaunch,
		true,
		true,
	},

	// 2.1.2
	route.Route{
		"GetUserInfo",
		"POST",
		"/info",
		handlerv2.UserInfo,
		true,
		true,
	},

	// 2.1.3
	route.Route{
		"UpdateUserInfo",
		"POST",
		"/info/update",
		handlerv2.UserInfoUpdate,
		true,
		true,
	},

	// 2.1.4
	route.Route{
		"UserGreeting",
		"POST",
		"/greeting",
		handlerv2.UserGreeting,
		true,
		true,
	},

	// 2.1.5
	route.Route{
		"UserNotification",
		"POST",
		"/notification",
		handlerv2.UserNotification,
		true,
		true,
	},

	// 2.1.6
	route.Route{
		"PromotionOnLogin",
		"POST",
		"/promotion/onlogin",
		handlerv2.UserPromotionOnLogin,
		true,
		true,
	},

	// 2.2.1
	route.Route{
		"UserProfile",
		"POST",
		"/profile",
		handlerv2.Dummy,
		true,
		true,
	},

	// 2.2.2
	route.Route{
		"TeacherProfile",
		"POST",
		"/teacher/profile",
		handlerv2.UserTeacherProfile,
		true,
		true,
	},

	// 2.2.3
	route.Route{
		"TeacherProfileCourse",
		"POST",
		"/teacher/profile/course",
		handlerv2.UserTeacherProfileCourse,
		true,
		true,
	},

	// 2.2.4
	route.Route{
		"TeacherProfileEvaluation",
		"POST",
		"/teacher/profile/evaluation",
		handlerv2.UserTeacherProfileEvalution,
		true,
		true,
	},

	// 2.2.5
	route.Route{
		"StudentProfile",
		"POST",
		"/student/profile",
		handlerv2.UserStudentProfile,
		true,
		true,
	},

	// 2.2.6
	route.Route{
		"UpdateStudentProfile",
		"POST",
		"/student/profile/update",
		handlerv2.UserStudentProfileUpdate,
		true,
		true,
	},

	// 2.2.7
	route.Route{
		"CompleteStudentProfile",
		"POST",
		"/student/profile/complete",
		handlerv2.UserStudentProfileComplete,
		true,
		true,
	},

	// 2.2.8
	route.Route{
		"TeacherProfileChecked",
		"POST",
		"/teacher/profile/checked",
		handlerv2.UserTeacherProfileChecked,
		true,
		true,
	},

	// 2.3.1
	route.Route{
		"UserSearch",
		"POST",
		"/search",
		handlerv2.UserSearch,
		true,
		true,
	},

	// 2.3.2
	route.Route{
		"UserTeacherSearch",
		"POST",
		"/teacher/search",
		handlerv2.UserTeacherSearch,
		true,
		true,
	},

	// 2.3.3
	route.Route{
		"UserStudentSearch",
		"POST",
		"/student/search",
		handlerv2.Dummy,
		true,
		true,
	},

	// 2.3.4
	route.Route{
		"TeacherRecent",
		"POST",
		"/teacher/recent",
		handlerv2.UserTeacherRecent,
		true,
		true,
	},

	// 2.3.5
	route.Route{
		"TeacherRecommendation",
		"POST",
		"/teacher/recommendation",
		handlerv2.UserTeacherRecommendation,
		true,
		true,
	},

	// 2.3.6
	route.Route{
		"ContactRecommendation",
		"POST",
		"/contact/recommendation",
		handlerv2.UserContactRecommendation,
		true,
		true,
	},

	// 2.4.1
	route.Route{
		"GetInvitationCode",
		"POST",
		"/invitation/code",
		handlerv2.Dummy,
		false,
		true,
	},

	// 2.4.2
	route.Route{
		"GetInvitationRecord",
		"POST",
		"/invitation/record",
		handlerv2.Dummy,
		false,
		true,
	},

	// 2.5.1
	route.Route{
		"GetUserDataUsage",
		"POST",
		"/data/usage",
		handlerv2.UserDataUsage,
		true,
		true,
	},

	// 2.5.2
	route.Route{
		"UpdateUserDataUsage",
		"POST",
		"/data/usage/update",
		handlerv2.UserDataUsageUpdate,
		true,
		true,
	},
}
