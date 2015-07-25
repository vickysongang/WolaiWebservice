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
	// Websocket
	Route{
		"V1WebSocket",
		"GET",
		"/v1/ws",
		V1WebSocketHandler,
	},

	// 1.1 Login
	Route{
		"V1LoginPOST",
		"POST",
		"/v1/login",
		V1Login,
	},
	Route{
		"V1LoginGET",
		"GET",
		"/v1/login",
		V1Login,
	},
	Route{
		"V1LoginGETURL",
		"GET",
		"/v1/login/{phone}",
		V1LoginGETURL,
	},

	// 1.2 Update profile
	Route{
		"V1UpdateProfilePOST",
		"POST",
		"/v1/update_profile",
		V1UpdateProfile,
	},
	Route{
		"V1UpdateProfileGET",
		"GET",
		"/v1/update_profile",
		V1UpdateProfile,
	},
	Route{
		"V1UpdateProfileGETURL",
		"GET",
		"/v1/update_profile/{userId}/{nickname}/{avatar}/{gender}",
		V1UpdateProfileGETURL,
	},

	// 1.3 Oauth Login
	Route{
		"V1OauthLoginPOST",
		"POST",
		"/v1/oauth/qq/login",
		V1OauthLogin,
	},
	Route{
		"V1OauthLoginGET",
		"GET",
		"/v1/oauth/qq/login",
		V1OauthLogin,
	},

	// 1.4 Oauth Register
	Route{
		"V1OauthRegisterPOST",
		"POST",
		"/v1/oauth/qq/register",
		V1OauthRegister,
	},
	Route{
		"V1OauthRegisterGET",
		"GET",
		"/v1/oauth/qq/register",
		V1OauthRegister,
	},

	// 2.1 Atrium
	Route{
		"V1AtriumPOST",
		"POST",
		"/v1/atrium",
		V1Atrium,
	},
	Route{
		"V1AtriumGET",
		"GET",
		"/v1/atrium",
		V1Atrium,
	},

	// 2.2 Feed Post
	Route{
		"V1FeedPostPOST",
		"POST",
		"/v1/feed/post",
		V1FeedPost,
	},
	Route{
		"V1FeedPostGET",
		"GET",
		"/v1/feed/post",
		V1FeedPost,
	},

	// 2.3 Feed Detial
	Route{
		"V1FeedDetailPOST",
		"POST",
		"/v1/feed/detail",
		V1FeedDetail,
	},
	Route{
		"V1FeedDetailGET",
		"GET",
		"/v1/feed/detail",
		V1FeedDetail,
	},

	// 2.4 Feed Like
	Route{
		"V1FeedLikePOST",
		"POST",
		"/v1/feed/like",
		V1FeedLike,
	},
	Route{
		"V1FeedLikeGET",
		"GET",
		"/v1/feed/like",
		V1FeedLike,
	},

	// 2.5 Feed Favorite
	/*
		Route{
			"V1FeedFavPOST",
			"POST",
			"/v1/feed/favorite",
			V1FeedFav,
		},
		Route{
			"V1FeedFavGET",
			"GET",
			"/v1/feed/favorite",
			V1FeedFav,
		},
	*/

	// 2.6 Feed Comment
	Route{
		"V1FeedCommentPOST",
		"POST",
		"/v1/feed/comment",
		V1FeedComment,
	},
	Route{
		"V1FeedCommentGET",
		"GET",
		"/v1/feed/comment",
		V1FeedComment,
	},

	// 2.7 Feed Comment Like
	/*
		Route{
			"V1FeedCommentLikePOST",
			"POST",
			"/v1/feed/comment/like",
			V1FeedCommentLike,
		},
		Route{
			"V1FeedCommentLikeGET",
			"GET",
			"/v1/feed/comment/like",
			V1FeedCommentLike,
		},
	*/
}
