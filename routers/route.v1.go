package routers

import (
	"WolaiWebservice/routers/route"
	"WolaiWebservice/websocket"
)

var routesV1 = route.Routes{
	// Websocket
	route.Route{
		"V1WebSocketHandler",
		"GET",
		"/v1/ws",
		websocket.V1WebSocketHandler,
	},
}
