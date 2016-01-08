package wrapper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	seelog "github.com/cihub/seelog"

	"WolaiWebservice/routers/token"
)

func HandlerWrapper(inner http.Handler, name string, logFlag bool, authFlag bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userId int64
		var success bool

		if authFlag {
			userId, success = authenticate(r)

			if !success {
				http.Error(w,
					http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

				if logFlag {
					logEntry := fmt.Sprintf("[%s]\t%s|%s|%d\t%s",
						r.Method, r.RequestURI, name, userId, "401 Unauthorized")
					seelog.Error(logEntry)
				}

				return
			}
		}

		start := time.Now()

		inner.ServeHTTP(w, r)

		if logFlag {
			formData, _ := json.Marshal(r.Form)

			logEntry := fmt.Sprintf("[%s]\t%s|%s|%d\t%s\t%s",
				r.Method, r.RequestURI, name, userId, time.Since(start).String(), string(formData))
			seelog.Info(logEntry)
		}
	})
}

func authenticate(r *http.Request) (int64, bool) {
	userIdStr := r.Header.Get("X-Wolai-ID")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		return 0, false
	}
	tokenString := r.Header.Get("X-Wolai-Token")

	manager := token.GetTokenManager()
	err = manager.TokenAuthenticate(userId, tokenString)

	if err != nil {
		return userId, false
	}

	return userId, true
}
