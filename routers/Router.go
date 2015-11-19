package routers

import (
	"encoding/json"
	"net/http"
	"time"

	"WolaiWebService/models"
	"WolaiWebService/utils"

	"WolaiWebService/handlers"

	seelog "github.com/cihub/seelog"
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
	//API V1
	for _, v1route := range V1Routes {
		var handler http.Handler
		handler = v1route.HandlerFunc
		handler = WebLogger(handler, v1route.Name)

		router.
			Methods(v1route.Method).
			Path(v1route.Pattern).
			Name(v1route.Name).
			Handler(handler)
	}

	//	//API V2
	//	for _, v2route := range V2Routes {
	//		var handler http.Handler
	//		handler = v2route.HandlerFunc
	//		//		handler = APIAuth(handler)
	//		handler = WebLogger(handler, v2route.Name)
	//		router.
	//			Methods(v2route.Method).
	//			Path(v2route.Pattern).
	//			Name(v2route.Name).
	//			Handler(handler)
	//	}

	return router
}

func WebLogger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		formData, _ := json.Marshal(r.Form)
		seelog.Info("[", r.Method, "] ", r.RequestURI, "|", name, "|", r.RemoteAddr, "|", r.UserAgent(), "\t", time.Since(start),
			"\t", string(formData))
	})
}

func APIAuth(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")
		encryptStr := utils.Encrypt(1, time.Now().Unix())
		if token != encryptStr {
			json.NewEncoder(w).Encode(models.NewPOIResponse(-1, "api auth fail", handlers.NullObject))
		} else {
			inner.ServeHTTP(w, r)
		}
	})
}
