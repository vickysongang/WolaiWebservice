package routev2

import (
	"github.com/gorilla/mux"

	"WolaiWebservice/handlers/handlerv2"
	"WolaiWebservice/routers/route"
	"WolaiWebservice/routers/wrapper"
)

func attachOauthRoute(router *mux.Router) {
	for _, r := range oauthRoutes {
		router.
			Methods(r.Method).
			Path(r.Pattern).
			Name(r.Name).
			Handler(wrapper.HandlerWrapper(r.HandlerFunc, r.Name, r.LogFlag, r.AuthFlag))
	}
}

var oauthRoutes = route.Routes{

	// 1.3.1
	route.Route{
		"QQLogin",
		"POST",
		"/qq/login",
		handlerv2.OauthQQLogin,
		true,
		false,
	},

	// 1.3.2
	route.Route{
		"QQRegister",
		"POST",
		"/qq/register",
		handlerv2.OauthQQRegister,
		true,
		false,
	},

	// 1.3.3
	route.Route{
		"QQBind",
		"POST",
		"/qq/bind",
		handlerv2.Dummy,
		false,
		false,
	},

	// 1.3.4
	route.Route{
		"QQUnbind",
		"POST",
		"/qq/unbind",
		handlerv2.Dummy,
		false,
		false,
	},
}
