package routev2

import (
	"github.com/gorilla/mux"

	"WolaiWebservice/handlers/handlerv2"
	"WolaiWebservice/routers/route"
	"WolaiWebservice/routers/wrapper"
)

func attachCourseRoute(router *mux.Router) {
	for _, r := range courseRoutes {
		router.
			Methods(r.Method).
			Path(r.Pattern).
			Name(r.Name).
			Handler(wrapper.HandlerWrapper(r.HandlerFunc, r.Name, r.LogFlag, r.AuthFlag))
	}
}

var courseRoutes = route.Routes{

	// 9.1.1
	route.Route{
		"CourseBanner",
		"POST",
		"/banner",
		handlerv2.CourseBanner,
		true,
		true,
	},

	// 9.1.2
	route.Route{
		"CourseHomePage",
		"POST",
		"/homepage",
		handlerv2.CourseHomePage,
		true,
		true,
	},

	// 9.1.3
	route.Route{
		"CourseModuleAll",
		"POST",
		"/module/all",
		handlerv2.CourseModuleAll,
		true,
		true,
	},

	// 9.2.1
	route.Route{
		"CourseListStudent",
		"POST",
		"/user/list/student",
		handlerv2.CourseListStudent,
		true,
		true,
	},

	// 9.2.2
	route.Route{
		"CourseListTeacher",
		"POST",
		"/user/list/teacher",
		handlerv2.CourseListTeacher,
		true,
		true,
	},

	// 9.3.1
	route.Route{
		"CourseDetailStudent",
		"POST",
		"/detail/student",
		handlerv2.CourseDetailStudent,
		true,
		true,
	},

	// 9.3.2
	route.Route{
		"CourseDetailTeacher",
		"POST",
		"/detail/teacher",
		handlerv2.CourseDetailTeacher,
		true,
		true,
	},

	// 9.3.3
	route.Route{
		"CourseDetailStudentUpgrade",
		"POST",
		"/upgrade/detail/student",
		handlerv2.CourseDetailStudentUpgrade,
		true,
		true,
	},

	// 9.3.4
	route.Route{
		"CourseDetailTeacherUpgrade",
		"POST",
		"/upgrade/detail/teacher",
		handlerv2.CourseDetailTeacherUpgrade,
		true,
		true,
	},

	// 9.4.1
	route.Route{
		"CourseActionProceed",
		"POST",
		"/action/proceed",
		handlerv2.CourseActionProceed,
		true,
		true,
	},

	// 9.4.2
	route.Route{
		"CourseActionQuickbuy",
		"POST",
		"/action/quickbuy",
		handlerv2.CourseActionQuickbuy,
		true,
		true,
	},

	// 9.4.3
	route.Route{
		"CourseActionPay",
		"POST",
		"/action/pay",
		handlerv2.CourseActionPay,
		true,
		true,
	},

	// 9.4.4
	route.Route{
		"CourseActionNextChapter",
		"POST",
		"/action/nextchapter",
		handlerv2.CourseActionNextChapter,
		true,
		true,
	},

	// 9.4.5
	route.Route{
		"CourseActionAuditionCheck",
		"POST",
		"/action/audition/check",
		handlerv2.CourseActionAuditionCheck,
		true,
		true,
	},

	// 9.4.6
	route.Route{
		"CourseAuditionActionProceed",
		"POST",
		"/action/audition/proceed",
		handlerv2.CourseAuditionActionProceed,
		true,
		false,
	},

	// 9.4.7
	route.Route{
		"CourseDeluxeActionProceed",
		"POST",
		"/action/deluxe/proceed",
		handlerv2.CourseDeluxeActionProceed,
		true,
		true,
	},

	// 9.5.1
	route.Route{
		"CourseAttachs",
		"POST",
		"/attachs",
		handlerv2.CourseAttachs,
		true,
		true,
	},

	// 9.5.2
	route.Route{
		"CourseChapterAttachs",
		"POST",
		"/chapter/attachs",
		handlerv2.CourseChapterAttachs,
		true,
		true,
	},

	// 9.6.1
	route.Route{
		"CourseCountOfConversation",
		"POST",
		"/conversation/count",
		handlerv2.CourseCountOfConversation,
		true,
		true,
	},

	// 9.6.2
	route.Route{
		"CourseListStudentOfConversation",
		"POST",
		"/conversation/user/list",
		handlerv2.CourseListStudentOfConversation,
		true,
		true,
	},

	// 9.7.1
	route.Route{
		"CourseRenewWaitingRecordDetail",
		"POST",
		"/renew/waiting/detail",
		handlerv2.CourseRenewWaitingRecordDetail,
		true,
		true,
	},
}
