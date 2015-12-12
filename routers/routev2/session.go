package routev2

import (
	"github.com/gorilla/mux"

	"WolaiWebservice/handlers/handlerv2"
	"WolaiWebservice/routers/route"
	"WolaiWebservice/routers/wrapper"
)

func attachSessionRoute(router *mux.Router) {
	for _, r := range sessionRoutes {
		router.
			Methods(r.Method).
			Path(r.Pattern).
			Name(r.Name).
			Handler(wrapper.HandlerWrapper(r.HandlerFunc, r.Name))
	}
}

var sessionRoutes = route.Routes{
	// 6.1.1
	route.Route{
		"GetSessionInfo",
		"POST",
		"/info",
		handlerv2.SessionInfo,
	},

	// 6.1.2
	route.Route{
		"GetUserSessionRecord",
		"POST",
		"/user/record",
		handlerv2.SessionUserRecord,
	},

	// 6.2.1
	route.Route{
		"SessionSeekHelp",
		"POST",
		"/seek_help",
		handlerv2.SessionSeekHelp,
	},

	// 6.3.1
	route.Route{
		"SessionEvaluationLabel",
		"POST",
		"/evaluation/label/list",
		handlerv2.SessionEvaluationLabelList,
	},

	// 6.3.2
	route.Route{
		"SessionEvaluationCreate",
		"POST",
		"/evaluation/label/post",
		handlerv2.SessionEvaluationLabelPost,
	},

	// 6.3.3
	route.Route{
		"SessionEvaluationResult",
		"POST",
		"/evaluation/label/result",
		handlerv2.SessionEvaluationLabelResult,
	},

	// 6.4.1
	route.Route{
		"SessionComplain",
		"POST",
		"/complain/post",
		handlerv2.SessionComplainPost,
	},

	// 6.4.2
	route.Route{
		"SessionComplainCheck",
		"POST",
		"/complain/check",
		handlerv2.SessionComplainCheck,
	},
}
