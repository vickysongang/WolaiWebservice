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
			Handler(wrapper.HandlerWrapper(r.HandlerFunc, r.Name, r.LogFlag, r.AuthFlag))
	}
}

var miscRoutes = route.Routes{
	// 10.1.1
	route.Route{
		"SendCloudHookPOST",
		"POST",
		"/hook/sendcloud",
		handlerv2.HookSendcloud,
		false,
		false,
	},
	route.Route{
		"SendCloudHookGET",
		"GET",
		"/hook/sendcloud",
		handlerv2.HookSendcloud,
		false,
		false,
	},

	// 10.1.2
	route.Route{
		"PingppHook",
		"POST",
		"/hook/pingpp",
		handlerv2.HookPingpp,
		false,
		false,
	},

	// 10.2.1
	route.Route{
		"HelpList",
		"POST",
		"/help/list",
		handlerv2.HelpList,
		true,
		false,
	},

	// 10.2.2
	route.Route{
		"GradeList",
		"POST",
		"/grade/list",
		handlerv2.GradeList,
		false,
		false,
	},

	// 10.2.3
	route.Route{
		"SubjectList",
		"POST",
		"/subject/list",
		handlerv2.SubjectList,
		false,
		false,
	},

	// 10.2.4
	route.Route{
		"AdvBanner",
		"POST",
		"/adv/banner",
		handlerv2.AdvBanner,
		false,
		false,
	},

	// 10.2.5
	route.Route{
		"VersionUpgrade",
		"POST",
		"/version/upgrade",
		handlerv2.VersionUpgrade,
		false,
		false,
	},

	// 10.2.6
	route.Route{
		"GetQiniuDownloadUrl",
		"GET",
		"/qiniu/url/download",
		handlerv2.GetQiniuDownloadUrl,
		false,
		false,
	},

	// 10.2.7
	route.Route{
		"GetQiniuUploadToken",
		"POST",
		"/qiniu/token/upload",
		handlerv2.GetQiniuUploadToken,
		false,
		true,
	},

	route.Route{
		"Dummy",
		"GET",
		"/dummy",
		handlerv2.Dummy,
		false,
		false,
	},

	route.Route{
		"Dummy2",
		"GET",
		"/dummy2",
		handlerv2.Dummy2,
		false,
		false,
	},

	route.Route{
		"Dummy3",
		"GET",
		"/dummy3",
		handlerv2.Dummy3,
		false,
		false,
	},
}
