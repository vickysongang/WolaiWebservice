package wrapper

import (
	"encoding/json"
	"net/http"
	"time"

	seelog "github.com/cihub/seelog"
)

func HandlerWrapper(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		formData, _ := json.Marshal(r.Form)
		seelog.Info("[", r.Method, "] ", r.RequestURI, "|", name, "\t", time.Since(start),
			"\t", string(formData))
	})
}

// func APIAuth(inner http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		token := r.Header.Get("token")
// 		encryptStr := utils.Encrypt(1, time.Now().Unix())
// 		if token != encryptStr {
// 			json.NewEncoder(w).Encode(models.NewPOIResponse(-1, "api auth fail", handlers.NullObject))
// 		} else {
// 			inner.ServeHTTP(w, r)
// 		}
// 	})
// }
