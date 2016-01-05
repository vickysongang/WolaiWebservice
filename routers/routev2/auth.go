package routev2

import (
	"github.com/gorilla/mux"

	"WolaiWebservice/handlers/handlerv2"
	"WolaiWebservice/routers/route"
	"WolaiWebservice/routers/wrapper"
)

func attachAuthRoute(router *mux.Router) {
	for _, r := range authRoutes {
		router.
			Methods(r.Method).
			Path(r.Pattern).
			Name(r.Name).
			Handler(wrapper.HandlerWrapper(r.HandlerFunc, r.Name))
	}
}

var authRoutes = route.Routes{

	// 1.1.1
	route.Route{
		"Register",
		"POST",
		"/register",
		handlerv2.Dummy,
	},

	// 1.1.2
	route.Route{
		"Login",
		"POST",
		"/login",
		handlerv2.Dummy,
	},

	// 1.1.3
	route.Route{
		"ForgotPassword",
		"POST",
		"/forgot_password",
		handlerv2.Dummy,
	},

	// 1.1.4
	route.Route{
		"SetPassword",
		"POST",
		"/set_password",
		handlerv2.Dummy,
	},

	// 1.1.5
	route.Route{
		"Logout",
		"POST",
		"/logout",
		handlerv2.Logout,
	},

	// 1.2.1
	route.Route{
		"SendSMSCode",
		"POST",
		"/phone/sms/code",
		handlerv2.AuthPhoneSMSCode,
	},

	// 1.2.2
	route.Route{
		"VerifySMSCode",
		"POST",
		"/phone/sms/verify",
		handlerv2.AuthPhoneSMSVerify,
	},

	// 1.2.3
	route.Route{
		"PhoneNumLogin",
		"POST",
		"/phone/login",
		handlerv2.AuthPhoneLogin,
	},
}
