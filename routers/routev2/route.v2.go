package routev2

import (
	"github.com/gorilla/mux"
)

func AttachRoute(r *mux.Router) {
	v2Router := r.PathPrefix("/v2/").Subrouter()

	authRouter := v2Router.PathPrefix("/auth/").Subrouter()
	attachAuthRoute(authRouter)

	oauthRouter := v2Router.PathPrefix("/oauth/").Subrouter()
	attachOauthRoute(oauthRouter)

	userRouter := v2Router.PathPrefix("/user/").Headers(
		"X-Wolai-Token", "",
		"X-Wolai-ID", "").Subrouter()
	attachUserRoute(userRouter)

	feedRouter := v2Router.PathPrefix("/feed/").Headers(
		"X-Wolai-Token", "",
		"X-Wolai-ID", "").Subrouter()
	attachFeedRoute(feedRouter)

	msgRouter := v2Router.PathPrefix("/message/").Headers(
		"X-Wolai-Token", "",
		"X-Wolai-ID", "").Subrouter()
	attachMessageRoute(msgRouter)

	orderRouter := v2Router.PathPrefix("/order/").Headers(
		"X-Wolai-Token", "",
		"X-Wolai-ID", "").Subrouter()
	attachOrderRoute(orderRouter)

	sessionRouter := v2Router.PathPrefix("/session/").Headers(
		"X-Wolai-Token", "",
		"X-Wolai-ID", "").Subrouter()
	attachSessionRoute(sessionRouter)

	tradeRouter := v2Router.PathPrefix("/trade/").Headers(
		"X-Wolai-Token", "",
		"X-Wolai-ID", "").Subrouter()
	attachTradeRoute(tradeRouter)

	pingppRouter := v2Router.PathPrefix("/pingpp/").Headers(
		"X-Wolai-Token", "",
		"X-Wolai-ID", "").Subrouter()
	attachPingppRoute(pingppRouter)

	courseRouter := v2Router.PathPrefix("/course/").Headers(
		"X-Wolai-Token", "",
		"X-Wolai-ID", "").Subrouter()
	attachCourseRoute(courseRouter)

	attachMiscRoute(v2Router)
}
