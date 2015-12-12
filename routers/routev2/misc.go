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
	// 9.1.1
	route.Route{
		"SendCloudHookPOST",
		"POST",
		"/hook/sendcloud",
		handlerv2.HookSendcloud,
	},
	route.Route{
		"SendCloudHookGET",
		"GET",
		"/hook/sendcloud",
		handlerv2.HookSendcloud,
	},

	// 9.1.2
	route.Route{
		"PingppHook",
		"POST",
		"/hook/pingpp",
		handlerv2.HookPingpp,
	},

	// 9.2.1
	route.Route{
		"HelpList",
		"POST",
		"/help/list",
		handlerv2.HelpList,
	},

	// 9.2.2
	route.Route{
		"GradeList",
		"POST",
		"/grade/list",
		handlerv2.GradeList,
	},

	// 9.2.3
	route.Route{
		"SubjectList",
		"POST",
		"/subject/list",
		handlerv2.SubjectList,
	},
}
