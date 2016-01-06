package routev2

import (
	"github.com/gorilla/mux"

	"WolaiWebservice/handlers/handlerv2"
	"WolaiWebservice/routers/route"
	"WolaiWebservice/routers/wrapper"
)

func attachOrderRoute(router *mux.Router) {
	for _, r := range orderRoutes {
		router.
			Methods(r.Method).
			Path(r.Pattern).
			Name(r.Name).
			Handler(wrapper.HandlerWrapper(r.HandlerFunc, r.Name, r.LogFlag, r.AuthFlag))
	}
}

var orderRoutes = route.Routes{
	// 5.1.1
	route.Route{
		"CreateOrder",
		"POST",
		"/create",
		handlerv2.OrderCreate,
		true,
		true,
	},

	// 5.1.2
	route.Route{
		"OrderExpectPrice",
		"POST",
		"/expectation",
		handlerv2.OrderExpectation,
		true,
		true,
	},

	// 5.1.3
	route.Route{
		"OrderInfo",
		"POST",
		"/info",
		handlerv2.OrderInfo,
		true,
		true,
	},
}
