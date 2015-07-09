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
}
