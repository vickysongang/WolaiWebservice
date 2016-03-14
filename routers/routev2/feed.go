package routev2

import (
	"github.com/gorilla/mux"

	"WolaiWebservice/handlers/handlerv2"
	"WolaiWebservice/routers/route"
	"WolaiWebservice/routers/wrapper"
)

func attachFeedRoute(router *mux.Router) {
	for _, r := range feedRoutes {
		router.
			Methods(r.Method).
			Path(r.Pattern).
			Name(r.Name).
			Handler(wrapper.HandlerWrapper(r.HandlerFunc, r.Name, r.LogFlag, r.AuthFlag))
	}
}

var feedRoutes = route.Routes{

	// 3.1.1
	route.Route{
		"GetFeedFlow",
		"POST",
		"/atrium",
		handlerv2.FeedAtrium,
		true,
		true,
	},

	// 3.1.2
	route.Route{
		"GetFeedFlowStick",
		"POST",
		"/atrium/stick",
		handlerv2.FeedAtriumStick,
		true,
		true,
	},

	// 3.1.3
	route.Route{
		"PostFeed",
		"POST",
		"/post",
		handlerv2.FeedPost,
		true,
		true,
	},

	// 3.1.4
	route.Route{
		"GetFeedDetail",
		"POST",
		"/detail",
		handlerv2.FeedDetail,
		true,
		false,
	},

	// 3.1.5
	route.Route{
		"LikeFeed",
		"POST",
		"/like",
		handlerv2.FeedLike,
		true,
		true,
	},

	// 3.1.6
	route.Route{
		"CommentFeed",
		"POST",
		"/comment",
		handlerv2.FeedComment,
		true,
		true,
	},

	// 3.2.1
	route.Route{
		"UserFeedHistory",
		"POST",
		"/user/history",
		handlerv2.FeedUserHistory,
		true,
		true,
	},

	// 3.2.2
	route.Route{
		"UserFeedLike",
		"POST",
		"/user/like",
		handlerv2.FeedUserLike,
		true,
		true,
	},
}
