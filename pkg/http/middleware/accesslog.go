package middleware

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func AccessLog(logger *logrus.Logger, next http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logger.Infof("[%s]: %s (%f sec)\n", r.Method, r.URL.Path, time.Since(start).Seconds())
	}
}
