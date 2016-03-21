package routev2

import (
	"github.com/gorilla/mux"

	"WolaiWebservice/handlers/handlerv2"
	"WolaiWebservice/routers/route"
	"WolaiWebservice/routers/wrapper"
)

func attachCometRoute(router *mux.Router) {
	for _, r := range cometRoutes {
		router.
			Methods(r.Method).
			Path(r.Pattern).
			Name(r.Name).
			Handler(wrapper.HandlerWrapperLogResponse(r.HandlerFunc, r.Name, r.LogFlag, r.AuthFlag))
	}
}

var cometRoutes = route.Routes{
	// 5.1.1
	route.Route{
		"HandleMessage",
		"POST",
		"/message/handle",
		handlerv2.HandleCometMessage,
		true,
		true,
	},
}
