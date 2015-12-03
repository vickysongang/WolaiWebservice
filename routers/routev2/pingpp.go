package routev2

import (
	"github.com/gorilla/mux"

	"WolaiWebservice/handlers/handlerv2"
	"WolaiWebservice/routers/route"
	"WolaiWebservice/routers/wrapper"
)

func attachPingppRoute(router *mux.Router) {
	for _, r := range pingppRoutes {
		router.
			Methods(r.Method).
			Path(r.Pattern).
			Name(r.Name).
			Handler(wrapper.HandlerWrapper(r.HandlerFunc, r.Name))
	}
}

var pingppRoutes = route.Routes{
	// 7.2.1
	route.Route{
		"PingppPay",
		"POST",
		"/pay",
		handlerv2.Dummy,
	},

	// 7.2.2
	route.Route{
		"PingppRefund",
		"POST",
		"/refund",
		handlerv2.Dummy,
	},

	// 7.2.3
	route.Route{
		"PingppPayQuery",
		"POST",
		"/pay/query",
		handlerv2.Dummy,
	},

	// 7.2.4
	route.Route{
		"PingppPayRecord",
		"POST",
		"/pay/record",
		handlerv2.Dummy,
	},

	// 7.2.5
	route.Route{
		"PingppRefundQuery",
		"POST",
		"/refund/query",
		handlerv2.Dummy,
	},

	// 7.2.6
	route.Route{
		"PingppRefundRecord",
		"POST",
		"/refund/record",
		handlerv2.Dummy,
	},

	// 7.2.7
	route.Route{
		"PingppResult",
		"POST",
		"/result",
		handlerv2.Dummy,
	},
}
