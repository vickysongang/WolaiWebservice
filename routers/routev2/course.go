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
			Handler(wrapper.HandlerWrapper(r.HandlerFunc, r.Name))
	}
}

var courseRoutes = route.Routes{

	// 9.1.1
	route.Route{
		"CourseBanner",
		"POST",
		"/banner",
		handlerv2.CourseBanner,
	},

	// 9.1.2
	route.Route{
		"CourseHomePage",
		"POST",
		"/homepage",
		handlerv2.CourseHomePage,
	},

	// 9.1.3
	route.Route{
		"CourseModuleAll",
		"POST",
		"/module/all",
		handlerv2.CourseModuleAll,
	},

	// 9.2.1
	route.Route{
		"CourseListStudent",
		"POST",
		"/user/list/student",
		handlerv2.Dummy,
	},

	// 9.2.2
	route.Route{
		"CourseListTeacher",
		"POST",
		"/user/list/teacher",
		handlerv2.Dummy,
	},

	// 9.3.1
	route.Route{
		"CourseDetailStudent",
		"POST",
		"/detail/student",
		handlerv2.Dummy,
	},

	// 9.3.2
	route.Route{
		"CourseDetailTeacher",
		"POST",
		"/detail/teacher",
		handlerv2.Dummy,
	},

	// 9.5.1
	route.Route{
		"CourseAttachs",
		"POST",
		"/attachs",
		handlerv2.CourseAttachs,
	},

	// 9.5.2
	route.Route{
		"CourseChapterAttachs",
		"POST",
		"/chapter/attachs",
		handlerv2.CourseChapterAttachs,
	},
}
