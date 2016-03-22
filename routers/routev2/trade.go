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
			Handler(wrapper.HandlerWrapper(r.HandlerFunc, r.Name, r.LogFlag, r.AuthFlag))
	}
}

var tradeRoutes = route.Routes{
	// 7.1.1
	route.Route{
		"UserBalance",
		"POST",
		"/user/balance",
		handlerv2.TradeUserBalance,
		true,
		true,
	},

	// 7.1.2
	route.Route{
		"UserTradeRecord",
		"POST",
		"/user/record",
		handlerv2.TradeUserRecord,
		true,
		true,
	},

	// 7.2.1
	route.Route{
		"TradeChargeBanner",
		"POST",
		"/charge/banner",
		handlerv2.TradeChargeBanner,
		true,
		true,
	},

	// 7.2.2
	route.Route{
		"TradeChargeShortcut",
		"POST",
		"/charge/shortcut",
		handlerv2.TradeChargeShortcut,
		true,
		true,
	},

	// 7.2.3
	route.Route{
		"TradeChargePremium",
		"POST",
		"/charge/premium",
		handlerv2.TradeChargePremium,
		true,
		true,
	},

	// 7.2.4
	route.Route{
		"TradeChargeCode",
		"POST",
		"/charge/code",
		handlerv2.TradeChargeCode,
		true,
		true,
	},

	// 7.3.1
	route.Route{
		"TradeQaPkgList",
		"POST",
		"/qapkg/list",
		handlerv2.TradeQaPkgList,
		true,
		false,
	},

	// 7.3.2
	route.Route{
		"TradeUserQaPkgDetail",
		"POST",
		"/qapkg/detail",
		handlerv2.TradeUserQaPkgDetail,
		true,
		true,
	},
}
