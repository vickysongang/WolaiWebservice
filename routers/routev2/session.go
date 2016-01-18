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
			Handler(wrapper.HandlerWrapper(r.HandlerFunc, r.Name, r.LogFlag, r.AuthFlag))
	}
}

var sessionRoutes = route.Routes{
	// 6.1.1
	route.Route{
		"GetSessionInfo",
		"POST",
		"/info",
		handlerv2.SessionInfo,
		true,
		true,
	},

	// 6.1.2
	route.Route{
		"GetUserSessionRecord",
		"POST",
		"/user/record",
		handlerv2.SessionUserRecord,
		true,
		true,
	},

	// 6.2.1
	route.Route{
		"SessionSeekHelp",
		"POST",
		"/seek_help",
		handlerv2.SessionSeekHelp,
		true,
		true,
	},

	// 6.2.2
	route.Route{
		"SessionQACardCatalog",
		"POST",
		"/qacard/catalog",
		handlerv2.SessionQACardCatalog,
		true,
		true,
	},

	// 6.2.3
	route.Route{
		"SessionQACardFetch",
		"POST",
		"/qacard/fetch",
		handlerv2.SessionQACardFetch,
		true,
		true,
	},

	// 6.3.1
	route.Route{
		"SessionEvaluationLabel",
		"POST",
		"/evaluation/label/list",
		handlerv2.SessionEvaluationLabelList,
		true,
		true,
	},

	// 6.3.2
	route.Route{
		"SessionEvaluationCreate",
		"POST",
		"/evaluation/label/post",
		handlerv2.SessionEvaluationLabelPost,
		true,
		true,
	},

	// 6.3.3
	route.Route{
		"SessionEvaluationResult",
		"POST",
		"/evaluation/label/result",
		handlerv2.SessionEvaluationLabelResult,
		true,
		true,
	},

	// 6.4.1
	route.Route{
		"SessionComplain",
		"POST",
		"/complain/post",
		handlerv2.SessionComplainPost,
		true,
		true,
	},

	// 6.4.2
	route.Route{
		"SessionComplainCheck",
		"POST",
		"/complain/check",
		handlerv2.SessionComplainCheck,
		true,
		true,
	},

	// 6.5.1
	route.Route{
		"SessionWhiteboardCall",
		"POST",
		"/whiteboard/call",
		handlerv2.SessionWhiteboardCall,
		true,
		true,
	},
}
