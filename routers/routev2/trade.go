package routev2

import (
	"github.com/gorilla/mux"

	"WolaiWebservice/handlers/handlerv2"
	"WolaiWebservice/routers/route"
	"WolaiWebservice/routers/wrapper"
)

func attachTradeRoute(router *mux.Router) {
	for _, r := range tradeRoutes {
		router.
			Methods(r.Method).
			Path(r.Pattern).
			Name(r.Name).
			Handler(wrapper.HandlerWrapper(r.HandlerFunc, r.Name))
	}
}

var tradeRoutes = route.Routes{
	// 7.1.1
	route.Route{
		"UserBalance",
		"POST",
		"/user/balance",
		handlerv2.Dummy,
	},

	// 7.1.2
	route.Route{
		"UserTradeRecord",
		"POST",
		"/user/record",
		handlerv2.Dummy,
	},
}
