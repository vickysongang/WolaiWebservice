package routev2

import (
	"github.com/gorilla/mux"

	"WolaiWebservice/handlers/handlerv2"
	"WolaiWebservice/routers/route"
	"WolaiWebservice/routers/wrapper"
)

func attachMessageRoute(router *mux.Router) {
	for _, r := range messageRoutes {
		router.
			Methods(r.Method).
			Path(r.Pattern).
			Name(r.Name).
			Handler(wrapper.HandlerWrapper(r.HandlerFunc, r.Name))
	}
}

var messageRoutes = route.Routes{
	// 4.1.1
	route.Route{
		"GetConversationID",
		"POST",
		"/conversation/create",
		handlerv2.Dummy,
	},
}
