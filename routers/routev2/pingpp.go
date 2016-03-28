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
			Handler(wrapper.HandlerWrapper(r.HandlerFunc, r.Name, r.LogFlag, r.AuthFlag))
	}
}

var pingppRoutes = route.Routes{
	// 8.1.1
	route.Route{
		"PingppPay",
		"POST",
		"/pay",
		handlerv2.PingppPay,
		true,
		false,
	},

	// 8.1.2
	route.Route{
		"PingppPayQuery",
		"POST",
		"/pay/query",
		handlerv2.PingppPayQuery,
		false,
		true,
	},

	// 8.1.3
	route.Route{
		"PingppPayRecord",
		"POST",
		"/pay/record",
		handlerv2.PingppPayRecord,
		false,
		true,
	},

	// 8.2.1
	route.Route{
		"PingppRefund",
		"POST",
		"/refund",
		handlerv2.PingppRefund,
		true,
		true,
	},

	// 8.2.2
	route.Route{
		"PingppRefundQuery",
		"POST",
		"/refund/query",
		handlerv2.PingppRefundQuery,
		false,
		true,
	},

	// 8.2.3
	route.Route{
		"PingppRefundRecord",
		"POST",
		"/refund/record",
		handlerv2.PingppRefundRecord,
		false,
		true,
	},

	// 8.3.1
	route.Route{
		"PingppResult",
		"POST",
		"/result",
		handlerv2.PingppResult,
		false,
		true,
	},
}
