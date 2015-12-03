package routev2

import (
	"github.com/gorilla/mux"

	"WolaiWebservice/handlers/handlerv2"
	"WolaiWebservice/routers/route"
	"WolaiWebservice/routers/wrapper"
)

func attachMiscRoute(router *mux.Router) {
	for _, r := range miscRoutes {
		router.
			Methods(r.Method).
			Path(r.Pattern).
			Name(r.Name).
			Handler(wrapper.HandlerWrapper(r.HandlerFunc, r.Name))
	}
}

var miscRoutes = route.Routes{
	// 8.1.1
	route.Route{
		"SendCloudHook",
		"POST",
		"/hook/sendcloud",
		handlerv2.Dummy,
	},

	// 8.1.2
	route.Route{
		"PingppHook",
		"POST",
		"/hook/pingpp",
		handlerv2.Dummy,
	},

	// 8.2.1
	route.Route{
		"HelpList",
		"POST",
		"/help/list",
		handlerv2.Dummy,
	},

	// 8.2.2
	route.Route{
		"GradeList",
		"POST",
		"/grade/list",
		handlerv2.Dummy,
	},

	// 8.2.3
	route.Route{
		"SubjectList",
		"POST",
		"/subject/list",
		handlerv2.Dummy,
	},
}
