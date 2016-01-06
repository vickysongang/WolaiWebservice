package route

import (
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
	LogFlag     bool
	AuthFlag    bool
}

type Routes []Route
