package routers

import (
	"net/http"

	"github.com/gorilla/mux"

	"WolaiWebservice/routers/routev2"
	"WolaiWebservice/routers/wrapper"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	//API V1
	for _, v1route := range routesV1 {
		var handler http.Handler
		handler = v1route.HandlerFunc
		handler = wrapper.HandlerWrapper(handler, v1route.Name)

		router.
			Methods(v1route.Method).
			Path(v1route.Pattern).
			Name(v1route.Name).
			Handler(handler)
	}

	routev2.AttachRoute(router)

	return router
}
