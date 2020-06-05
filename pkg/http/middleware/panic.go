package middleware

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func Panic(logger *logrus.Logger, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("panic: [%s] - %s", r.URL.Path, err)
				http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	}
}
