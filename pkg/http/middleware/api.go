package middleware

import (
	"net/http"
)

var URNs = map[string]struct{}{
	"/register": struct{}{},
}

func API(next http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if _, ok := URNs[r.URL.Path]; ok {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
		}
		next.ServeHTTP(w, r)
	}
}
