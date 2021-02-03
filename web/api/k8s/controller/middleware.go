package api

import (
	"net/http"
	"net/http/httptest"

	"github.com/sirupsen/logrus"
)

var whitelistPaths = []string{
	"/health",
}

func getWhiteList(path string) string {
	for _, p := range whitelistPaths {
		if p == path {
			return p
		}
	}

	return ""
}

func ValidatorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		next.ServeHTTP(w, r)
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		recorderWrite := httptest.NewRecorder()
		next.ServeHTTP(recorderWrite, r)

		for key := range recorderWrite.Header() {
			w.Header().Add(key, recorderWrite.Header().Get(key))
		}

		if recorderWrite.Code < 200 || recorderWrite.Code > 210 {
			logrus.WithFields(logrus.Fields{
				"err": recorderWrite.Body.String(),
			}).Warnln()
		}

		w.WriteHeader(recorderWrite.Code)
		recorderWrite.Body.WriteTo(w)
	})

}
