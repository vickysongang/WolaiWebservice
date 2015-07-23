package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

var routes = Routes{
	/*
		Route{
			"V1LoginPOST",
			"POST",
			"/v1/login",
			V1LoginPOST,
		},
	*/
	Route{
		"V1LoginGET",
		"GET",
		"/v1/login",
		V1LoginGET,
	},
	Route{
		"V1LoginGETURL",
		"GET",
		"/v1/login/{phone}",
		V1LoginGETURL,
	},
	/*
		Route{
			"V1UpdateProfilePOST",
			"POST",
			"/v1/update_profile",
			V1UpdateProfilePOST,
		},
	*/
	Route{
		"V1UpdateProfileGET",
		"GET",
		"/v1/update_profile",
		V1UpdateProfileGET,
	},
	Route{
		"V1UpdateProfileGETURL",
		"GET",
		"/v1/update_profile/{userId}/{nickname}/{avatar}/{gender}",
		V1UpdateProfileGETURL,
	},
	Route{
		"V1WebSocket",
		"GET",
		"/v1/ws",
		V1WebSocketHandler,
	},

	Route{
		"V1AtriumGET",
		"GET",
		"/v1/atrium",
		V1AtriumGET,
	},
	Route{
		"V1FeedPostGET",
		"GET",
		"/v1/feed/post",
		V1FeedPostGET,
	},
	Route{
		"V1FeedDetailGET",
		"GET",
		"/v1/feed/detail",
		V1FeedDetailGET,
	},
	Route{
		"V1FeedCommentGET",
		"GET",
		"/v1/feed/comment",
		V1FeedCommentGET,
	},
	Route{
		"V1FeedLikeGET",
		"GET",
		"/v1/feed/like",
		V1FeedLikeGET,
	},
	Route{
		"V1FeedFavGET",
		"GET",
		"/v1/feed/favorite",
		V1FeedFavGET,
	},
	Route{
		"V1FeedCommentLikeGET",
		"GET",
		"/v1/feed/comment/like",
		V1FeedCommentLikeGET,
	},
}
